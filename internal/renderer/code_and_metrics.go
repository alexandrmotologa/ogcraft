package renderer

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/fogleman/gg"
)

// CodeSnippetSpec defines parameters for rendering a terminal code window.
type CodeSnippetSpec struct {
	Code         string
	Language     string
	X            float64
	Y            float64
	Width        float64
	Height       float64
	FontSize     float64
	ShowControls bool
	LineNumbers  bool
}

// RenderCodeSnippet renders a syntax-highlighted code block with a terminal window frame.
func RenderCodeSnippet(dc *gg.Context, spec CodeSnippetSpec) error {
	dc.Push()
	defer dc.Pop()

	w := spec.Width
	if w <= 0 {
		w = 600
	}
	h := spec.Height
	if h <= 0 {
		h = 320
	}
	fontSize := spec.FontSize
	if fontSize <= 0 {
		fontSize = 20
	}

	// 1. Draw Window Frame Background
	dc.SetRGBA(0.05, 0.07, 0.1, 0.95)
	dc.DrawRoundedRectangle(spec.X, spec.Y, w, h, 12)
	dc.Fill()

	// Window Border
	dc.SetRGBA(0.2, 0.25, 0.35, 0.8)
	dc.SetLineWidth(1.5)
	dc.DrawRoundedRectangle(spec.X, spec.Y, w, h, 12)
	dc.Stroke()

	headerH := 40.0
	currY := spec.Y + 20.0

	// 2. Draw Header Window Controls (macOS Red, Yellow, Green dots)
	if spec.ShowControls {
		dotRadius := 6.0
		// Red
		dc.SetRGB(1.0, 0.37, 0.34)
		dc.DrawCircle(spec.X+20, currY, dotRadius)
		dc.Fill()
		// Yellow
		dc.SetRGB(1.0, 0.74, 0.18)
		dc.DrawCircle(spec.X+40, currY, dotRadius)
		dc.Fill()
		// Green
		dc.SetRGB(0.15, 0.79, 0.25)
		dc.DrawCircle(spec.X+60, currY, dotRadius)
		dc.Fill()

		// Optional Title in header center
		if spec.Language != "" {
			face, err := fonts.GetFontFace(fonts.JetBrainsMonoMedium, 14)
			if err == nil {
				dc.SetFontFace(face)
				dc.SetRGBA(0.55, 0.6, 0.68, 0.9)
				title := spec.Language
				tW, _ := dc.MeasureString(title)
				dc.DrawString(title, spec.X+(w-tW)/2, currY+5)
			}
		}

		// Divider under header
		dc.SetRGBA(0.2, 0.25, 0.35, 0.5)
		dc.SetLineWidth(1.0)
		dc.DrawLine(spec.X, spec.Y+headerH, spec.X+w, spec.Y+headerH)
		dc.Stroke()
	}

	// 3. Render Code Lines
	lines := strings.Split(spec.Code, "\n")
	face, err := fonts.GetFontFace(fonts.JetBrainsMonoMedium, fontSize)
	if err != nil {
		return err
	}
	dc.SetFontFace(face)

	lineHeight := fontSize * 1.5
	codeStartY := spec.Y + headerH + fontSize + 10

	gutterW := 0.0
	if spec.LineNumbers {
		gutterW = 45.0
	}

	for i, line := range lines {
		if codeStartY > spec.Y+h-15 {
			break
		}

		// Line Number
		if spec.LineNumbers {
			dc.SetRGBA(0.4, 0.45, 0.52, 0.6)
			lineNumStr := fmt.Sprintf("%2d", i+1)
			dc.DrawString(lineNumStr, spec.X+16, codeStartY)
		}

		// Tokenize and render code line
		renderSyntaxHighlightedLine(dc, line, spec.X+16+gutterW, codeStartY)

		codeStartY += lineHeight
	}

	return nil
}

// renderSyntaxHighlightedLine renders keywords, strings, comments, and symbols with syntax colors.
func renderSyntaxHighlightedLine(dc *gg.Context, line string, startX, startY float64) {
	currX := startX

	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
		// Comment line
		dc.SetRGBA(0.55, 0.58, 0.62, 0.9) // gray
		dc.DrawString(line, currX, startY)
		return
	}

	words := splitTokens(line)
	for _, tok := range words {
		// Set color based on token type
		switch {
		case isKeyword(tok):
			dc.SetRGBA(1.0, 0.48, 0.45, 1.0) // salmon / pink
		case isType(tok):
			dc.SetRGBA(0.47, 0.75, 1.0, 1.0) // light blue
		case strings.HasPrefix(tok, "\"") || strings.HasPrefix(tok, "'"):
			dc.SetRGBA(0.65, 0.84, 1.0, 1.0) // cyan string
		case isOperator(tok):
			dc.SetRGBA(0.82, 0.66, 1.0, 1.0) // purple
		case isNumber(tok):
			dc.SetRGBA(0.98, 0.8, 0.08, 1.0) // yellow
		default:
			dc.SetRGBA(0.9, 0.94, 0.98, 1.0) // light text
		}

		dc.DrawString(tok, currX, startY)
		tokW, _ := dc.MeasureString(tok)
		currX += tokW
	}
}

