package renderer

import (
	"image/color"
	"strings"

	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
	"github.com/fogleman/gg"
)

// ApplyBackground paints the background and optional pattern overlay onto the canvas context.
func ApplyBackground(dc *gg.Context, width, height float64, bg templates.BackgroundSpec) {
	dc.Push()
	defer dc.Pop()

	switch strings.ToLower(bg.Type) {
	case "gradient", "linear":
		applyLinearGradient(dc, width, height, bg)
	case "radial":
		applyRadialGradient(dc, width, height, bg)
	case "solid":
		fallthrough
	default:
		r, g, b, a := icons.ParseHexColor(bg.Color)
		if bg.Color == "" {
			r, g, b, a = 0.05, 0.07, 0.1, 1.0 // default dark slate
		}
		dc.SetRGBA(r, g, b, a)
		dc.Clear()
	}

	// Apply pattern overlay if configured
	if bg.Pattern != "" && bg.Pattern != "none" {
		applyPatternOverlay(dc, width, height, bg)
	}
}

func applyPatternOverlay(dc *gg.Context, width, height float64, bg templates.BackgroundSpec) {
	dc.Push()
	defer dc.Pop()

	pColor := bg.PatternColor
	if pColor == "" {
		pColor = "#ffffff0f" // very subtle 6% white
	}
	r, g, b, a := icons.ParseHexColor(pColor)
	dc.SetRGBA(r, g, b, a)

	spacing := bg.PatternSpacing
	if spacing <= 0 {
		spacing = 30.0
	}

	switch strings.ToLower(bg.Pattern) {
	case "dots":
		dotRadius := 1.5
		for x := spacing / 2; x < width; x += spacing {
			for y := spacing / 2; y < height; y += spacing {
				dc.DrawCircle(x, y, dotRadius)
			}
		}
		dc.Fill()

	case "grid":
		dc.SetLineWidth(1.0)
		for x := 0.0; x <= width; x += spacing {
			dc.MoveTo(x, 0)
			dc.LineTo(x, height)
		}
		for y := 0.0; y <= height; y += spacing {
			dc.MoveTo(0, y)
			dc.LineTo(width, y)
		}
		dc.Stroke()

	case "lines":
		dc.SetLineWidth(1.0)
		for d := -height; d <= width+height; d += spacing {
			dc.MoveTo(d, 0)
			dc.LineTo(d+height, height)
		}
		dc.Stroke()
	}
}

func applyLinearGradient(dc *gg.Context, width, height float64, bg templates.BackgroundSpec) {
	var x0, y0, x1, y1 float64

	switch bg.Direction {
	case "to-right":
		x0, y0, x1, y1 = 0, height/2, width, height/2
	case "to-left":
		x0, y0, x1, y1 = width, height/2, 0, height/2
	case "to-bottom":
		x0, y0, x1, y1 = width/2, 0, width/2, height
	case "to-top":
		x0, y0, x1, y1 = width/2, height, width/2, 0
	case "to-bottom-left":
		x0, y0, x1, y1 = width, 0, 0, height
	case "to-bottom-right":
		fallthrough
	default:
		x0, y0, x1, y1 = 0, 0, width, height
	}

	grad := gg.NewLinearGradient(x0, y0, x1, y1)

	if len(bg.Stops) == 0 {
		grad.AddColorStop(0.0, toRGBAColor(0.05, 0.07, 0.09, 1.0))
		grad.AddColorStop(1.0, toRGBAColor(0.09, 0.11, 0.14, 1.0))
	} else {
		for _, stop := range bg.Stops {
			r, g, b, a := icons.ParseHexColor(stop.Color)
			grad.AddColorStop(stop.Offset, toRGBAColor(r, g, b, a))
		}
	}

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()
}

func applyRadialGradient(dc *gg.Context, width, height float64, bg templates.BackgroundSpec) {
	cx := bg.CenterX
	if cx == 0 {
		cx = width * 0.8
	}
	cy := bg.CenterY
	if cy == 0 {
		cy = height * 0.2
	}
	r0 := 0.0
	r1 := bg.Radius
	if r1 == 0 {
		r1 = 500.0
	}

	grad := gg.NewRadialGradient(cx, cy, r0, cx, cy, r1)

	cr, cg, cb, ca := icons.ParseHexColor(bg.CenterColor)
	if bg.CenterColor == "" {
		cr, cg, cb, ca = 0.2, 0.4, 0.9, 0.35 // soft blue glow
	}

	or, og, ob, oa := icons.ParseHexColor(bg.OuterColor)
	if bg.OuterColor == "" {
		or, og, ob, oa = 0.04, 0.05, 0.07, 1.0
	}

	grad.AddColorStop(0.0, toRGBAColor(cr, cg, cb, ca))
	grad.AddColorStop(1.0, toRGBAColor(or, og, ob, oa))

	// Base fill
	dc.SetRGBA(or, og, ob, oa)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()

	// Radial blend
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()
}

func toRGBAColor(r, g, b, a float64) color.Color {
	return color.NRGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}
}
