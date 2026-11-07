package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed builtin/*.json
var builtinFS embed.FS

// Manager manages both built-in and user custom templates.
type Manager struct {
	customDir string
	templates map[string]*Template
	mu        sync.RWMutex
}

// NewManager creates and initializes a Template Manager.
func NewManager(customDir string) (*Manager, error) {
	m := &Manager{
		customDir: customDir,
		templates: make(map[string]*Template),
	}

	if err := m.loadBuiltinTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load built-in templates: %w", err)
	}

	if customDir != "" {
		if err := os.MkdirAll(customDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create custom templates dir %s: %w", customDir, err)
		}
		if err := m.loadCustomTemplates(); err != nil {
			return nil, fmt.Errorf("failed to load custom templates: %w", err)
		}
	}

	return m, nil
}

func (m *Manager) loadBuiltinTemplates() error {
	entries, err := builtinFS.ReadDir("builtin")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := builtinFS.ReadFile("builtin/" + entry.Name())
		if err != nil {
			return err
		}

		var tmpl Template
		if err := json.Unmarshal(data, &tmpl); err != nil {
			return fmt.Errorf("invalid json in builtin template %s: %w", entry.Name(), err)
		}

		tmpl.Builtin = true
		if tmpl.Width <= 0 {
			tmpl.Width = 1200
		}
		if tmpl.Height <= 0 {
			tmpl.Height = 630
		}

		m.templates[tmpl.ID] = &tmpl
	}

	return nil
}

func (m *Manager) loadCustomTemplates() error {
	entries, err := os.ReadDir(m.customDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		fullPath := filepath.Join(m.customDir, entry.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		var tmpl Template
		if err := json.Unmarshal(data, &tmpl); err != nil {
			continue
		}

		tmpl.Builtin = false
		if tmpl.Width <= 0 {
			tmpl.Width = 1200
		}
		if tmpl.Height <= 0 {
			tmpl.Height = 630
		}

		m.templates[tmpl.ID] = &tmpl
	}

	return nil
}

// List returns a summary list of all registered templates.
func (m *Manager) List() []TemplateSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]TemplateSummary, 0, len(m.templates))
	for _, t := range m.templates {
		list = append(list, TemplateSummary{
			ID:            t.ID,
			Name:          t.Name,
			Description:   t.Description,
			Builtin:       t.Builtin,
			ElementsCount: len(t.Elements),
		})
	}
	return list
}

// Get retrieves a template by ID.
func (m *Manager) Get(id string) (*Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Normalize ID (trim .png or .json if present)
	id = strings.TrimSuffix(id, ".png")
	id = strings.TrimSuffix(id, ".webp")
	id = strings.TrimSuffix(id, ".json")

	tmpl, exists := m.templates[id]
	if !exists {
		return nil, fmt.Errorf("template %q not found", id)
	}

	return tmpl, nil
}

// Save stores or updates a custom template both in memory and on disk.
func (m *Manager) Save(tmpl Template) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.templates[tmpl.ID]; ok && existing.Builtin {
		return fmt.Errorf("cannot overwrite built-in template %q", tmpl.ID)
	}

	tmpl.Builtin = false
	if tmpl.Width <= 0 {
		tmpl.Width = 1200
	}
	if tmpl.Height <= 0 {
		tmpl.Height = 630
	}

	if m.customDir != "" {
		filePath := filepath.Join(m.customDir, tmpl.ID+".json")
		data, err := json.MarshalIndent(tmpl, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write custom template file: %w", err)
		}
	}

	m.templates[tmpl.ID] = &tmpl
	return nil
}

// Delete removes a custom template.
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tmpl, exists := m.templates[id]
	if !exists {
		return fmt.Errorf("template %q not found", id)
	}

	if tmpl.Builtin {
		return fmt.Errorf("cannot delete built-in template %q", id)
	}

	delete(m.templates, id)

	if m.customDir != "" {
		filePath := filepath.Join(m.customDir, id+".json")
		_ = os.Remove(filePath)
	}

	return nil
}
