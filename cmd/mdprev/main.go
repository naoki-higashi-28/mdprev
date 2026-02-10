package main

import (
	"context"
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
	"time"

	"github.com/naoki-higashi-28/mdprev/internal/dependency"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure"
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
	autoClose := flag.Bool("auto-close", true, "shutdown when all connections close")
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

	var onConnect, onDisconnect func()
	var tracker *infrastructure.ConnectionTracker
	if *autoClose {
		tracker = infrastructure.NewConnectionTracker(5 * time.Second)
		onConnect = tracker.Add
		onDisconnect = tracker.Remove
	}

	mux, watcher, err := dependency.NewServeMux(absRoot, sub, onConnect, onDisconnect)
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

	srv := &http.Server{Handler: mux}

	if tracker != nil {
		go func() {
			<-tracker.Done()
			fmt.Println("\nNo active connections. Shutting down...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}()
	}

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
