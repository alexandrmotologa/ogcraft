package renderer

import (
	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/fogleman/gg"
)

// BadgeSpec specifies badge / pill rendering.
type BadgeSpec struct {
	Text            string  `json:"text"`
	Icon            string  `json:"icon,omitempty"`
	Font            string  `json:"font"`
	FontSize        float64 `json:"font_size"`
	Color           string  `json:"color"`
	BackgroundColor string  `json:"background_color"`
	BorderColor     string  `json:"border_color,omitempty"`
	BorderWidth     float64 `json:"border_width,omitempty"`
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
	PaddingX        float64 `json:"padding_x"`
	PaddingY        float64 `json:"padding_y"`
	BorderRadius    float64 `json:"border_radius"`
}

// RenderBadge renders a badge pill with text and optional icon.
func RenderBadge(dc *gg.Context, spec BadgeSpec) error {
	dc.Push()
	defer dc.Pop()

	fontSize := spec.FontSize
	if fontSize <= 0 {
		fontSize = 20
	}
	padX := spec.PaddingX
	if padX <= 0 {
		padX = 14
	}
	padY := spec.PaddingY
	if padY <= 0 {
		padY = 7
	}
	radius := spec.BorderRadius
	if radius <= 0 {
		radius = 6
	}

	face, err := fonts.GetFontFace(spec.Font, fontSize)
	if err != nil {
		return err
	}
	dc.SetFontFace(face)

	textW, textH := dc.MeasureString(spec.Text)

	iconSize := fontSize * 0.95
	iconSpacing := 8.0
	extraIconW := 0.0
	if spec.Icon != "" {
		extraIconW = iconSize + iconSpacing
	}

	totalW := textW + extraIconW + (padX * 2)
	totalH := fontSize + (padY * 2)

	// Background fill
	bgR, bgG, bgB, bgA := icons.ParseHexColor(spec.BackgroundColor)
	dc.SetRGBA(bgR, bgG, bgB, bgA)
	dc.DrawRoundedRectangle(spec.X, spec.Y, totalW, totalH, radius)
	dc.Fill()

	// Optional border
	if spec.BorderColor != "" && spec.BorderWidth > 0 {
		bColorR, bColorG, bColorB, bColorA := icons.ParseHexColor(spec.BorderColor)
		dc.SetRGBA(bColorR, bColorG, bColorB, bColorA)
		dc.SetLineWidth(spec.BorderWidth)
		dc.DrawRoundedRectangle(spec.X, spec.Y, totalW, totalH, radius)
		dc.Stroke()
	}

	currX := spec.X + padX

	// Draw Icon if specified
	if spec.Icon != "" {
		iconY := spec.Y + (totalH-iconSize)/2
		icons.DrawIcon(dc, spec.Icon, currX, iconY, iconSize, spec.Color)
		currX += iconSize + iconSpacing
	}

	// Text draw
	txtR, txtG, txtB, txtA := icons.ParseHexColor(spec.Color)
	dc.SetRGBA(txtR, txtG, txtB, txtA)
	textBaselineY := spec.Y + padY + textH*0.82
	dc.DrawString(spec.Text, currX, textBaselineY)

	return nil
}

// DividerSpec specifies a separator line.
type DividerSpec struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Color  string  `json:"color"`
}

// RenderDivider renders a divider line.
func RenderDivider(dc *gg.Context, spec DividerSpec) {
	dc.Push()
	defer dc.Pop()

	r, g, b, a := icons.ParseHexColor(spec.Color)
	dc.SetRGBA(r, g, b, a)

	h := spec.Height
	if h <= 0 {
		h = 1
	}
	w := spec.Width
	if w <= 0 {
		w = float64(dc.Width()) - spec.X*2
	}

	dc.DrawRectangle(spec.X, spec.Y, w, h)
	dc.Fill()
}

// ContainerSpec specifies a card or frame background container.
type ContainerSpec struct {
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	BorderRadius    float64 `json:"border_radius"`
	BackgroundColor string  `json:"background_color"`
	BorderColor     string  `json:"border_color,omitempty"`
	BorderWidth     float64 `json:"border_width,omitempty"`
}

// RenderContainer renders a rounded card frame.
func RenderContainer(dc *gg.Context, spec ContainerSpec) {
	dc.Push()
	defer dc.Pop()

	bgR, bgG, bgB, bgA := icons.ParseHexColor(spec.BackgroundColor)
	dc.SetRGBA(bgR, bgG, bgB, bgA)
	dc.DrawRoundedRectangle(spec.X, spec.Y, spec.Width, spec.Height, spec.BorderRadius)
	dc.Fill()

	if spec.BorderColor != "" && spec.BorderWidth > 0 {
		bColorR, bColorG, bColorB, bColorA := icons.ParseHexColor(spec.BorderColor)
		dc.SetRGBA(bColorR, bColorG, bColorB, bColorA)
		dc.SetLineWidth(spec.BorderWidth)
		dc.DrawRoundedRectangle(spec.X, spec.Y, spec.Width, spec.Height, spec.BorderRadius)
		dc.Stroke()
	}
}
