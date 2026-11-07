package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexandrmotologa/ogcraft/internal/api"
	"github.com/alexandrmotologa/ogcraft/internal/cache"
	"github.com/alexandrmotologa/ogcraft/internal/config"
	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/renderer"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
	"github.com/alexandrmotologa/ogcraft/web"
)

var (
	// Version is injected during build via -ldflags
	Version = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		runServe(os.Args[1:])
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "serve":
		runServe(args)
	case "render":
		runRender(args)
	case "version", "--version", "-v":
		fmt.Printf("ogcraft v%s (Go vector engine, React 19 visual designer)\n", Version)
	case "help", "--help", "-h":
		printUsage()
	default:
		// If first arg starts with flag (e.g. --port 8080), treat as serve
		if strings.HasPrefix(command, "-") {
			runServe(os.Args[1:])
			return
		}
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`ogcraft — Self-hosted dynamic OpenGraph image and social preview card generator

Usage:
  ogcraft [command] [options]

Commands:
  serve     Start the HTTP daemon and embedded visual designer (default)
  render    Render a social card directly from CLI to file
  version   Print version information
  help      Show this help message

Options for 'serve':
  --port, -p           HTTP listening port (default: 8080)
  --host, -h           HTTP listening host (default: 0.0.0.0)
  --cache-dir          L2 disk cache directory (default: ./data/cache)
  --templates-dir      Custom templates directory (default: ./data/templates)
  --l1-bytes           L1 memory cache limit in bytes (default: 25165824 [24MB])
  --secret             HMAC-SHA256 signature secret for URL protection
  --rate-limit         Max requests per second per IP (default: 100)

Options for 'render':
  --template, -t       Template ID to render (default: dark_tech)
  --param, -p          Dynamic parameter formatted as key=value (repeatable)
  --output, -o         Output file path (default: card.png)
  --templates-dir      Directory containing custom templates`)
}

type paramSlice []string

func (p *paramSlice) String() string {
	return strings.Join(*p, ", ")
}

func (p *paramSlice) Set(val string) error {
	*p = append(*p, val)
	return nil
}

func runRender(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	tmplID := fs.String("template", "dark_tech", "Template ID")
	fs.StringVar(tmplID, "t", "dark_tech", "Template ID (shorthand)")

	output := fs.String("output", "card.png", "Output file path")
	fs.StringVar(output, "o", "card.png", "Output file path (shorthand)")

	templatesDir := fs.String("templates-dir", "./data/templates", "Templates directory")

	var params paramSlice
	fs.Var(&params, "param", "Parameter key=value")
	fs.Var(&params, "p", "Parameter key=value (shorthand)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if err := fonts.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize fonts: %v\n", err)
		os.Exit(1)
	}

	tmplMgr, err := templates.NewManager(*templatesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load templates: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := tmplMgr.Get(*tmplID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	paramMap := make(map[string]string)
	for _, p := range params {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			paramMap[parts[0]] = parts[1]
		}
	}

	start := time.Now()
	imgBytes, err := renderer.RenderCard(tmpl, paramMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rendering failed: %v\n", err)
		os.Exit(1)
	}
	duration := time.Since(start)

	if err := os.WriteFile(*output, imgBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rendered %s -> %s (%d KB, %6.2fms)\n", *tmplID, *output, len(imgBytes)/1024, float64(duration.Microseconds())/1000.0)
}

func runServe(args []string) {
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()

	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.IntVar(&cfg.Port, "port", cfg.Port, "HTTP listening port")
	fs.IntVar(&cfg.Port, "p", cfg.Port, "HTTP listening port (shorthand)")
	fs.StringVar(&cfg.Host, "host", cfg.Host, "HTTP listening host")
	fs.StringVar(&cfg.Host, "h", cfg.Host, "HTTP listening host (shorthand)")
	fs.StringVar(&cfg.CacheDir, "cache-dir", cfg.CacheDir, "Disk cache directory")
	fs.StringVar(&cfg.TemplatesDir, "templates-dir", cfg.TemplatesDir, "Custom templates directory")
	fs.Int64Var(&cfg.L1MaxBytes, "l1-bytes", cfg.L1MaxBytes, "L1 RAM cache byte limit")
	fs.StringVar(&cfg.Secret, "secret", cfg.Secret, "HMAC signing secret")
	fs.IntVar(&cfg.RateLimit, "rate-limit", cfg.RateLimit, "Max requests per second per IP")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if err := fonts.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize fonts: %v\n", err)
		os.Exit(1)
	}

	tmplMgr, err := templates.NewManager(cfg.TemplatesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize template manager: %v\n", err)
		os.Exit(1)
	}

	cacheMgr, err := cache.NewTwoTierCache(cfg.L1MaxBytes, cfg.CacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize two-tier cache: %v\n", err)
		os.Exit(1)
	}

	webFS := web.GetFS()
	server := api.NewServer(cfg, tmplMgr, cacheMgr, webFS)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      server.Router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Print startup cockpit banner
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("  │  ogcraft v%-48s│\n", Version)
	fmt.Println("  │  Self-Hosted Dynamic OpenGraph Image & Preview Generator    │")
	fmt.Println("  ├─────────────────────────────────────────────────────────────┤")
	fmt.Printf("  │  Visual Studio Cockpit : http://localhost:%-5d              │\n", cfg.Port)
	fmt.Printf("  │  Sample Card Endpoint  : /v1/cards/dark_tech.png            │\n")
	fmt.Printf("  │  L1 Memory Cache Cap   : %-35s│\n", fmt.Sprintf("%.1f MB RAM", float64(cfg.L1MaxBytes)/1024/1024))
	fmt.Printf("  │  L2 Disk Store Path    : %-35s│\n", cfg.CacheDir)
	if cfg.Secret != "" {
		fmt.Printf("  │  HMAC URL Signing      : ACTIVE (secret configured)         │\n")
	} else {
		fmt.Printf("  │  HMAC URL Signing      : DISABLED (open public endpoint)    │\n")
	}
	fmt.Println("  └─────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-shutdownChan
	fmt.Println("\nReceived shutdown signal, stopping daemon...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}

	fmt.Println("ogcraft daemon stopped cleanly.")
}
