package renderer

import (
	"bytes"
	"image"
	"image/png"

	"github.com/fogleman/gg"
)

// Canvas encapsulates a 2D drawing context for OpenGraph cards.
type Canvas struct {
	Width  int
	Height int
	dc     *gg.Context
}

// NewCanvas creates a new Canvas with the specified dimensions.
func NewCanvas(width, height int) *Canvas {
	if width <= 0 {
		width = 1200
	}
	if height <= 0 {
		height = 630
	}
	dc := gg.NewContext(width, height)
	return &Canvas{
		Width:  width,
		Height: height,
		dc:     dc,
	}
}

// Context returns the underlying *gg.Context.
func (c *Canvas) Context() *gg.Context {
	return c.dc
}

// Image returns the rendered image.Image.
func (c *Canvas) Image() image.Image {
	return c.dc.Image()
}

// EncodePNG encodes the canvas image into a PNG byte buffer using fast compression.
func (c *Canvas) EncodePNG() ([]byte, error) {
	var buf bytes.Buffer
	enc := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	if err := enc.Encode(&buf, c.dc.Image()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
