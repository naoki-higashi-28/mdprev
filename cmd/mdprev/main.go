package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/naoki-higashi-28/mdprev/internal/dependency"
)

//go:embed all:dist
var distFS embed.FS

var defaultPort = "0"

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	default:
		return
	}

	args = append(args, url)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func main() {
	host := flag.String("host", "127.0.0.1", "bind host")
	port := flag.String("port", defaultPort, "listen port (0 for random)")
	open := flag.Bool("open", true, "open browser automatically")
	flag.Parse()

	// Determine root directory
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Fatalf("Failed to resolve root path: %v", err)
	}

	// Verify root directory exists
	info, err := os.Stat(absRoot)
	if err != nil || !info.IsDir() {
		log.Fatalf("Root path is not a valid directory: %s", absRoot)
	}

	// Setup dependencies
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}
	mux, watcher, err := dependency.NewServeMux(absRoot, sub)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer watcher.Close()

	// Start server
	addr := net.JoinHostPort(*host, *port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	url := fmt.Sprintf("http://%s", ln.Addr().String())
	fmt.Printf("mdprev serving %s on %s\n", absRoot, url)

	if *open {
		openBrowser(url)
	}

	if err := http.Serve(ln, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
