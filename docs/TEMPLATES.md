# ogcraft Template Specification

Templates in `ogcraft` are declared in JSON. Each template defines a 1200x630 pixel canvas, background styling, vector pattern overlays, and a list of visual elements that bind to URL query parameters.

## Canvas Specifications
- Default Width: 1200 pixels
- Default Height: 630 pixels (standard OpenGraph 1.904:1 aspect ratio)
- Multi-Aspect Overrides: Supports `?aspect=square` (800x800), `?aspect=story` (1080x1920), or custom `?w=...&h=...`.

---

## Template Schema

```json
{
  "id": "dark_tech",
  "name": "Dark Tech and Developer",
  "width": 1200,
  "height": 630,
  "background": {
    "type": "gradient",
    "direction": "to-bottom-right",
    "pattern": "dots",
    "pattern_color": "#ffffff0f",
    "pattern_spacing": 28,
    "stops": [
      { "offset": 0.0, "color": "#090d16" },
      { "offset": 1.0, "color": "#121a29" }
    ]
  },
  "elements": []
}
```

### Background Types & Overlays
1. **Solid Color**: `{"type": "solid", "color": "#0a0a0c"}`
2. **Linear Gradient**: `{"type": "gradient", "direction": "to-bottom-right", "stops": [...]}`
3. **Radial Accent**: `{"type": "radial", "center_x": 1050, "center_y": 120, "radius": 600, "center_color": "#8b5cf644", "outer_color": "#05070d"}`
4. **Pattern Overlays**:
   - `dots`: Dot grid overlay across the canvas
   - `grid`: Tech grid lines
   - `lines`: Diagonal 45-degree pinstripes
   - `none`: Smooth background without texture

---

## Element Types

### 1. `text`
Typography with automatic line-wrapping, line truncation, and dynamic `fit_text` downscaling.

```json
{
  "type": "text",
  "field": "title",
  "default": "Default Article Title",
  "font": "Inter-Bold",
  "font_size": 56,
  "color": "#f3f4f6",
  "x": 80,
  "y": 200,
  "max_width": 1040,
  "line_height": 1.25,
  "max_lines": 3,
  "fit_text": true
}
```

---

### 2. `badge`
Pill with background color, padding, rounded corners, and optional vector icon.

```json
{
  "type": "badge",
  "field": "tag",
  "default": "ENGINEERING",
  "icon": "bolt",
  "font": "JetBrainsMono-Medium",
  "font_size": 20,
  "color": "#ffffff",
  "background_color": "#2563eb",
  "x": 80,
  "y": 80,
  "padding_x": 16,
  "padding_y": 8,
  "border_radius": 6
}
```

---

### 3. `code`
Terminal code snippet window with macOS window controls, line numbers, and syntax highlighting.

```json
{
  "type": "code",
  "field": "code",
  "default": "package main\n\nfunc main() {\n    println(\"Zero-CGo OpenGraph Engine\")\n}",
  "language": "go",
  "x": 80,
  "y": 180,
  "width": 680,
  "height": 340,
  "font_size": 20,
  "show_controls": true,
  "line_numbers": true
}
```

---

### 4. `metric`
Stat callout card with big numerical value, subtext label, optional icon, and progress bar.

```json
{
  "type": "metric",
  "field": "stars",
  "value": "14.2k",
  "subtext": "GitHub Stargazers",
  "icon": "star",
  "x": 800,
  "y": 180,
  "width": 320,
  "height": 140,
  "background_color": "#161e2e",
  "color": "#38bdf8",
  "border_color": "#223048",
  "border_width": 1,
  "border_radius": 12,
  "progress": 0.85
}
```

---

### 5. `avatar`
Author profile picture, logo, or icon with SSRF protection, circular clipping, and monogram fallback.

```json
{
  "type": "avatar",
  "field": "avatar",
  "default": "https://avatars.githubusercontent.com/u/583231",
  "x": 80,
  "y": 480,
  "size": 68,
  "shape": "circle",
  "border_color": "#374151",
  "border_width": 2
}
```

---

### 6. `icon`
Vector icon (`github`, `twitter`, `x`, `star`, `clock`, `calendar`, `user`, `tag`, `bolt`, `check`).

```json
{
  "type": "icon",
  "icon": "github",
  "x": 80,
  "y": 550,
  "size": 32,
  "color": "#9ca3af"
}
```

---

### 7. `divider` & `container`
Horizontal/vertical dividers and background frames for visual separation.
