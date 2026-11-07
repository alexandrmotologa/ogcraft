package templates

// GradientStop defines a color stop in a gradient.
type GradientStop struct {
	Offset float64 `json:"offset"`
	Color  string  `json:"color"`
}

// BackgroundSpec defines the canvas background style.
type BackgroundSpec struct {
	Type        string         `json:"type"`                   // "solid", "gradient", "radial"
	Color       string         `json:"color,omitempty"`        // For solid
	Direction   string         `json:"direction,omitempty"`    // "to-bottom", "to-right", "to-bottom-right", "to-bottom-left"
	Stops       []GradientStop `json:"stops,omitempty"`        // For linear gradient
	CenterX     float64        `json:"center_x,omitempty"`     // For radial
	CenterY     float64        `json:"center_y,omitempty"`     // For radial
	Radius      float64        `json:"radius,omitempty"`       // For radial
	CenterColor string         `json:"center_color,omitempty"` // For radial
	OuterColor  string         `json:"outer_color,omitempty"`  // For radial
}

// Element represents any visual item inside a template.
type Element struct {
	Type            string  `json:"type"`                        // "text", "badge", "avatar", "icon", "divider", "container"
	Field           string  `json:"field,omitempty"`            // Query parameter key (e.g. "title", "author")
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
	Width           float64 `json:"width,omitempty"`            // Width for divider / container
	Height          float64 `json:"height,omitempty"`           // Height for divider / container
	MaxWidth        float64 `json:"max_width,omitempty"`        // Max width for text wrapping
	LineHeight      float64 `json:"line_height,omitempty"`      // Line height multiplier
	MaxLines        int     `json:"max_lines,omitempty"`        // Max lines before truncation
	FitText         bool    `json:"fit_text,omitempty"`         // Automatically downscale font to fit
	PaddingX        float64 `json:"padding_x,omitempty"`        // Horizontal padding for badge
	PaddingY        float64 `json:"padding_y,omitempty"`        // Vertical padding for badge
	Icon            string  `json:"icon,omitempty"`             // Icon name (e.g. "github", "star")
	Size            float64 `json:"size,omitempty"`             // Size for avatar or icon
	Shape           string  `json:"shape,omitempty"`            // Shape for avatar ("circle", "rounded", "square")
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
