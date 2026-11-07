package renderer_test

import (
	"bytes"
	"image/png"
	"testing"
	"time"

	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/renderer"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
)

func TestRenderAllBuiltinTemplates(t *testing.T) {
	if err := fonts.Init(); err != nil {
		t.Fatalf("failed to init fonts: %v", err)
	}

	mgr, err := templates.NewManager("")
	if err != nil {
		t.Fatalf("failed to create template manager: %v", err)
	}

	list := mgr.List()
	if len(list) < 6 {
		t.Fatalf("expected at least 6 built-in templates, got %d", len(list))
	}

	for _, summary := range list {
		tmpl, err := mgr.Get(summary.ID)
		if err != nil {
			t.Fatalf("failed to get template %s: %v", summary.ID, err)
		}

		customParams := map[string]string{
			"title":       "Dynamic Title for " + summary.Name,
			"author":      "Test Engineer",
			"tag":         "PERFORMANCE",
			"description": "Validating rendering speed and memory safety.",
			"date":        "2026-11-07",
		}

		start := time.Now()
		imgBytes, err := renderer.RenderCard(tmpl, customParams)
		dur := time.Since(start)

		if err != nil {
			t.Errorf("rendering template %s failed: %v", summary.ID, err)
			continue
		}

		if len(imgBytes) == 0 {
			t.Errorf("template %s rendered empty byte array", summary.ID)
			continue
		}

		// Verify valid PNG
		cfg, err := png.DecodeConfig(bytes.NewReader(imgBytes))
		if err != nil {
			t.Errorf("template %s produced invalid PNG: %v", summary.ID, err)
			continue
		}

		if cfg.Width != 1200 || cfg.Height != 630 {
			t.Errorf("template %s wrong dimensions: got %dx%d, want 1200x630", summary.ID, cfg.Width, cfg.Height)
		}

		t.Logf("Template %-15s rendered in %6.2fms (PNG size: %d KB)", summary.ID, float64(dur.Microseconds())/1000.0, len(imgBytes)/1024)
	}
}

func TestRenderFitTextDownscaling(t *testing.T) {
	mgr, _ := templates.NewManager("")
	tmpl, _ := mgr.Get("dark_tech")

	// Extremely long title that requires multiple wrappings and downscaling
	veryLongTitle := "This is an extraordinarily long article title written specifically to stress-test the dynamic FitText typography downscaling algorithm and assert that it fits within designated canvas boundaries without overflowing or panicking"

	imgBytes, err := renderer.RenderCard(tmpl, map[string]string{
		"title": veryLongTitle,
	})
	if err != nil {
		t.Fatalf("FitText render failed: %v", err)
	}

	if len(imgBytes) == 0 {
		t.Fatal("empty image returned")
	}
}
