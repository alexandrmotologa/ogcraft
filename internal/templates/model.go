package templates

// GradientStop defines a color stop in a gradient.
type GradientStop struct {
	Offset float64 `json:"offset"`
	Color  string  `json:"color"`
}

// BackgroundSpec defines the canvas background style.
type BackgroundSpec struct {
	Type           string         `json:"type"`                      // "solid", "gradient", "radial"
	Color          string         `json:"color,omitempty"`           // For solid
	Direction      string         `json:"direction,omitempty"`       // "to-bottom", "to-right", "to-bottom-right", "to-bottom-left"
	Stops          []GradientStop `json:"stops,omitempty"`           // For linear gradient
	CenterX        float64        `json:"center_x,omitempty"`        // For radial
	CenterY        float64        `json:"center_y,omitempty"`        // For radial
	Radius         float64        `json:"radius,omitempty"`          // For radial
	CenterColor    string         `json:"center_color,omitempty"`    // For radial
	OuterColor     string         `json:"outer_color,omitempty"`     // For radial
	Pattern        string         `json:"pattern,omitempty"`         // "none", "dots", "grid", "lines"
	PatternColor   string         `json:"pattern_color,omitempty"`   // Hex color for pattern overlay
	PatternSpacing float64        `json:"pattern_spacing,omitempty"` // Distance between dots/lines in pixels
}

// Element represents any visual item inside a template.
type Element struct {
	Type            string  `json:"type"`                        // "text", "badge", "avatar", "icon", "divider", "container", "code", "metric"
	Field           string  `json:"field,omitempty"`            // Query parameter key (e.g. "title", "author", "code", "metric")
	Default         string  `json:"default,omitempty"`          // Fallback value if query parameter is empty
	Font            string  `json:"font,omitempty"`             // Font family name
	FontSize        float64 `json:"font_size,omitempty"`        // Font size in points
	Color           string  `json:"color,omitempty"`            // Hex color
	BackgroundColor string  `json:"background_color,omitempty"` // Hex background color
	BorderColor     string  `json:"border_color,omitempty"`     // Hex border color
	BorderWidth     float64 `json:"border_width,omitempty"`     // Border width
	BorderRadius    float64 `json:"border_radius,omitempty"`    // Corner radius
	X               float64 `json:"x"`                          // X coordinate
	Y               float64 `json:"y"`                          // Y coordinate
	Width           float64 `json:"width,omitempty"`            // Width for divider / container / code / metric
	Height          float64 `json:"height,omitempty"`           // Height for divider / container / code / metric
	MaxWidth        float64 `json:"max_width,omitempty"`        // Max width for text wrapping
	LineHeight      float64 `json:"line_height,omitempty"`      // Line height multiplier
	MaxLines        int     `json:"max_lines,omitempty"`        // Max lines before truncation
	FitText         bool    `json:"fit_text,omitempty"`         // Automatically downscale font to fit
	PaddingX        float64 `json:"padding_x,omitempty"`        // Horizontal padding for badge
	PaddingY        float64 `json:"padding_y,omitempty"`        // Vertical padding for badge
	Icon            string  `json:"icon,omitempty"`             // Icon name (e.g. "github", "star")
	Size            float64 `json:"size,omitempty"`             // Size for avatar or icon
	Shape           string  `json:"shape,omitempty"`            // Shape for avatar ("circle", "rounded", "square")

	// Code snippet specific properties
	Language     string `json:"language,omitempty"`      // "go", "typescript", "bash", "python"
	ShowControls bool   `json:"show_controls,omitempty"` // Show macOS red/yellow/green window dots
	LineNumbers  bool   `json:"line_numbers,omitempty"`  // Show line numbers in gutter

	// Metric specific properties
	Value    string  `json:"value,omitempty"`    // Big stat value (e.g. "99.9%", "12.4k")
	Subtext  string  `json:"subtext,omitempty"`  // Label below value (e.g. "Uptime", "GitHub Stars")
	Progress float64 `json:"progress,omitempty"` // Optional progress bar (0.0 to 1.0)
}

// Template defines the structure of an OpenGraph card canvas.
type Template struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	Background  BackgroundSpec `json:"background"`
	Elements    []Element      `json:"elements"`
	Builtin     bool           `json:"builtin,omitempty"`
}

// TemplateSummary provides metadata for template listing.
type TemplateSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Builtin       bool   `json:"builtin"`
	ElementsCount int    `json:"elements_count"`
}
