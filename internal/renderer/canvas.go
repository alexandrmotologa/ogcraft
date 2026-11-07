package renderer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/fogleman/gg"
	webp "github.com/mayahiro/go-webp"
)

// Supported output formats
const (
	FormatPNG  = "png"
	FormatWebP = "webp"
	FormatJPEG = "jpeg"
	FormatJPG  = "jpg"
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

// EncodeWebP encodes the canvas into a high-efficiency WebP byte buffer.
func (c *Canvas) EncodeWebP(quality int) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	var buf bytes.Buffer
	err := webp.Encode(&buf, c.dc.Image(), &webp.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, fmt.Errorf("webp encoding failed: %w", err)
	}
	return buf.Bytes(), nil
}

// EncodeJPEG encodes the canvas into a standard JPEG byte buffer.
func (c *Canvas) EncodeJPEG(quality int) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, c.dc.Image(), &jpeg.Options{
		Quality: quality,
	})
	if err != nil {
		return nil, fmt.Errorf("jpeg encoding failed: %w", err)
	}
	return buf.Bytes(), nil
}

// EncodeFormat encodes the canvas into the requested format (png, webp, jpg, jpeg).
func (c *Canvas) EncodeFormat(format string) ([]byte, string, error) {
	switch strings.ToLower(strings.TrimPrefix(format, ".")) {
	case FormatWebP:
		data, err := c.EncodeWebP(85)
		return data, "image/webp", err
	case FormatJPEG, FormatJPG:
		data, err := c.EncodeJPEG(85)
		return data, "image/jpeg", err
	case FormatPNG:
		fallthrough
	default:
		data, err := c.EncodePNG()
		return data, "image/png", err
	}
}
