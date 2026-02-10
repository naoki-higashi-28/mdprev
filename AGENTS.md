# AGENTS.md

## Project Overview

mdprev — A local Markdown preview tool. Go backend + React frontend.

## Repository Structure

```
cmd/mdprev/           # Go entrypoint
internal/
  domain/             # Domain layer (entities, value objects)
  usecase/            # Use case layer
  interface/
    handler/          # HTTP handlers
    middleware/       # Middleware (security, etc.)
  infrastructure/     # Infrastructure layer (filesystem operations)
web/                  # React frontend (Vite + pnpm)
web/dist/             # Build artifacts (.gitignore target)
docs/                 # Specifications
```

## Architecture

### Backend (Go)

- **Clean Architecture + DDD**
- Dependency direction: handler → usecase → domain ← infrastructure
- Domain layer has no external dependencies
- Infrastructure layer implements domain interfaces

```
handler (interface layer)
  ↓ calls
usecase (use case layer)
  ↓ depends on
domain (domain layer: entities, repository interfaces)
  ↑ implements
infrastructure (infrastructure layer: filesystem)
```

### Frontend (React)

- **Vertical Slice Architecture**
- Organized by feature (shared parts placed in `shared/`)

```
web/src/
  features/
    tree/             # Tree feature (components, hooks, API calls)
    preview/          # Preview feature
    pathbar/          # Path bar feature
  shared/
    components/       # Shared UI components
    hooks/            # Shared hooks
    lib/              # Utilities
    types/            # Shared type definitions
  App.tsx
  main.tsx
```

## Development Environment

### Version Management: mise

```toml
# .mise.toml
[tools]
go = "1.24"
node = "22"
pnpm = "10"
```

- Run `mise install` to set up Go, Node.js, and pnpm

### Backend

- **Language**: Go 1.24
- **Router**: Standard library (`net/http`)

### Frontend

- **Language**: TypeScript (strict mode)
- **Framework**: React
- **Build**: Vite
- **Package Manager**: pnpm
- **CSS**: Tailwind CSS

## Linter / Formatter

### Frontend: Biome

```jsonc
// web/biome.json
{
  "formatter": {
    "indentStyle": "space",
    "indentWidth": 2
  },
  "linter": {
    "enabled": true
  }
}
```

- **Format**: `pnpm biome format --write .`
- **Lint**: `pnpm biome lint .`
- **Both**: `pnpm biome check --write .`

### Backend: gofmt / go vet

- **Format**: `gofmt -w .`
- **Static Analysis**: `go vet ./...`

## Common Commands

```bash
# --- Setup ---
mise install                      # Install Go, Node.js, pnpm
cd web && pnpm install             # Install frontend dependencies

# --- Frontend ---
cd web && pnpm dev                 # Start Vite dev server (localhost:3000, API proxied to Go)
cd web && pnpm build               # Build to web/dist/
cd web && pnpm biome check --write . # lint + format

# --- Backend ---
go run ./cmd/mdprev                # Start server (random port, requires web/dist built)
go test ./...                      # Run tests
gofmt -w .                         # Format
go vet ./...                       # Static analysis

# --- Build ---
cd web && pnpm build && cd .. && go build -o mdprev ./cmd/mdprev
```

## Coding Conventions

### Go

- Follow standard Go style (gofmt compliant)
- Return errors to callers (`fmt.Errorf("xxx: %w", err)`)
- No external package dependencies in the domain layer

### TypeScript / React

- Follow Biome rules
- Use function components + hooks
- Prefer `interface` for types (use `type` only when needed, e.g., unions)
- Keep feature-scoped logic within the feature directory

## Git Workflow

### Commit Messages

- Follow Semantic Commit Messages (written in English)
- Title only, no description body required

```
feat: add user authentication API
fix: correct file path validation
refactor: simplify error handling in handler layer
docs: add Git workflow rules to CLAUDE.md
chore: update Biome version
```
