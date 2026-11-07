package renderer

import (
	"strings"

	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// TextSpec specifies text rendering parameters.
type TextSpec struct {
	Text       string  `json:"text"`
	Font       string  `json:"font"`
	FontSize   float64 `json:"font_size"`
	Color      string  `json:"color"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	MaxWidth   float64 `json:"max_width"`
	LineHeight float64 `json:"line_height"`
	MaxLines   int     `json:"max_lines"`
	FitText    bool    `json:"fit_text"`
}

// RenderText renders wrapped and fitted text into the canvas.
func RenderText(dc *gg.Context, spec TextSpec) error {
	dc.Push()
	defer dc.Pop()

	r, g, b, a := icons.ParseHexColor(spec.Color)
	dc.SetRGBA(r, g, b, a)

	fontSize := spec.FontSize
	if fontSize <= 0 {
		fontSize = 32
	}
	lineHeightMult := spec.LineHeight
	if lineHeightMult <= 0 {
		lineHeightMult = 1.25
	}
	maxLines := spec.MaxLines
	if maxLines <= 0 {
		maxLines = 10
	}
	maxWidth := spec.MaxWidth
	if maxWidth <= 0 {
		maxWidth = float64(dc.Width()) - spec.X - 40
	}

	var face font.Face
	var lines []string
	var err error

	// FitText downscaling loop
	minFontSize := 18.0
	for {
		face, err = fonts.GetFontFace(spec.Font, fontSize)
		if err != nil {
			return err
		}
		dc.SetFontFace(face)

		lines = wrapText(dc, spec.Text, maxWidth)

		if !spec.FitText || len(lines) <= maxLines || fontSize <= minFontSize {
			break
		}

		fontSize -= 2.0
	}

	// Truncate lines if still exceeding maxLines
	if len(lines) > maxLines {
		truncated := lines[:maxLines]
		lastLine := truncated[maxLines-1]

		// Add ellipsis
		for len(lastLine) > 0 {
			w, _ := dc.MeasureString(lastLine + "...")
			if w <= maxWidth {
				truncated[maxLines-1] = lastLine + "..."
				break
			}
			lastLine = strings.TrimSpace(lastLine[:len(lastLine)-1])
		}
		lines = truncated
	}

	// Draw lines
	lineSpacing := fontSize * lineHeightMult
	currY := spec.Y + fontSize // align baseline

	for _, line := range lines {
		dc.DrawString(line, spec.X, currY)
		currY += lineSpacing
	}

	return nil
}

// wrapText breaks a string into lines that fit within maxWidth.
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var result []string
	paragraphs := strings.Split(text, "\n")

	for _, p := range paragraphs {
		words := strings.Fields(p)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		currLine := words[0]
		for _, w := range words[1:] {
			candidate := currLine + " " + w
			width, _ := dc.MeasureString(candidate)
			if width <= maxWidth {
				currLine = candidate
			} else {
				result = append(result, currLine)
				currLine = w
			}
		}
		if currLine != "" {
			result = append(result, currLine)
		}
	}

	return result
}
