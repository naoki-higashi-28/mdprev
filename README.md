# mdprev

A local Markdown preview tool. Browse and preview Markdown files in your browser — no internet required.

<video src="https://github.com/user-attachments/assets/2e9fa2a7-f5e9-444c-a2c6-d95fc24bb8e2" controls width="100%"></video>

## Features

- **Directory tree** — Navigate your folder structure with an expandable sidebar
- **Markdown preview** — Full GitHub Flavored Markdown support (tables, task lists, strikethrough, etc.)
- **Math rendering** — Inline `$...$` and block `$$...$$` equations via KaTeX
- **Table of contents** — Auto-generated from headings for quick navigation
- **File search** — Find Markdown files by name
- **Live reload** — Preview updates automatically when files change
- **Frontmatter** — YAML frontmatter parsing and display
- **Offline** — Everything runs locally with no external dependencies

## Install

### Go

```bash
go install github.com/naoki-higashi-28/mdprev/cmd/mdprev@latest
```

### Download binary

Pre-built binaries for Linux, macOS, and Windows are available on the [Releases](https://github.com/naoki-higashi-28/mdprev/releases) page.

### Build from source

Requirements: Go 1.24+, Node.js 22+, pnpm 10+ (or install all with `mise install`)

```bash
cd web && pnpm install && pnpm build && cd ..
cp -r web/dist cmd/mdprev/dist
go build -o mdprev ./cmd/mdprev
```

## Usage

```bash
# Preview current directory
mdprev

# Specify host and port
mdprev --host 0.0.0.0 --port 8080 ~/docs
```

Open the printed URL in your browser.

## License

MIT
