package fonts

import (
	"embed"
	"fmt"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/*.ttf
var fontFS embed.FS

// Known font identifiers
const (
	InterRegular        = "Inter-Regular"
	InterSemiBold       = "Inter-SemiBold"
	InterBold           = "Inter-Bold"
	JetBrainsMonoMedium = "JetBrainsMono-Medium"
	JetBrainsMonoBold   = "JetBrainsMono-Bold"
)

var (
	parsedFonts = make(map[string]*truetype.Font)
	faceCache   = make(map[string]font.Face)
	cacheMu     sync.RWMutex
	initOnce    sync.Once
)

// Init loads and parses all embedded fonts into memory.
func Init() error {
	var initErr error
	initOnce.Do(func() {
		fontsToLoad := map[string]string{
			InterRegular:        "assets/Inter-Regular.ttf",
			InterSemiBold:       "assets/Inter-SemiBold.ttf",
			InterBold:           "assets/Inter-Bold.ttf",
			JetBrainsMonoMedium: "assets/JetBrainsMono-Medium.ttf",
			JetBrainsMonoBold:   "assets/JetBrainsMono-Bold.ttf",
		}

		for name, path := range fontsToLoad {
			data, err := fontFS.ReadFile(path)
			if err != nil {
				initErr = fmt.Errorf("failed to read embedded font %s: %w", path, err)
				return
			}
			f, err := truetype.Parse(data)
			if err != nil {
				initErr = fmt.Errorf("failed to parse font %s: %w", path, err)
				return
			}
			parsedFonts[name] = f
		}
	})
	return initErr
}

// GetFontFace returns a cached font.Face for the given font name and point size.
func GetFontFace(name string, size float64) (font.Face, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	// Normalization & Fallbacks
	normalized := normalizeFontName(name)

	cacheKey := fmt.Sprintf("%s:%.1f", normalized, size)

	cacheMu.RLock()
	face, exists := faceCache[cacheKey]
	cacheMu.RUnlock()
	if exists {
		return face, nil
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	// Double check after lock
	if face, exists := faceCache[cacheKey]; exists {
		return face, nil
	}

	ttfFont, ok := parsedFonts[normalized]
	if !ok {
		ttfFont = parsedFonts[InterRegular]
	}

	opt := &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}

	face = truetype.NewFace(ttfFont, opt)
	faceCache[cacheKey] = face

	return face, nil
}

// normalizeFontName resolves aliases and case differences to canonical font names.
func normalizeFontName(name string) string {
	switch name {
	case "Inter-Bold", "inter-bold", "bold", "InterBold":
		return InterBold
	case "Inter-SemiBold", "inter-semibold", "semibold", "InterSemiBold":
		return InterSemiBold
	case "Inter-Regular", "inter", "inter-regular", "regular", "default":
		return InterRegular
	case "JetBrainsMono-Medium", "jetbrains-mono", "mono", "monospace", "JetBrainsMono":
		return JetBrainsMonoMedium
	case "JetBrainsMono-Bold", "jetbrains-mono-bold", "mono-bold":
		return JetBrainsMonoBold
	default:
		if _, ok := parsedFonts[name]; ok {
			return name
		}
		return InterRegular
	}
}
