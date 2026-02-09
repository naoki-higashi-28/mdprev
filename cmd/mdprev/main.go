package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

func main() {
	addr := flag.String("addr", ":0", "listen address (e.g. :8080)")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(sub)))

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting server on http://localhost:%d", ln.Addr().(*net.TCPAddr).Port)
	if err := http.Serve(ln, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
