# ogcraft HTTP API Reference

`ogcraft` exposes an HTTP service for rendering dynamic social preview cards and managing template configurations.

## Endpoints

### 1. Render Card Image

```http
GET /v1/cards/{template}.png
GET /v1/cards/{template}
```

Renders a 1200x630 card using the specified template and query parameters.

#### Path Parameters
- `template` (string, required): Template identifier (e.g. `dark_tech`, `minimal_light`, `glow_dev`, `product_launch`, `podcast_event`, `github_repo`, or a custom saved template id).

#### Query Parameters
Any query parameter passed to this endpoint is matched against element bindings in the template:

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `title` | string | Template default | Primary title text |
| `description` | string | Template default | Secondary subtitle or summary text |
| `author` | string | Template default | Author or speaker name |
| `tag` | string | Template default | Category, topic, or badge text |
| `avatar` | string (URL) | Template default | Remote image URL for author avatar or logo |
| `date` | string | Template default | Publication date or event timestamp |
| `sig` | string | optional | HMAC-SHA256 signature when signature verification is enabled |
| `_refresh` | boolean | false | Bypass cache and re-render card |

#### Response Headers
- `Content-Type`: `image/png`
- `ETag`: Deterministic hash of template id, template version, and query parameters
- `Cache-Control`: `public, max-age=31536000, immutable`
- `X-Cache`: `HIT-L1`, `HIT-L2`, or `MISS`
- `X-Render-Time`: Execution time in milliseconds (e.g. `8.4ms`)

#### HTTP Status Codes
- `200 OK`: Image rendered or returned from cache.
- `304 Not Modified`: Returned when request header `If-None-Match` matches the computed ETag.
- `400 Bad Request`: Unknown template or missing required parameters.
- `401 Unauthorized`: Missing or invalid HMAC signature when signature requirement is active.
- `500 Internal Server Error`: Rendering or font layout failure.

---

### 2. Live Render Payload

```http
POST /v1/render
```

Renders a complete JSON template payload directly into a 1200x630 image. Used by the visual designer for real-time live preview before saving.

#### Request Body
```json
{
  "template": {
    "id": "preview",
    "name": "Live Preview",
    "width": 1200,
    "height": 630,
    "background": {
      "type": "gradient",
      "direction": "to-bottom-right",
      "stops": [
        { "offset": 0.0, "color": "#0d1117" },
        { "offset": 1.0, "color": "#161b22" }
      ]
    },
    "elements": [
      {
        "type": "text",
        "field": "title",
        "default": "Live Preview Card",
        "font": "Inter-Bold",
        "font_size": 56,
        "color": "#f0f6fc",
        "x": 80,
        "y": 200,
        "max_width": 1040,
        "line_height": 1.2,
        "max_lines": 3
      }
    ]
  },
  "params": {
    "title": "Preview with custom text"
  }
}
```

#### Response
- `Content-Type`: `image/png`
- Binary image data.

---

### 3. List Templates

```http
GET /v1/templates
```

Returns all registered templates, including built-in presets and user-created custom templates.

#### Response
```json
[
  {
    "id": "dark_tech",
    "name": "Dark Tech and Developer",
    "builtin": true,
    "elements_count": 5
  },
  {
    "id": "minimal_light",
    "name": "Minimal Editorial Light",
    "builtin": true,
    "elements_count": 4
  }
]
```

---

### 4. Get Template Definition

```http
GET /v1/templates/{id}
```

Retrieves the complete JSON schema for a template.

---

### 5. Create or Update Custom Template

```http
POST /v1/templates
```

Saves a user template to the custom template directory on disk.

---

### 6. Delete Custom Template

```http
DELETE /v1/templates/{id}
```

Deletes a custom template. Built-in templates cannot be deleted and return `403 Forbidden`.

---

### 7. Service Diagnostics

```http
GET /health
```

Returns service health, uptime, memory statistics, and cache hit metrics.

#### Response
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 1420,
  "memory": {
    "alloc_mb": 18.4,
    "sys_mb": 34.2,
    "num_gc": 12
  },
  "cache": {
    "l1_items": 142,
    "l1_bytes": 6842100,
    "l1_hit_rate": 0.84,
    "l2_items": 420
  },
  "renders_total": 560
}
```

---

## HMAC Signature Security

When `OGCRAFT_SECRET` is configured, public clients must sign URL parameters to prevent unauthorized rendering:

```
sig = hex(hmac_sha256(secret, template + "?" + sorted_query_params_without_sig))
```

Example request:
```
GET /v1/cards/dark_tech.png?author=Alexandr&tag=Go&title=Scaling+Systems&sig=a8f2e9104b...
```
