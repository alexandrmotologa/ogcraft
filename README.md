# ogcraft

Self-hosted dynamic OpenGraph image and social preview card generator in Go and React 19.

Generates 1200x630 pixel social cards from URL query parameters in under 15ms with two-tier caching and sub-40MB memory footprint.

---

## Highlights

- **Fast native rendering**: Pure Go vector graphics pipeline rendering 1200x630 cards in 8 to 15 milliseconds.
- **Low memory consumption**: Under 40MB RAM under sustained traffic. Zero Chromium or headless browser dependencies.
- **Two-tier cache engine**: Thread-safe in-memory LRU (L1) backed by content-hashed disk cache (L2) with RFC 9111 ETag validation.
- **Embedded typography**: TrueType fonts embedded directly in the binary (Inter Regular, SemiBold, Bold, and JetBrains Mono Medium).
- **Dynamic FitText**: Automatic multi-line wrapping and downscaling so long titles fit without clipping.
- **Visual card designer**: Drag-and-drop React 19 cockpit embedded inside the Go executable.
- **SSRF protection**: Remote avatar and logo fetching guards against private IP ranges and internal network requests.
- **HMAC security**: Optional URL signature checking prevents cache poisoning and parameter tampering on public endpoints.
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
curl -o card.png "http://localhost:8080/v1/cards/dark_tech.png?title=High+Performance+Go+Services&author=Alexandr&tag=Go"
```

---

## CLI Usage

`ogcraft` can run as a daemon or render images directly from the command line:

```bash
# Start HTTP daemon on custom port
ogcraft serve --port 3000 --cache-dir /var/cache/ogcraft

# Render single card to disk without running server
ogcraft render --template dark_tech \
  --param title="Building Reliable Systems" \
  --param author="Alexandr Motologa" \
  --param tag="Architecture" \
  --output card.png

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

Or build the local container:

```bash
docker build -t ogcraft .
docker run -p 8080:8080 ogcraft
```

---

## Documentation

- [HTTP API Reference](docs/API.md)
- [Template Schema and Presets](docs/TEMPLATES.md)
- [Framework Integration Guides](docs/INTEGRATIONS.md)

---

## Architecture Overview

```
URL Query Parameters ──► HMAC Verification ──► L1 Memory LRU ──► [HIT: 0.2ms]
                                                     │
                                                  [MISS]
                                                     ▼
                                              L2 Disk Cache  ──► [HIT: 1.8ms]
                                                     │
                                                  [MISS]
                                                     ▼
                                            Go 2D Vector Engine
                                            - Text Wrap & FitText
                                            - TrueType Fonts
                                            - SVG & Shapes
                                            - SSRF Safe Avatars
                                                     │
                                                     ▼
                                             Fast PNG Encoder ──► [11.2ms]
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
