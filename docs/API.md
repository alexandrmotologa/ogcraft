# ogcraft HTTP API Reference

`ogcraft` exposes an HTTP service for rendering dynamic social preview cards, serving Prometheus metrics, and managing template configurations.

## Endpoints

### 1. Render Card Image

```http
GET /v1/cards/{template}.png
GET /v1/cards/{template}.webp
GET /v1/cards/{template}.jpg
GET /v1/cards/{template}
```

Renders a social card using the specified template, format, and query parameters.

#### Path Parameters
- `template` (string, required): Template identifier with optional extension (e.g. `dark_tech`, `dark_tech.webp`, `minimal_light.png`, `product_launch.jpg`).

#### Query Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `format` | string | `png` | Output format: `png`, `webp`, `jpeg`, `jpg` |
| `aspect` | string | `og` | Aspect preset: `og` (1200x630), `square` (800x800), `story` (1080x1920) |
| `w` | integer | `1200` | Custom canvas width (200 to 3840 pixels) |
| `h` | integer | `630` | Custom canvas height (200 to 3840 pixels) |
| `theme` | string | - | Theme override: `dark` or `light` |
| `bg` | hex string | - | Custom background color override (e.g. `0d1117`) |
| `accent` | hex string | - | Accent color override for badges, icons, and highlights |
| `pattern` | string | - | Background pattern overlay: `dots`, `grid`, `lines`, `none` |
| `title` | string | Template default | Primary title text |
| `description` | string | Template default | Secondary subtitle or summary text |
| `author` | string | Template default | Author or speaker name |
| `tag` | string | Template default | Category, topic, or badge text |
| `avatar` | string (URL) | Template default | Remote image URL for author avatar or logo |
| `date` | string | Template default | Publication date or event timestamp |
| `code` | string | Template default | Code snippet text for terminal code blocks |
| `metric` | string | Template default | Numerical or metric stat value |
| `sig` | string | optional | HMAC-SHA256 signature when signature verification is active |
| `_refresh` | boolean | false | Bypass cache and force re-render |

#### Response Headers
- `Content-Type`: `image/png`, `image/webp`, or `image/jpeg`
- `ETag`: Deterministic RFC 9111 hash `W/"..."`
- `Cache-Control`: `public, max-age=31536000, immutable`
- `X-Cache`: `HIT-L1`, `HIT-L2`, or `MISS`
- `X-Render-Time`: Execution latency in milliseconds

---

### 2. Live Render Payload

```http
POST /v1/render
```

Renders a complete JSON template definition into an image buffer on the fly. Used by the visual designer for zero-latency live preview.

---

### 3. List Templates

```http
GET /v1/templates
```

Returns all registered templates, including built-in presets and user-created custom templates.

---

### 4. Template CRUD

- `GET /v1/templates/{id}`: Fetch template schema definition
- `POST /v1/templates`: Save custom template to disk
- `DELETE /v1/templates/{id}`: Delete custom template

---

### 5. Cache Purge

```http
POST /v1/cache/purge
```

Flushes all cached images from both L1 RAM and L2 disk store.
When `OGCRAFT_SECRET` is set, requires `Authorization: Bearer <secret>` or `?secret=<secret>`.

---

### 6. Prometheus Metrics

```http
GET /metrics
```

Returns standard Prometheus metrics in text format version 0.0.4:
- `ogcraft_uptime_seconds`: Process uptime
- `ogcraft_renders_total`: Total number of cards rendered
- `ogcraft_cache_hits_total{tier="l1|l2"}`: Cache hits by tier
- `ogcraft_cache_misses_total{tier="l1|l2"}`: Cache misses by tier
- `ogcraft_cache_bytes{tier="l1"}`: Current RAM cache byte usage
- `ogcraft_cache_items{tier="l1"}`: Number of items cached in RAM
- `ogcraft_memory_alloc_bytes`: Memory allocated and currently in use
- `ogcraft_goroutines`: Active goroutine count

---

### 7. Diagnostics

```http
GET /health
```

Returns JSON health status, memory allocation, and cache ratios.