func splitTokens(s string) []string {
	var tokens []string
	var current strings.Builder

	inString := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if (ch == '"' || ch == '\'') && !inString {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			inString = true
			quoteChar = ch
			current.WriteByte(ch)
			continue
		}

		if inString {
			current.WriteByte(ch)
			if ch == quoteChar {
				tokens = append(tokens, current.String())
				current.Reset()
				inString = false
			}
			continue
		}

		if isSeparator(ch) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func isSeparator(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '(' || ch == ')' || ch == '{' ||
		ch == '}' || ch == '[' || ch == ']' || ch == ',' || ch == '.' ||
		ch == ';' || ch == ':' || ch == '+' || ch == '=' || ch == '-' ||
		ch == '*' || ch == '/' || ch == '&' || ch == '|' || ch == '!'
}

func isKeyword(tok string) bool {
	keywords := map[string]bool{
		"package": true, "import": true, "func": true, "return": true,
		"var": true, "const": true, "type": true, "struct": true,
		"interface": true, "if": true, "else": true, "for": true,
		"range": true, "go": true, "select": true, "case": true,
		"defer": true, "switch": true, "default": true,
		"let": true, "export": true, "async": true, "await": true,
		"from": true, "function": true, "class": true, "def": true,
	}
	return keywords[tok]
}

func isType(tok string) bool {
	types := map[string]bool{
		"string": true, "int": true, "int64": true, "float64": true,
		"bool": true, "error": true, "byte": true, "rune": true,
		"any": true, "nil": true, "true": true, "false": true,
		"void": true, "number": true, "boolean": true, "Promise": true,
	}
	return types[tok]
}

func isOperator(tok string) bool {
	ops := map[string]bool{
		":=": true, "==": true, "!=": true, "&&": true, "||": true,
		"->": true, "=>": true, "+=": true, "-=": true,
	}
	return ops[tok]
}

func isNumber(tok string) bool {
	if len(tok) == 0 {
		return false
	}
	for i := 0; i < len(tok); i++ {
		if (tok[i] < '0' || tok[i] > '9') && tok[i] != '.' {
			return false
		}
	}
	return true
}

// MetricSpec defines parameters for rendering a big stat or rating card.
type MetricSpec struct {
	Value           string
	Subtext         string
	Icon            string
	X               float64
	Y               float64
	Width           float64
	Height          float64
	ValueColor      string
	SubtextColor    string
	BackgroundColor string
	BorderColor     string
	BorderWidth     float64
	BorderRadius    float64
	Progress        float64 // 0.0 to 1.0
}

// RenderMetricCard renders a stat callout with big value, subtext, and optional progress bar.
func RenderMetricCard(dc *gg.Context, spec MetricSpec) error {
	dc.Push()
	defer dc.Pop()

	w := spec.Width
	if w <= 0 {
		w = 280
	}
	h := spec.Height
	if h <= 0 {
		h = 120
	}
	radius := spec.BorderRadius
	if radius <= 0 {
		radius = 12
	}

	// 1. Draw Card Background
	if spec.BackgroundColor != "" {
		bgR, bgG, bgB, bgA := icons.ParseHexColor(spec.BackgroundColor)
		dc.SetRGBA(bgR, bgG, bgB, bgA)
		dc.DrawRoundedRectangle(spec.X, spec.Y, w, h, radius)
		dc.Fill()
	}

	// 2. Draw Border
	if spec.BorderColor != "" && spec.BorderWidth > 0 {
		bR, bG, bB, bA := icons.ParseHexColor(spec.BorderColor)
		dc.SetRGBA(bR, bG, bB, bA)
		dc.SetLineWidth(spec.BorderWidth)
		dc.DrawRoundedRectangle(spec.X, spec.Y, w, h, radius)
		dc.Stroke()
	}

	contentX := spec.X + 24
	contentY := spec.Y + 48

	// Optional Icon
	if spec.Icon != "" {
		iconSize := 28.0
		icons.DrawIcon(dc, spec.Icon, contentX, spec.Y+24, iconSize, spec.ValueColor)
		contentX += iconSize + 12
	}

	// 3. Render Big Stat Value
	valColor := spec.ValueColor
	if valColor == "" {
		valColor = "#ffffff"
	}
	vR, vG, vB, vA := icons.ParseHexColor(valColor)
	dc.SetRGBA(vR, vG, vB, vA)

	valFace, err := fonts.GetFontFace(fonts.InterBold, 44)
	if err == nil {
		dc.SetFontFace(valFace)
		dc.DrawString(spec.Value, contentX, contentY)
	}

	// 4. Render Subtext / Label
	subColor := spec.SubtextColor
	if subColor == "" {
		subColor = "#94a3b8"
	}
	sR, sG, sB, sA := icons.ParseHexColor(subColor)
	dc.SetRGBA(sR, sG, sB, sA)

	subFace, err := fonts.GetFontFace(fonts.InterSemiBold, 18)
	if err == nil {
		dc.SetFontFace(subFace)
		dc.DrawString(spec.Subtext, spec.X+24, contentY+32)
	}

	// 5. Render Optional Progress Bar
	if spec.Progress > 0.0 {
		barX := spec.X + 24
		barY := spec.Y + h - 18
		barW := w - 48
		barH := 6.0

		// Background Track
		dc.SetRGBA(0.2, 0.25, 0.35, 0.6)
		dc.DrawRoundedRectangle(barX, barY, barW, barH, 3)
		dc.Fill()

		// Fill
		fillW := barW * spec.Progress
		if fillW > barW {
			fillW = barW
		}
		dc.SetRGBA(vR, vG, vB, vA)
		dc.DrawRoundedRectangle(barX, barY, fillW, barH, 3)
		dc.Fill()
	}

	return nil
}
