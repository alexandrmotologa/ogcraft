package renderer_test

import (
	"bytes"
	"image/jpeg"
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

func TestRenderMultiFormatWebPAndJPEG(t *testing.T) {
	mgr, _ := templates.NewManager("")
	tmpl, _ := mgr.Get("dark_tech")

	params := map[string]string{"title": "Testing Multi-Format Encoding"}

	// 1. Test WebP
	webpBytes, mime, err := renderer.RenderCardFormat(tmpl, params, "webp")
	if err != nil {
		t.Fatalf("webp render failed: %v", err)
	}
	if mime != "image/webp" {
		t.Errorf("expected mime image/webp, got %s", mime)
	}
	if len(webpBytes) < 1000 {
		t.Errorf("webp output unexpectedly small: %d bytes", len(webpBytes))
	}
	t.Logf("Rendered WebP size: %d KB", len(webpBytes)/1024)

	// 2. Test JPEG
	jpegBytes, mime, err := renderer.RenderCardFormat(tmpl, params, "jpeg")
	if err != nil {
		t.Fatalf("jpeg render failed: %v", err)
	}
	if mime != "image/jpeg" {
		t.Errorf("expected mime image/jpeg, got %s", mime)
	}
	cfg, err := jpeg.DecodeConfig(bytes.NewReader(jpegBytes))
	if err != nil {
		t.Fatalf("invalid JPEG generated: %v", err)
	}
	if cfg.Width != 1200 || cfg.Height != 630 {
		t.Errorf("wrong JPEG dimensions: %dx%d", cfg.Width, cfg.Height)
	}
	t.Logf("Rendered JPEG size: %d KB", len(jpegBytes)/1024)
}

func TestRenderAspectRatiosAndCustomDimensions(t *testing.T) {
	mgr, _ := templates.NewManager("")
	tmpl, _ := mgr.Get("minimal_light")

	// Square 800x800
	squareBytes, err := renderer.RenderCard(tmpl, map[string]string{"aspect": "square"})
	if err != nil {
		t.Fatalf("square render failed: %v", err)
	}
	cfg, _ := png.DecodeConfig(bytes.NewReader(squareBytes))
	if cfg.Width != 800 || cfg.Height != 800 {
		t.Errorf("expected 800x800 for square, got %dx%d", cfg.Width, cfg.Height)
	}

	// Custom dimensions 600x400
	customBytes, err := renderer.RenderCard(tmpl, map[string]string{"w": "600", "h": "400"})
	if err != nil {
		t.Fatalf("custom dimensions render failed: %v", err)
	}
	cfg2, _ := png.DecodeConfig(bytes.NewReader(customBytes))
	if cfg2.Width != 600 || cfg2.Height != 400 {
		t.Errorf("expected 600x400, got %dx%d", cfg2.Width, cfg2.Height)
	}
}

func TestRenderCodeSnippetAndMetricElements(t *testing.T) {
	customTmpl := &templates.Template{
		ID:     "code_demo",
		Name:   "Code Demo",
		Width:  1200,
		Height: 630,
		Background: templates.BackgroundSpec{
			Type:    "solid",
			Color:   "#090d16",
			Pattern: "dots",
		},
		Elements: []templates.Element{
			{
				Type:         "code",
				Default:      "package main\n\nfunc main() {\n    println(\"Hello ogcraft!\")\n}",
				Language:     "go",
				X:            60,
				Y:            60,
				Width:        650,
				Height:       380,
				ShowControls: true,
				LineNumbers:  true,
			},
			{
				Type:            "metric",
				Value:           "99.9%",
				Subtext:         "Service Uptime",
				Icon:            "bolt",
				X:               760,
				Y:               60,
				Width:           380,
				Height:          140,
				BackgroundColor: "#161e2e",
				Color:           "#38bdf8",
				BorderColor:     "#223048",
				BorderWidth:     1,
				Progress:        0.999,
			},
		},
	}

	imgBytes, err := renderer.RenderCard(customTmpl, nil)
	if err != nil {
		t.Fatalf("failed rendering code and metric elements: %v", err)
	}
	if len(imgBytes) == 0 {
		t.Fatal("empty image generated")
	}
}

func TestRenderFitTextDownscaling(t *testing.T) {
	mgr, _ := templates.NewManager("")
	tmpl, _ := mgr.Get("dark_tech")

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
