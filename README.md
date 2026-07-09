# Blogger XML Exporter

A web application that exports Blogger posts to XML. Select a blog post (with full-text search), fill a configurable form, and generate a custom XML document. Works with or without a blog post selected.

## Features

- **Blogger Integration**: Fetch posts from Google Blogger API with searchable dropdown
- **Flexible Form**: Define form layout in `config.yaml` (groups, fields, layout, presets)
- **Separate Concerns**: Form layout and XML output structure are configured independently
- **Pre-filling**: Auto-populate form fields from post data, templates, or static values
- **Presets**: Quick-fill templates that populate multiple form fields at once
- **Configurable Theme**: Set brand colors via `config.yaml`
- **No Build Step**: Frontend uses Alpine.js, Tom Select, Flatpickr (all vendored locally)
- **Minimal Docker Image**: ~10 MB using `distroless` base
- **Optional Assets**: Mount custom favicon/logo via external volume

## Quick Start

### Prerequisites

- Go 1.18+
- Blogger API enabled (see [Google Cloud Setup](#google-cloud-setup))
- `make` (for build tasks)

### Installation

```bash
git clone https://github.com/yourusername/blogger-xml-exporter
cd blogger-xml-exporter

cp config.example.yaml config.yaml
# Edit config.yaml: set blogger.blogId

export BLOGGER_API_KEY="your-api-key"
make build
./blogger-xml-exporter
```

Server runs on `http://localhost:8080`.

### Docker

```bash
docker build -t blogger-xml-exporter .

docker run \
  -e BLOGGER_API_KEY="your-key" \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -p 8080:8080 \
  blogger-xml-exporter
```

With custom assets:
```bash
docker run \
  -e BLOGGER_API_KEY="your-key" \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/assets:/app/assets \
  -p 8080:8080 \
  blogger-xml-exporter
```

## Configuration

### Structure

`config.yaml` separates two independent concerns:

- **`form.items`**: Form layout tree (groups, fields, presets, rows, widths)
- **`xml.fields`**: XML output structure (element paths and value sources)

This separation avoids duplicating structure. A form field pulls from a post, but the XML output can reshape it completely. Static XML values (e.g., company metadata) don't need hidden form fields.

### Blogger Settings

```yaml
blogger:
  blogId: "1234567890123456789"
  maxResults: 20
```

### Server & Branding

```yaml
server:
  port: 8080

site:
  title: "XML Exporter"
  heading: "XML Exporter"

theme:
  primaryColor: "#2563eb"        # Brand color (hex)
  darkColor: "#1e40af"           # Dark variant (optional)
  lightColor: "#3b82f6"          # Light variant (optional)
```

Omit theme colors to use defaults (blue shades). Colors are applied as CSS variables at runtime.

### Assets

```yaml
assets:
  dir: "/app/assets"             # Optional external directory
  favicon: "favicon.ico"         # Filename within dir
  logo: "logo.png"               # Filename within dir
```

### Form Definition

```yaml
form:
  items:
    - type: group
      title: "Article"
      items:
        - type: text
          name: title
          label: "Title"
          source: title            # Pull from post.title
          row: 1
          width: 8
        - type: textarea
          name: description
          label: "Description"
          template: '{{ .title }} summary'  # Template-based
          row: 1
          width: 4
```

Field types: `text`, `textarea`, `date`, `list`, `select`, `array`, `group`.

Each field has exactly one of:
- `source: <dot.path>` - pull directly from post JSON
- `template: <go-template>` - computed value
- Neither - user fills manually

### XML Output

```yaml
xml:
  root: "item"
  filename: '{{ .title }}'  # Download filename (without .xml)
  namespaces:
    - name: ex
      value: "http://example.org/schema"
  fields:
    - xmlPath: "header/title"
      formField: title
    - xmlPath: "body/summary"
      template: '{{ source "content" | stripHTML }}'
    - xmlPath: "meta/generator"
      template: "Blogger XML Exporter"
```

Each XML field references a form field, a template, or is static.

### Presets (Quick-Fill Templates)

Groups can define reusable preset templates:

```yaml
- type: group
  title: "Author"
  presets:
    - label: "John Doe"
      values:
        first_name: "John"
        last_name: "Doe"
        email: "john@example.com"
```

See `config.example.yaml` for more detailed examples.

## Development

### Build Tasks

```bash
make setup       # Install Go tools + Tailwind CLI
make dev         # Run with live-reload + CSS watch
make build       # Build binary + CSS
make build-css   # Compile Tailwind CSS only
make test        # Run tests
make lint        # Run linter
```

### Environment Variables

```bash
export BLOGGER_API_KEY="your-api-key"
export CONFIG_PATH="config.yaml"
export STATIC_DIR="web/static"
```

### Project Structure

```
.
├── main.go                      # Entry point
├── config.yaml                  # Instance config (not in repo)
├── config.example.yaml          # Template with docs
├── internal/
│   ├── config/      config.go   # Config loading
│   ├── blogger/     client.go   # Blogger API client
│   ├── httpapi/     handlers.go # HTTP handlers
│   └── xmlgen/      render.go   # XML generation
├── web/
│   ├── static/      index.html, js/, css/
│   └── tailwind.src.css
└── Dockerfile, Makefile, go.mod
```

### CSS Compilation

Tailwind CSS v4 Standalone (no Node.js):

```bash
make build-css
```

Outputs to `web/static/css/style.css` (versioned in repo).

## API Endpoints

| Method | Path                    | Returns                           |
|--------|-------------------------|-----------------------------------|
| GET    | `/healthz`              | Health check                      |
| GET    | `/api/form-schema`      | Form layout + site config + theme |
| GET    | `/api/defaults`         | Pre-fill without post selection   |
| GET    | `/api/posts`            | List recent posts                 |
| GET    | `/api/posts?q=<search>` | Full-text search                  |
| GET    | `/api/posts/<id>`       | Single post + pre-filled form     |
| POST   | `/api/generate`         | Generate & download XML           |

## Google Cloud Setup

1. Create a Google Cloud project
2. Enable Blogger API
3. Create an API key (OAuth 2.0)
4. Export as `BLOGGER_API_KEY`

See [Google Cloud Docs](https://developers.google.com/blogger/docs/3.0/getting_started) for details.

## License

MIT
