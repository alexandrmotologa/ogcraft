# `ogcraft` Architecture and Implementation Plan

> Self-hosted dynamic OpenGraph image and social preview card generator in Go and React 19. Renders 1200x630 cards in under 15ms with two-tier caching and sub-40MB memory footprint.

---

## 1. System Overview

`ogcraft` generates 1200x630 social preview cards on the fly from URL query parameters. It targets self-hosted blogs, docs, SaaS marketing sites, and e-commerce stores that need dynamic social images without paying for metered SaaS or running heavy headless browsers.

### Performance and Resource Targets
- Rendering latency: Under 15ms per 1200x630 card on standard CPU hardware.
- Memory consumption: Under 40MB RAM under sustained load.
- Output formats: Optimized PNG with sub-50KB average asset size.
- Binary size: Single self-contained static binary under 30MB including embedded fonts, icons, and React 19 designer UI.
- Zero external runtime dependencies: Pure Go rendering, no CGo, no Chromium, no external graphics libraries required on host.

---

## 2. Rendering Pipeline Architecture

```
[ Social Crawler / Browser / Client ]
                  │
                  │ GET /v1/cards/:template.png?title=...&author=...&tag=...
                  ▼
┌────────────────────────────────────────────────────────────────────────┐
│                        ogcraft HTTP Service                            │
│                                                                        │
│  ┌───────────────────────┐         ┌────────────────────────────────┐  │
│  │ HMAC Signature Guard  │ ──────► │ Template Selector and Query    │  │
│  │ (Optional DoS shield) │         │ Parameter Normalizer           │  │
│  └───────────────────────┘         └────────────────┬───────────────┘  │
│                                                     │                  │
│                                                     ▼                  │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ L1 In-Memory LRU Cache (SHA-256 of template + sorted parameters) │  │
│  └────────────────────────┬─────────────────────────────────────────┘  │
│               Cache HIT   │                     │ Cache MISS           │
│       ┌───────────────────┘                     ▼                      │
│       │                    ┌────────────────────────────────────────┐  │
│       │                    │ L2 Deterministic Disk Cache            │  │
│       │                    └────────────────────┬───────────────────┘  │
│       │                             Cache HIT   │   │ Cache MISS       │
│       │                         ┌───────────────┘   ▼                  │
│       │                         │    ┌──────────────────────────────┐  │
│       │                         │    │ Native 2D Vector Engine      │  │
│       │                         │    │ - Dynamic FitText algorithm  │  │
│       │                         │    │ - TrueType font rasterizer   │  │
│       │                         │    │ - SSRF-safe avatar loader    │  │
│       │                         │    │ - Gradient and shape shaders │  │
│       │                         │    │ - Vector icon library        │  │
│       │                         │    └──────────────┬───────────────┘  │
│       │                         │                   │                  │
│       │                         │                   ▼                  │
│       │                         │    ┌──────────────────────────────┐  │
│       │                         │    │ Fast PNG Encoder             │  │
│       │                         │    └──────────────┬───────────────┘  │
│       │                         ▼                   │                  │
│       └────────────────► [ HTTP Response Formatter ]◄┘                  │
│                          - RFC 9111 ETag / 304 Not Modified            │
│                          - Cache-Control: public, max-age=31536000     │
│                                                                        │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ Embedded React 19 Visual Card Designer                           │  │
│  │ - 1200x630 live canvas preview with responsive scaling          │  │
│  │ - Real-time typography, gradient, and layout controls            │  │
│  │ - Dynamic parameter simulator and meta tag generator             │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Technology Stack

| Layer | Component | Implementation |
| :--- | :--- | :--- |
| Core Language | Go 1.24+ | Zero CGo, cross-platform compilation, low memory footprint |
| 2D Vector Engine | Pure Go Graphics | `fogleman/gg` with custom clipping, shape extensions, and gradient rasterizer |
| Typography | TrueType / OpenType Parser | `golang.org/x/image/font` with embedded Inter and JetBrains Mono fonts |
| Caching | Sharded Two-Tier | Thread-safe in-memory LRU (L1) and atomic hashed filesystem cache (L2) |
| Web Server | Standard HTTP / Router | `go-chi/chi/v5` with RFC 9111 ETag middleware, CORS, and rate limiting |
| Visual Designer | React 19 + TypeScript | Vite, Tailwind CSS, Lucide icons, embedded into Go binary via `embed.FS` |
| Security | SSRF and HMAC | Private IP validation for remote avatars and optional HMAC-SHA256 query signatures |

---

## 4. Repository Structure

```
B:\workgit\ogcraft\
├── cmd\
│   └── ogcraft\
│       └── main.go                 # Daemon and CLI entrypoint (serve, render, version)
├── internal\
│   ├── api\
│   │   ├── handler.go              # HTTP routes: /v1/cards, /v1/templates, /v1/render, /health
│   │   ├── middleware.go           # Logging, CORS, rate limiting, and RFC 9111 ETag
│   │   └── signature.go            # HMAC-SHA256 URL token validation
│   ├── cache\
│   │   ├── cache.go                # Unified two-tier cache coordinator
│   │   ├── l1_memory.go            # Thread-safe in-memory LRU with byte limit
│   │   └── l2_disk.go              # Content-hashed atomic file cache
│   ├── config\
│   │   └── config.go               # Server, cache, security, and font configurations
│   ├── fonts\
│   │   ├── embedded.go             # Embedded TrueType font files and font face registry
│   │   └── assets\                 # Inter (Regular, SemiBold, Bold) and JetBrains Mono (Medium)
│   ├── icons\
│   │   └── icons.go                # Embedded vector icons (GitHub, X, Calendar, Star, Clock, etc.)
│   ├── renderer\
│   │   ├── canvas.go               # 1200x630 rendering canvas wrapper
│   │   ├── text.go                 # Dynamic FitText, word wrapping, truncation, and line metrics
│   │   ├── shapes.go               # Badges, rounded boxes, circles, drop shadows, borders
│   │   ├── gradients.go            # Linear, radial, and mesh accent shaders
│   │   ├── images.go               # SSRF-protected remote image fetcher with local caching
│   │   └── render.go               # Master template execution pipeline
│   └── templates\
│       ├── model.go                # JSON template schema specification
│       ├── manager.go              # Template loader, registry, and validation
│       └── builtin\                # Preset templates (dark_tech, minimal_light, glow_dev, product_launch, podcast_event, github_repo)
├── web\
│   ├── src\
│   │   ├── components\
│   │   │   ├── CanvasPreview.tsx   # Pixel-accurate 1200x630 card preview
│   │   │   ├── ElementInspector.tsx# Property controls for selected elements
│   │   │   ├── TemplateSelector.tsx# Built-in template selector
│   │   │   ├── ParamTester.tsx     # Live URL parameter testing form
│   │   │   └── CodeExporter.tsx    # Next.js, HTML meta, Astro, and Markdown snippet generator
│   │   ├── lib\
│   │   │   ├── api.ts              # Client communicating with Go backend
│   │   │   └── types.ts            # Shared template and element types
│   │   ├── App.tsx                 # Root cockpit interface
│   │   ├── main.tsx                # React 19 entrypoint
│   │   └── index.css               # Tailwind CSS theme declarations
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── docs\
│   ├── API.md                      # Complete HTTP endpoint documentation
│   ├── TEMPLATES.md                # Template schema and element reference
│   └── INTEGRATIONS.md             # Integration guides for Next.js, Astro, Hugo, Remix
├── Makefile                        # Build, test, cross-compilation, and benchmark scripts
├── Dockerfile                      # Multi-stage build producing static image
├── README.md                       # Comprehensive project guide and reference
└── MASTER_PLAN.md                  # This document
```

---

## 5. Detailed Implementation Phases

### Phase 1: Go Engine Scaffolding, Fonts, and 2D Vector Engine
- [ ] Initialize Go module `github.com/alexandrmotologa/ogcraft`.
- [ ] Embed high-quality TrueType fonts: Inter (Regular, SemiBold, Bold) and JetBrains Mono (Medium).
- [ ] Build font registry supporting variable size lookup and caching font faces.
- [ ] Implement text layout engine with:
  - Automatic line wrapping based on pixel width.
  - Line truncation with trailing ellipsis when exceeding `max_lines`.
  - Dynamic `FitText` downscaling to guarantee text never overflows designated bounds.
- [ ] Implement shape rasterizer: rounded rectangles, pills, circular avatars, shadows, and borders.
- [ ] Implement linear and radial gradient generators.
- [ ] Implement built-in vector icon renderer for developer and social badges (GitHub, X, Clock, Calendar, Star, Tag, Link).

### Phase 2: Template System and Image Fetcher
- [ ] Define JSON template schema supporting background styles, text blocks, badges, avatars, and icons.
- [ ] Implement built-in template repository:
  1. `dark_tech`: Developer dark mode card with monospace tag and avatar.
  2. `minimal_light`: High-contrast editorial card with clean typography.
  3. `glow_dev`: Cyberpunk radial accent with glassmorphism card container.
  4. `product_launch`: Feature highlight card with CTA pill and metric badge.
  5. `podcast_event`: Event announcement with date badge, host avatar, and episode number.
  6. `github_repo`: Repository card with star count, fork count, and language badge.
- [ ] Implement SSRF-safe remote image downloader:
  - Reject private, loopback, multicast, and link-local IP addresses.
  - 2-second timeout and 2MB payload cap.
  - Image cropping: circle, squircle, and rounded rectangle.
  - Fallback avatar with author initials on download error.

### Phase 3: Two-Tier Cache and HTTP API
- [ ] Build L1 in-memory LRU cache with configurable byte limit (default 24MB).
- [ ] Build L2 disk cache with deterministic SHA-256 keys and atomic rename writes.
- [ ] Build HTTP service with `go-chi/chi/v5`:
  - `GET /v1/cards/{template}.png`: Render social image from URL parameters.
  - `GET /v1/templates`: List all available templates.
  - `GET /v1/templates/{id}`: Fetch template definition.
  - `POST /v1/templates`: Save custom template.
  - `DELETE /v1/templates/{id}`: Remove custom template.
  - `POST /v1/render`: Render dynamic template definition sent in request body.
  - `GET /health`: Diagnostic metrics (memory stats, cache hit ratios, render latency).
- [ ] Implement RFC 9111 HTTP cache headers (`ETag`, `304 Not Modified`, `Cache-Control`).
- [ ] Implement optional HMAC-SHA256 parameter signature guard.

### Phase 4: Embedded React 19 Visual Designer
- [ ] Scaffold React 19 + TypeScript + Vite project inside `web/`.
- [ ] Build live 1200x630 preview canvas with responsive aspect ratio scaling.
- [ ] Build interactive controls:
  - Canvas background: solid color, linear gradient, radial gradient.
  - Element properties: text content, font family, font size, color, coordinates, padding.
  - Add / remove elements (Text, Badge, Avatar, Metric, Icon).
- [ ] Build live URL parameter simulator: adjust query values and see real-time updates.
- [ ] Build code generator: export snippets for Next.js 15, HTML `<meta>`, Astro, Hugo, and Markdown.
- [ ] Package production build and embed into Go binary via `//go:embed`.

### Phase 5: CLI, Testing, Verification, and Documentation
- [ ] Build CLI commands (`ogcraft serve`, `ogcraft render --template ... --output ...`, `ogcraft version`).
- [ ] Write comprehensive unit tests for font metrics, text wrapping, cache eviction, and SSRF prevention.
- [ ] Run latency benchmarks verifying sub-15ms card generation.
- [ ] Write documentation: `README.md`, `docs/API.md`, `docs/TEMPLATES.md`, `docs/INTEGRATIONS.md`.
- [ ] Create multi-stage `Dockerfile` producing static container image under 25MB.
- [ ] Verify clean execution, verify all endpoints, and push commits.
