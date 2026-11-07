package icons

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
)

// DrawIcon renders a named vector icon at the given position and scale.
func DrawIcon(dc *gg.Context, name string, x, y, size float64, hexColor string) bool {
	dc.Push()
	defer dc.Pop()

	r, g, b, a := ParseHexColor(hexColor)
	dc.SetRGBA(r, g, b, a)

	scale := size / 24.0
	dc.Translate(x, y)
	dc.Scale(scale, scale)

	switch strings.ToLower(strings.TrimSpace(name)) {
	case "github":
		drawGitHub(dc)
	case "twitter", "x":
		drawX(dc)
	case "calendar":
		drawCalendar(dc)
	case "clock":
		drawClock(dc)
	case "star":
		drawStar(dc)
	case "tag":
		drawTag(dc)
	case "user":
		drawUser(dc)
	case "eye":
		drawEye(dc)
	case "check":
		drawCheck(dc)
	case "bolt":
		drawBolt(dc)
	default:
		return false
	}

	return true
}

func drawGitHub(dc *gg.Context) {
	dc.DrawCircle(12, 12, 10)
	dc.Fill()

	dc.Push()
	dc.SetRGBA(0, 0, 0, 0.9)
	dc.DrawCircle(8, 12, 2.5)
	dc.DrawCircle(16, 12, 2.5)
	dc.Fill()
	dc.Pop()
}

func drawX(dc *gg.Context) {
	dc.SetLineWidth(2.2)
	dc.MoveTo(4, 4)
	dc.LineTo(20, 20)
	dc.MoveTo(20, 4)
	dc.LineTo(4, 20)
	dc.Stroke()
}

func drawCalendar(dc *gg.Context) {
	dc.SetLineWidth(2.0)
	dc.DrawRoundedRectangle(3, 4, 18, 17, 3)
	dc.Stroke()

	dc.MoveTo(3, 9)
	dc.LineTo(21, 9)
	dc.Stroke()

	dc.MoveTo(8, 2)
	dc.LineTo(8, 5)
	dc.MoveTo(16, 2)
	dc.LineTo(16, 5)
	dc.Stroke()
}

func drawClock(dc *gg.Context) {
	dc.SetLineWidth(2.0)
	dc.DrawCircle(12, 12, 9)
	dc.Stroke()

	dc.MoveTo(12, 7)
	dc.LineTo(12, 12)
	dc.LineTo(16, 14)
	dc.Stroke()
}

func drawStar(dc *gg.Context) {
	dc.MoveTo(12, 2)
	dc.LineTo(15, 9)
	dc.LineTo(22, 9.5)
	dc.LineTo(17, 14.5)
	dc.LineTo(18.5, 21.5)
	dc.LineTo(12, 17.5)
	dc.LineTo(5.5, 21.5)
	dc.LineTo(7, 14.5)
	dc.LineTo(2, 9.5)
	dc.LineTo(9, 9)
	dc.ClosePath()
	dc.Fill()
}

func drawTag(dc *gg.Context) {
	dc.SetLineWidth(2.0)
	dc.MoveTo(12, 2)
	dc.LineTo(20, 2)
	dc.LineTo(20, 10)
	dc.LineTo(10, 20)
	dc.LineTo(2, 12)
	dc.ClosePath()
	dc.Stroke()

	dc.DrawCircle(16, 6, 1.5)
	dc.Fill()
}

func drawUser(dc *gg.Context) {
	dc.DrawCircle(12, 7, 4)
	dc.Fill()

	dc.SetLineWidth(2.2)
	dc.DrawEllipticalArc(12, 20, 8, 6, gg.Radians(180), gg.Radians(360))
	dc.Stroke()
}

func drawEye(dc *gg.Context) {
	dc.SetLineWidth(2.0)
	dc.MoveTo(2, 12)
	dc.QuadraticTo(12, 4, 22, 12)
	dc.QuadraticTo(12, 20, 2, 12)
	dc.Stroke()

	dc.DrawCircle(12, 12, 3)
	dc.Fill()
}

func drawCheck(dc *gg.Context) {
	dc.SetLineWidth(2.5)
	dc.MoveTo(4, 12)
	dc.LineTo(9, 17)
	dc.LineTo(20, 6)
	dc.Stroke()
}

func drawBolt(dc *gg.Context) {
	dc.MoveTo(13, 2)
	dc.LineTo(4, 13)
	dc.LineTo(11, 13)
	dc.LineTo(10, 22)
	dc.LineTo(20, 10)
	dc.LineTo(13, 10)
	dc.ClosePath()
	dc.Fill()
}

// ParseHexColor converts hex string (#RGB, #RRGGBB, #RRGGBBAA) to RGBA components in 0.0-1.0
func ParseHexColor(hex string) (r, g, b, a float64) {
	hex = strings.TrimPrefix(hex, "#")
	a = 1.0

	var c uint64
	switch len(hex) {
	case 3:
		// #RGB -> #RRGGBB
		expanded := fmt.Sprintf("%c%c%c%c%c%c", hex[0], hex[0], hex[1], hex[1], hex[2], hex[2])
		fmt.Sscanf(expanded, "%x", &c)
		r = float64((c>>16)&0xFF) / 255.0
		g = float64((c>>8)&0xFF) / 255.0
		b = float64(c&0xFF) / 255.0
	case 6:
		fmt.Sscanf(hex, "%x", &c)
		r = float64((c>>16)&0xFF) / 255.0
		g = float64((c>>8)&0xFF) / 255.0
		b = float64(c&0xFF) / 255.0
	case 8:
		fmt.Sscanf(hex, "%x", &c)
		r = float64((c>>24)&0xFF) / 255.0
		g = float64((c>>16)&0xFF) / 255.0
		b = float64((c>>8)&0xFF) / 255.0
		a = float64(c&0xFF) / 255.0
	default:
		return 1.0, 1.0, 1.0, 1.0
	}
	return
}
