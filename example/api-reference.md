# API Reference

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/tree` | Get directory tree |
| `GET` | `/api/preview` | Get Markdown content |
| `GET` | `/api/raw` | Get raw file content |
| `WS` | `/ws` | WebSocket for live reload |

## Code Examples

### Go

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, mdprev!")
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
```

### TypeScript

```typescript
interface FileNode {
  name: string;
  path: string;
  isDir: boolean;
  children?: FileNode[];
}

async function fetchTree(root: string): Promise<FileNode> {
  const response = await fetch(`/api/tree?root=${encodeURIComponent(root)}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch tree: ${response.statusText}`);
  }
  return response.json();
}
```

### Shell

```bash
# Install mdprev
go install github.com/naoki-higashi-28/mdprev/cmd/mdprev@latest

# Preview current directory
mdprev

# Specify port
mdprev --port 8080 ~/docs
```

### JSON

```json
{
  "name": "example",
  "path": "/example",
  "isDir": true,
  "children": [
    {
      "name": "README.md",
      "path": "/example/README.md",
      "isDir": false
    }
  ]
}
```

## Error Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | Request succeeded |
| 400 | Bad Request | Missing or invalid parameters |
| 403 | Forbidden | Path traversal attempt blocked |
| 404 | Not Found | File does not exist |
| 500 | Internal Server Error | Unexpected server error |
