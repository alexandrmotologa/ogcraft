package renderer

import (
	"strconv"
	"strings"

	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
)

// RenderCard renders a complete template into PNG bytes (default format).
func RenderCard(tmpl *templates.Template, params map[string]string) ([]byte, error) {
	data, _, err := RenderCardFormat(tmpl, params, "png")
	return data, err
}

// RenderCardFormat renders a complete template into the specified format (png, webp, jpeg).
func RenderCardFormat(tmpl *templates.Template, params map[string]string, format string) ([]byte, string, error) {
	width := tmpl.Width
	if width <= 0 {
		width = 1200
	}
	height := tmpl.Height
	if height <= 0 {
		height = 630
	}

	// 1. Check for dynamic aspect ratio or custom dimension overrides from query params
	if params != nil {
		if aspect := params["aspect"]; aspect != "" {
			switch strings.ToLower(aspect) {
			case "square", "1:1":
				width, height = 800, 800
			case "story", "vertical", "9:16":
				width, height = 1080, 1920
			case "og", "1.91:1":
				width, height = 1200, 630
			}
		}
		if customW := params["w"]; customW != "" {
			if w, err := strconv.Atoi(customW); err == nil && w >= 200 && w <= 3840 {
				width = w
			}
		}
		if customH := params["h"]; customH != "" {
			if h, err := strconv.Atoi(customH); err == nil && h >= 200 && h <= 3840 {
				height = h
			}
		}
	}

	canvas := NewCanvas(width, height)
	dc := canvas.Context()

	// 2. Clone background and apply dynamic palette overrides
	bgSpec := tmpl.Background
	applyDynamicTheme(&bgSpec, params)

	// 3. Render Background & Patterns
	ApplyBackground(dc, float64(width), float64(height), bgSpec)

	// Accent and text color overrides
	accentColor := ""
	textColor := ""
	if params != nil {
		accentColor = params["accent"]
		textColor = params["text_color"]
	}

	// 4. Render Elements
	for _, el := range tmpl.Elements {
		val := resolveFieldValue(el.Field, el.Default, params)

		elColor := el.Color
		if textColor != "" && (el.Type == "text" || el.Type == "badge") {
			elColor = textColor
		}

		elBgColor := el.BackgroundColor
		if accentColor != "" && (el.Type == "badge" || el.Type == "metric") {
			elBgColor = accentColor
		}

		switch strings.ToLower(el.Type) {
		case "text":
			textSpec := TextSpec{
				Text:       val,
				Font:       el.Font,
				FontSize:   el.FontSize,
				Color:      elColor,
				X:          el.X,
				Y:          el.Y,
				MaxWidth:   el.MaxWidth,
				LineHeight: el.LineHeight,
				MaxLines:   el.MaxLines,
				FitText:    el.FitText,
			}
			if err := RenderText(dc, textSpec); err != nil {
				return nil, "", err
			}

		case "badge":
			badgeSpec := BadgeSpec{
				Text:            val,
				Icon:            el.Icon,
				Font:            el.Font,
				FontSize:        el.FontSize,
				Color:           elColor,
				BackgroundColor: elBgColor,
				BorderColor:     el.BorderColor,
				BorderWidth:     el.BorderWidth,
				X:               el.X,
				Y:               el.Y,
				PaddingX:        el.PaddingX,
				PaddingY:        el.PaddingY,
				BorderRadius:    el.BorderRadius,
			}
			if err := RenderBadge(dc, badgeSpec); err != nil {
				return nil, "", err
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
				return nil, "", err
			}

		case "icon":
			iconColor := elColor
			if accentColor != "" {
				iconColor = accentColor
			}
			icons.DrawIcon(dc, el.Icon, el.X, el.Y, el.Size, iconColor)

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

		case "code":
			codeContent := val
			if codeContent == "" {
				codeContent = el.Default
			}
			codeSpec := CodeSnippetSpec{
				Code:         codeContent,
				Language:     el.Language,
				X:            el.X,
				Y:            el.Y,
				Width:        el.Width,
				Height:       el.Height,
				FontSize:     el.FontSize,
				ShowControls: el.ShowControls,
				LineNumbers:  el.LineNumbers,
			}
			if err := RenderCodeSnippet(dc, codeSpec); err != nil {
				return nil, "", err
			}

		case "metric":
			metricVal := val
			if metricVal == "" {
				metricVal = el.Value
			}
			subtext := el.Subtext
			if s := resolveFieldValue("metric_subtext", "", params); s != "" {
				subtext = s
			}
			metricSpec := MetricSpec{
				Value:           metricVal,
				Subtext:         subtext,
				Icon:            el.Icon,
				X:               el.X,
				Y:               el.Y,
				Width:           el.Width,
				Height:          el.Height,
				ValueColor:      elColor,
				SubtextColor:    el.Color,
				BackgroundColor: elBgColor,
				BorderColor:     el.BorderColor,
				BorderWidth:     el.BorderWidth,
				BorderRadius:    el.BorderRadius,
				Progress:        el.Progress,
			}
			if err := RenderMetricCard(dc, metricSpec); err != nil {
				return nil, "", err
			}
		}
	}

	return canvas.EncodeFormat(format)
}

func applyDynamicTheme(bg *templates.BackgroundSpec, params map[string]string) {
	if params == nil {
		return
	}

	if theme := params["theme"]; theme == "light" {
		bg.Type = "solid"
		bg.Color = "#f8fafc"
	} else if theme == "dark" {
		bg.Type = "gradient"
		bg.Direction = "to-bottom-right"
		bg.Stops = []templates.GradientStop{
			{Offset: 0.0, Color: "#090d16"},
			{Offset: 1.0, Color: "#121a29"},
		}
	}

	if bgOverride := params["bg"]; bgOverride != "" {
		if !strings.HasPrefix(bgOverride, "#") {
			bgOverride = "#" + bgOverride
		}
		bg.Type = "solid"
		bg.Color = bgOverride
	}

	if pattern := params["pattern"]; pattern != "" {
		bg.Pattern = pattern
	}
}

func resolveFieldValue(field, defaultValue string, params map[string]string) string {
	if field != "" && params != nil {
		if val, exists := params[field]; exists {
			return val
		}
	}
	return defaultValue
}
