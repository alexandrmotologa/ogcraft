package renderer

import (
	"strings"

	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
)

// RenderCard renders a complete template with query parameter substitutions into PNG bytes.
func RenderCard(tmpl *templates.Template, params map[string]string) ([]byte, error) {
	width := tmpl.Width
	if width <= 0 {
		width = 1200
	}
	height := tmpl.Height
	if height <= 0 {
		height = 630
	}

	canvas := NewCanvas(width, height)
	dc := canvas.Context()

	// 1. Render Background
	ApplyBackground(dc, float64(width), float64(height), tmpl.Background)

	// 2. Render Elements
	for _, el := range tmpl.Elements {
		val := resolveFieldValue(el.Field, el.Default, params)

		switch strings.ToLower(el.Type) {
		case "text":
			textSpec := TextSpec{
				Text:       val,
				Font:       el.Font,
				FontSize:   el.FontSize,
				Color:      el.Color,
				X:          el.X,
				Y:          el.Y,
				MaxWidth:   el.MaxWidth,
				LineHeight: el.LineHeight,
				MaxLines:   el.MaxLines,
				FitText:    el.FitText,
			}
			if err := RenderText(dc, textSpec); err != nil {
				return nil, err
			}

		case "badge":
			badgeSpec := BadgeSpec{
				Text:            val,
				Icon:            el.Icon,
				Font:            el.Font,
				FontSize:        el.FontSize,
				Color:           el.Color,
				BackgroundColor: el.BackgroundColor,
				BorderColor:     el.BorderColor,
				BorderWidth:     el.BorderWidth,
				X:               el.X,
				Y:               el.Y,
				PaddingX:        el.PaddingX,
				PaddingY:        el.PaddingY,
				BorderRadius:    el.BorderRadius,
			}
			if err := RenderBadge(dc, badgeSpec); err != nil {
				return nil, err
			}

		case "avatar":
			authorName := resolveFieldValue("author", "", params)
			avatarSpec := AvatarSpec{
				URL:          val,
				FallbackText: authorName,
				X:            el.X,
				Y:            el.Y,
				Size:         el.Size,
				Shape:        el.Shape,
				BorderColor:  el.BorderColor,
				BorderWidth:  el.BorderWidth,
			}
			if err := RenderAvatar(dc, avatarSpec); err != nil {
				return nil, err
			}

		case "icon":
			icons.DrawIcon(dc, el.Icon, el.X, el.Y, el.Size, el.Color)

		case "divider":
			dividerSpec := DividerSpec{
				X:      el.X,
				Y:      el.Y,
				Width:  el.Width,
				Height: el.Height,
				Color:  el.Color,
			}
			RenderDivider(dc, dividerSpec)

		case "container":
			containerSpec := ContainerSpec{
				X:               el.X,
				Y:               el.Y,
				Width:           el.Width,
				Height:          el.Height,
				BorderRadius:    el.BorderRadius,
				BackgroundColor: el.BackgroundColor,
				BorderColor:     el.BorderColor,
				BorderWidth:     el.BorderWidth,
			}
			RenderContainer(dc, containerSpec)
		}
	}

	return canvas.EncodePNG()
}

func resolveFieldValue(field, defaultValue string, params map[string]string) string {
	if field != "" && params != nil {
		if val, exists := params[field]; exists {
			return val
		}
	}
	return defaultValue
}
