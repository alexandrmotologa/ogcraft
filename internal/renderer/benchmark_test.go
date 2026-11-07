package renderer_test

import (
	"testing"

	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/renderer"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
)

func BenchmarkRenderDarkTechWarm(b *testing.B) {
	_ = fonts.Init()
	mgr, _ := templates.NewManager("")
	tmpl, _ := mgr.Get("dark_tech")

	// Override avatar to empty so no HTTP network request happens during benchmark
	params := map[string]string{
		"title":  "Sub-15ms Social Card Rendering in Go",
		"author": "Alexandr Motologa",
		"tag":    "BENCHMARK",
		"avatar": "",
	}

	// Warmup
	_, _ = renderer.RenderCard(tmpl, params)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderCard(tmpl, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
