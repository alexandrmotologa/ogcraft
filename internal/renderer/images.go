package renderer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/icons"
	"github.com/fogleman/gg"
)

// AvatarSpec specifies image/avatar rendering.
type AvatarSpec struct {
	URL          string  `json:"url"`
	FallbackText string  `json:"fallback_text"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Size         float64 `json:"size"`
	Shape        string  `json:"shape"` // "circle", "rounded", "square"
	BorderColor  string  `json:"border_color,omitempty"`
	BorderWidth  float64 `json:"border_width,omitempty"`
}

type cachedAvatar struct {
	img       image.Image
	fetchedAt time.Time
	failed    bool
}

var (
	avatarMemoryCache = make(map[string]cachedAvatar)
	avatarCacheMu     sync.RWMutex
	httpClient        = &http.Client{
		Timeout: 1500 * time.Millisecond,
		Transport: &http.Transport{
			DialContext: safeDialContext,
		},
	}
)

// RenderAvatar renders an image or fallback monogram onto the canvas.
func RenderAvatar(dc *gg.Context, spec AvatarSpec) error {
	dc.Push()
	defer dc.Pop()

	size := spec.Size
	if size <= 0 {
		size = 64
	}

	var img image.Image
	if spec.URL != "" && spec.URL != "none" {
		img = fetchImageSafe(spec.URL)
	}

	if img != nil {
		drawClippedImage(dc, img, spec)
	} else {
		drawFallbackMonogram(dc, spec)
	}

	if spec.BorderColor != "" && spec.BorderWidth > 0 {
		dc.Push()
		bR, bG, bB, bA := icons.ParseHexColor(spec.BorderColor)
		dc.SetRGBA(bR, bG, bB, bA)
		dc.SetLineWidth(spec.BorderWidth)

		switch strings.ToLower(spec.Shape) {
		case "circle":
			radius := size / 2
			dc.DrawCircle(spec.X+radius, spec.Y+radius, radius)
		case "rounded":
			dc.DrawRoundedRectangle(spec.X, spec.Y, size, size, 12)
		default:
			dc.DrawRectangle(spec.X, spec.Y, size, size)
		}
		dc.Stroke()
		dc.Pop()
	}

	return nil
}

func drawClippedImage(dc *gg.Context, img image.Image, spec AvatarSpec) {
	dc.Push()
	size := spec.Size

	switch strings.ToLower(spec.Shape) {
	case "circle":
		radius := size / 2
		dc.DrawCircle(spec.X+radius, spec.Y+radius, radius)
	case "rounded":
		dc.DrawRoundedRectangle(spec.X, spec.Y, size, size, 12)
	default:
		dc.DrawRectangle(spec.X, spec.Y, size, size)
	}
	dc.Clip()

	bounds := img.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	scale := size / imgW
	if size/imgH > scale {
		scale = size / imgH
	}

	targetW := imgW * scale
	targetH := imgH * scale
	offsetX := spec.X + (size-targetW)/2
	offsetY := spec.Y + (size-targetH)/2

	dc.DrawImageAnchored(img, int(offsetX), int(offsetY), 0, 0)
	dc.Pop()
}

func drawFallbackMonogram(dc *gg.Context, spec AvatarSpec) {
	dc.Push()
	size := spec.Size
	radius := size / 2

	bgR, bgG, bgB := 0.18, 0.22, 0.30
	dc.SetRGB(bgR, bgG, bgB)

	switch strings.ToLower(spec.Shape) {
	case "circle":
		dc.DrawCircle(spec.X+radius, spec.Y+radius, radius)
	case "rounded":
		dc.DrawRoundedRectangle(spec.X, spec.Y, size, size, 12)
	default:
		dc.DrawRectangle(spec.X, spec.Y, size, size)
	}
	dc.Fill()

	initials := extractInitials(spec.FallbackText)
	if initials == "" {
		initials = "?"
	}

	fontSize := size * 0.42
	face, err := fonts.GetFontFace(fonts.InterBold, fontSize)
	if err == nil {
		dc.SetFontFace(face)
		dc.SetRGB(0.9, 0.94, 1.0)
		tW, tH := dc.MeasureString(initials)
		dc.DrawString(initials, spec.X+(size-tW)/2, spec.Y+(size+tH*0.75)/2)
	}

	dc.Pop()
}

func extractInitials(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		runes := []rune(parts[0])
		if len(runes) > 2 {
			return strings.ToUpper(string(runes[:2]))
		}
		return strings.ToUpper(string(runes))
	}
	r1 := []rune(parts[0])
	r2 := []rune(parts[len(parts)-1])
	if len(r1) > 0 && len(r2) > 0 {
		return strings.ToUpper(string(r1[0]) + string(r2[0]))
	}
	return ""
}

func fetchImageSafe(rawURL string) image.Image {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil
	}

	cacheKey := hashString(rawURL)

	avatarCacheMu.RLock()
	cached, ok := avatarMemoryCache[cacheKey]
	avatarCacheMu.RUnlock()
	if ok {
		// Cache valid for 30 minutes
		if time.Since(cached.fetchedAt) < 30*time.Minute {
			return cached.img
		}
	}

	resp, err := httpClient.Get(rawURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		// Cache negative lookup for 5 minutes to avoid spamming network
		avatarCacheMu.Lock()
		avatarMemoryCache[cacheKey] = cachedAvatar{img: nil, fetchedAt: time.Now(), failed: true}
		avatarCacheMu.Unlock()
		return nil
	}
	defer resp.Body.Close()

	limitReader := io.LimitReader(resp.Body, 2*1024*1024)
	decoded, _, err := image.Decode(limitReader)
	if err != nil {
		avatarCacheMu.Lock()
		avatarMemoryCache[cacheKey] = cachedAvatar{img: nil, fetchedAt: time.Now(), failed: true}
		avatarCacheMu.Unlock()
		return nil
	}

	avatarCacheMu.Lock()
	if len(avatarMemoryCache) > 300 {
		avatarMemoryCache = make(map[string]cachedAvatar)
	}
	avatarMemoryCache[cacheKey] = cachedAvatar{img: decoded, fetchedAt: time.Now(), failed: false}
	avatarCacheMu.Unlock()

	return decoded
}

func safeDialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return nil, fmt.Errorf("SSRF blocked: host %s resolves to private/restricted IP %s", host, ip)
		}
	}

	dialer := &net.Dialer{Timeout: 1500 * time.Millisecond}
	return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}

	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
