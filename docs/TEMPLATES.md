# ogcraft Template Specification

Templates in `ogcraft` are declared in JSON. Each template defines a 1200x630 pixel canvas, background styling, and a list of visual elements that bind to URL query parameters.

## Canvas Specifications
- Width: 1200 pixels
- Height: 630 pixels (standard OpenGraph 1.904:1 aspect ratio)
- Coordinate system: Origin `(0, 0)` is the top-left corner.

---

## Template Schema

```json
{
  "id": "template_id",
  "name": "Human Readable Name",
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
  "elements": []
}
```

### Background Types
1. **Solid Color**:
   ```json
   {
     "type": "solid",
     "color": "#0a0a0c"
   }
   ```
2. **Linear Gradient**:
   ```json
   {
     "type": "gradient",
     "direction": "to-bottom-right",
     "stops": [
       { "offset": 0.0, "color": "#111827" },
       { "offset": 1.0, "color": "#1f2937" }
     ]
   }
   ```
   Supported directions: `to-right`, `to-bottom`, `to-bottom-right`, `to-bottom-left`.
3. **Radial Accent**:
   ```json
   {
     "type": "radial",
     "center_x": 1000,
     "center_y": 100,
     "radius": 450,
     "center_color": "#3b82f633",
     "outer_color": "#090d16"
   }
   ```

---

## Element Types

### 1. `text`
Renders dynamic or static typography with automatic wrapping, truncation, and dynamic `fit_text` scaling.

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

- `field`: URL query parameter key that supplies the text.
- `default`: Fallback text when parameter is omitted.
- `font`: `Inter-Regular`, `Inter-SemiBold`, `Inter-Bold`, or `JetBrainsMono-Medium`.
- `fit_text`: When true, font size automatically decreases if the text exceeds `max_lines` within `max_width`.

---

### 2. `badge`
Renders a rounded pill or tag with background color, padding, and text.

```json
{
  "type": "badge",
  "field": "tag",
  "default": "ENGINEERING",
  "font": "JetBrainsMono-Medium",
  "font_size": 22,
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

### 3. `avatar`
Renders an author profile image, logo, or icon.

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

Shapes supported: `circle`, `rounded` (uses border_radius 12), `square`.
When remote download fails, `ogcraft` automatically renders a fallback monogram containing the author's initials.

---

### 4. `icon`
Renders embedded vector icons.

```json
{
  "type": "icon",
  "name": "github",
  "x": 80,
  "y": 550,
  "size": 28,
  "color": "#9ca3af"
}
```

Available icons: `github`, `twitter`, `x`, `star`, `clock`, `calendar`, `user`, `tag`, `bolt`, `check`.

---

### 5. `divider`
Renders a horizontal or vertical rule separator.

```json
{
  "type": "divider",
  "x": 80,
  "y": 440,
  "width": 1040,
  "height": 1,
  "color": "#1f2937"
}
```

---

## Built-In Templates

`ogcraft` ships with six built-in presets:

1. **`dark_tech`**: Deep obsidian background with gradient undertone, monospace category badge, primary headline, author avatar, and author name.
2. **`minimal_light`**: Clean off-white background with bold high-contrast typography, category pill, and date.
3. **`glow_dev`**: Cyberpunk dark canvas with radial corner light source and accent cards.
4. **`product_launch`**: Product announcement card with CTA badge and metric badge.
5. **`podcast_event`**: Event and podcast card with date pill, speaker avatar, and episode indicator.
6. **`github_repo`**: Open source repository card with language dot, star count, and fork indicators.
