<p align="center">
  <img src="docs/images/logo.png?raw=true" alt="ogcraft Logo" width="130" style="border-radius: 24px;" />
</p>
<h1 align="center">ogcraft</h1>
<p align="center">
  <strong>Self-hosted dynamic OpenGraph image and social preview card generator in Go and React 19.</strong><br />
  Generates 1200x630 pixel social cards from URL query parameters in under 15ms with two-tier caching and sub-40MB memory footprint.
</p>

<p align="center">
  <img src="docs/images/ogcraft_demo.gif?raw=true" alt="ogcraft Studio and Social Feed Simulator Demo" width="100%" />
</p>

---

## Highlights

- **Fast native rendering**: Pure Go vector graphics pipeline rendering cards in 8 to 15 milliseconds.
- **Multiple output formats**: Native support for **WebP** (sub-25KB), **PNG**, and **JPEG** with automatic content negotiation.
- **Multi-aspect presets**: Standard OpenGraph (1200x630), Square (800x800 for Twitter/Instagram), and Vertical Story (1080x1920).
- **Two-tier cache engine**: Thread-safe in-memory LRU (L1) backed by content-hashed disk cache (L2) with RFC 9111 ETag validation.
- **Embedded typography**: TrueType fonts embedded directly in the binary (Inter Regular, SemiBold, Bold, and JetBrains Mono Medium, Bold).
- **Dynamic FitText**: Automatic multi-line wrapping and downscaling so long titles fit without clipping.
- **Advanced visual elements**: Terminal code windows with macOS controls and syntax coloring, metric cards with progress bars, vector badges, and pattern overlays (dots grid, tech lines).
- **Dynamic palette overrides**: Change colors on the fly via URL query parameters (`?bg=#0f172a`, `?accent=#38bdf8`, `?theme=light`).
- **Visual card designer**: Drag-and-drop React 19 cockpit with a live **Social Feed Simulator** (preview on Twitter/X, LinkedIn, Discord, Slack), JSON import/export, and undo/redo history.
- **Batch SSG CLI**: Pre-render hundreds of cards from JSON or CSV files in seconds for static site generation.
- **Production observability**: Prometheus metrics endpoint (`/metrics`) and authenticated cache purge API (`POST /v1/cache/purge`).
- **SSRF protection**: Remote avatar and logo fetching guards against private IP ranges and internal network requests.
- **Single static binary**: Ships as a standalone executable under 30MB with zero external runtime requirements.

---

## Quickstart

### 1. Download and Run

```bash
# Clone the repository
git clone https://github.com/alexandrmotologa/ogcraft.git
cd ogcraft

# Build single binary
go build -o ogcraft ./cmd/ogcraft

# Start the server (default port: 8080)
./ogcraft serve
```

Open `http://localhost:8080` in your browser to access the visual designer cockpit.

### 2. Request an Image

```bash
# WebP format (sub-25KB)
curl -o card.webp "http://localhost:8080/v1/cards/dark_tech.webp?title=High+Performance+Go+Services&author=Alexandr&tag=Go"

# Standard PNG
curl -o card.png "http://localhost:8080/v1/cards/dark_tech.png?title=High+Performance+Go+Services&author=Alexandr&tag=Go"

# Square aspect ratio for Twitter / Instagram
curl -o square.png "http://localhost:8080/v1/cards/minimal_light.png?aspect=square&title=Square+Post"
```

---

## CLI Usage

`ogcraft` can run as an HTTP daemon, render one-shot images directly, or batch pre-render collections for static site generators:

```bash
# Start HTTP daemon on custom port
ogcraft serve --port 3000 --cache-dir /var/cache/ogcraft

# Render single card to disk without running server
ogcraft render --template dark_tech \
  --format webp \
  --param title="Building Reliable Systems" \
  --param author="Alexandr Motologa" \
  --param tag="Architecture" \
  --output card.webp

# Batch pre-render social cards for static site deployment (SSG)
ogcraft batch --input posts.json \
  --template dark_tech \
  --format webp \
  --out-dir ./dist/og

# Print version and build information
ogcraft version
```

---

## Configuration

Set options via command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
| :--- | :--- | :--- | :--- |
| `--port` | `OGCRAFT_PORT` | `8080` | HTTP listening port |
| `--host` | `OGCRAFT_HOST` | `0.0.0.0` | HTTP listening host address |
| `--cache-dir` | `OGCRAFT_CACHE_DIR` | `./data/cache` | L2 disk cache directory |
| `--l1-bytes` | `OGCRAFT_L1_BYTES` | `25165824` (24MB) | L1 memory cache byte limit |
| `--templates-dir` | `OGCRAFT_TEMPLATES_DIR` | `./data/templates` | User-defined custom templates directory |
| `--secret` | `OGCRAFT_SECRET` | empty | HMAC secret for URL signing (signing disabled when empty) |
| `--rate-limit` | `OGCRAFT_RATE_LIMIT` | `100` | Max requests per second per IP |

---

## Docker Deployment

Run with Docker in seconds:

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name ogcraft \
  ghcr.io/alexandrmotologa/ogcraft:latest
```

---

## Visual Studio & Social Previews

### Interactive Canvas Editor
The embedded React 19 studio provides live element dragging, real-time coordinate inspection, typography adjustment, and URL parameter simulation.

<p align="center">
  <img src="docs/images/ogcraft_dashboard.png?raw=true" alt="ogcraft Interactive Canvas Editor" width="100%" />
</p>

### Social Feed Simulator
Preview how generated OpenGraph preview cards display across Twitter / X, LinkedIn, Discord, and Slack feed cards before deploying to production.

<p align="center">
  <img src="docs/images/ogcraft_social_simulator.png?raw=true" alt="ogcraft Social Feed Simulator" width="100%" />
</p>

### Responsive Mobile Interface
The design studio adapts cleanly to tablets and smartphones with a dedicated bottom navigation bar, touch-friendly parameter controls, and responsive preview scaling.

<p align="center">
  <img src="docs/images/ogcraft_mobile.png?raw=true" alt="ogcraft Responsive Mobile Interface" width="380" />
</p>

---

## Documentation

- [HTTP API Reference](docs/API.md)
- [Template Schema and Presets](docs/TEMPLATES.md)
- [Framework Integration Guides](docs/INTEGRATIONS.md)

---

## Architecture Overview

```
URL Query Parameters ──► HMAC Verification ──► L1 Memory LRU ──► [HIT: 0.1ms]
                                                     │
                                                  [MISS]
                                                     ▼
                                              L2 Disk Cache  ──► [HIT: 1.5ms]
                                                     │
                                                  [MISS]
                                                     ▼
                                            Go 2D Vector Engine
                                            - Text Wrap & FitText
                                            - Code Windows & Syntax
                                            - Metric & Rating Cards
                                            - Pattern Overlays (Dots/Grid)
                                            - SSRF Safe Avatars
                                                     │
                                                     ▼
                                     WebP / PNG / JPEG Fast Encoder ──► [11.2ms]
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
