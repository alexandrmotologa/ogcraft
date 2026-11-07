package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/alexandrmotologa/ogcraft/internal/cache"
	"github.com/alexandrmotologa/ogcraft/internal/config"
	"github.com/alexandrmotologa/ogcraft/internal/renderer"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server coordinates HTTP routing, API handlers, caching, and template rendering.
type Server struct {
	cfg          *config.Config
	templates    *templates.Manager
	cache        *cache.TwoTierCache
	router       *chi.Mux
	startTime    time.Time
	rendersTotal atomic.Uint64
}

// RenderPayload represents the JSON body sent to POST /v1/render for live previews.
type RenderPayload struct {
	Template templates.Template `json:"template"`
	Params   map[string]string  `json:"params"`
	Format   string             `json:"format,omitempty"`
}

// NewServer creates a new API Server instance.
func NewServer(cfg *config.Config, tmplMgr *templates.Manager, cacheMgr *cache.TwoTierCache, webFS fs.FS) *Server {
	s := &Server{
		cfg:       cfg,
		templates: tmplMgr,
		cache:     cacheMgr,
		router:    chi.NewRouter(),
		startTime: time.Now(),
	}

	s.setupRoutes(webFS)
	return s
}

// Router returns the chi router.
func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) setupRoutes(webFS fs.FS) {
	r := s.router

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"ETag", "X-Cache", "X-Render-Time"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	rl := newRateLimiter(s.cfg.RateLimit)
	r.Use(RateLimitMiddleware(rl))

	// Health and Prometheus diagnostics
	r.Get("/health", s.handleHealth)
	r.Get("/metrics", s.handlePrometheusMetrics)

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/cards/{template}", s.handleRenderCard)
		r.Post("/render", s.handleLiveRender)

		r.Get("/templates", s.handleListTemplates)
		r.Get("/templates/{id}", s.handleGetTemplate)
		r.Post("/templates", s.handleSaveTemplate)
		r.Delete("/templates/{id}", s.handleDeleteTemplate)

		// Cache maintenance
		r.Post("/cache/purge", s.handleCachePurge)
	})

	// Static frontend
	if webFS != nil {
		fileServer := http.FileServer(http.FS(webFS))
		r.Handle("/*", http.StripPrefix("/", fileServer))
	} else {
		r.Get("/", s.handleFallbackIndex)
	}
}

func (s *Server) handleRenderCard(w http.ResponseWriter, r *http.Request) {
	templateParam := chi.URLParam(r, "template")

	format := "png"
	templateID := templateParam

	switch {
	case strings.HasSuffix(templateParam, ".webp"):
		format = "webp"
		templateID = strings.TrimSuffix(templateParam, ".webp")
	case strings.HasSuffix(templateParam, ".jpeg"):
		format = "jpeg"
		templateID = strings.TrimSuffix(templateParam, ".jpeg")
	case strings.HasSuffix(templateParam, ".jpg"):
		format = "jpeg"
		templateID = strings.TrimSuffix(templateParam, ".jpg")
	case strings.HasSuffix(templateParam, ".png"):
		format = "png"
		templateID = strings.TrimSuffix(templateParam, ".png")
	default:
		if qFormat := r.URL.Query().Get("format"); qFormat != "" {
			format = qFormat
		}
	}

	queryParams := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}

	if s.cfg.Secret != "" {
		sig := queryParams["sig"]
		if !VerifySignature(s.cfg.Secret, templateID, queryParams, sig) {
			http.Error(w, "invalid or missing signature", http.StatusUnauthorized)
			return
		}
	}

	tmpl, err := s.templates.Get(templateID)
	if err != nil {
		http.Error(w, fmt.Sprintf("template not found: %s", templateID), http.StatusNotFound)
		return
	}

	cacheKey, etag := cache.ComputeKey(tmpl.ID+"."+format, queryParams)

	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.WriteHeader(http.StatusNotModified)
		return
	}

	refresh := r.URL.Query().Get("_refresh") == "true"

	if !refresh {
		if cached, hit := s.cache.Get(cacheKey); hit {
			contentType := "image/png"
			if format == "webp" {
				contentType = "image/webp"
			} else if format == "jpeg" || format == "jpg" {
				contentType = "image/jpeg"
			}

			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Content-Length", strconv.Itoa(len(cached.Data)))
			w.Header().Set("ETag", cached.ETag)
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("X-Cache", cached.Source)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(cached.Data)
			return
		}
	}

	start := time.Now()
	imgBytes, contentType, err := renderer.RenderCardFormat(tmpl, queryParams, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("render error: %v", err), http.StatusInternalServerError)
		return
	}
	renderDuration := time.Since(start)
	s.rendersTotal.Add(1)

	s.cache.Set(cacheKey, imgBytes)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Cache", cache.SourceMiss)
	w.Header().Set("X-Render-Time", fmt.Sprintf("%.2fms", float64(renderDuration.Microseconds())/1000.0))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(imgBytes)
}

func (s *Server) handleLiveRender(w http.ResponseWriter, r *http.Request) {
	var payload RenderPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("invalid json payload: %v", err), http.StatusBadRequest)
		return
	}

	format := payload.Format
	if format == "" {
		format = "png"
	}

	start := time.Now()
	imgBytes, contentType, err := renderer.RenderCardFormat(&payload.Template, payload.Params, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("render error: %v", err), http.StatusInternalServerError)
		return
	}
	duration := time.Since(start)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.Header().Set("X-Render-Time", fmt.Sprintf("%.2fms", float64(duration.Microseconds())/1000.0))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(imgBytes)
}

func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	list := s.templates.List()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleGetTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tmpl, err := s.templates.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tmpl)
}

func (s *Server) handleSaveTemplate(w http.ResponseWriter, r *http.Request) {
	var tmpl templates.Template
	if err := json.NewDecoder(r.Body).Decode(&tmpl); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}

	if err := s.templates.Save(tmpl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "saved", "id": tmpl.ID})
}

func (s *Server) handleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.templates.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted", "id": id})
}

func (s *Server) handleCachePurge(w http.ResponseWriter, r *http.Request) {
	if s.cfg.Secret != "" {
		authHeader := r.Header.Get("Authorization")
		expected := "Bearer " + s.cfg.Secret
		if authHeader != expected && r.URL.Query().Get("secret") != s.cfg.Secret {
			http.Error(w, "unauthorized cache purge", http.StatusUnauthorized)
			return
		}
	}

	if err := s.cache.Clear(); err != nil {
		http.Error(w, fmt.Sprintf("cache clear error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "cache cleared successfully",
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	l1Items, l1Bytes, l1Hits, l1Misses, l1Evicts := s.cache.L1.Stats()
	l2Hits, l2Misses, l2Writes := s.cache.L2.Stats()

	l1HitRate := 0.0
	if total := l1Hits + l1Misses; total > 0 {
		l1HitRate = float64(l1Hits) / float64(total)
	}

	diag := map[string]interface{}{
		"status":         "healthy",
		"version":        "1.0.0",
		"uptime_seconds": int(time.Since(s.startTime).Seconds()),
		"memory": map[string]interface{}{
			"alloc_mb":   float64(m.Alloc) / 1024 / 1024,
			"total_mb":   float64(m.TotalAlloc) / 1024 / 1024,
			"sys_mb":     float64(m.Sys) / 1024 / 1024,
			"heap_inuse": float64(m.HeapInuse) / 1024 / 1024,
			"num_gc":     m.NumGC,
		},
		"cache": map[string]interface{}{
			"l1_items":     l1Items,
			"l1_bytes":     l1Bytes,
			"l1_hits":      l1Hits,
			"l1_misses":    l1Misses,
			"l1_hit_rate":  l1HitRate,
			"l1_evictions": l1Evicts,
			"l2_hits":      l2Hits,
			"l2_misses":    l2Misses,
			"l2_writes":    l2Writes,
		},
		"renders_total": s.rendersTotal.Load(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(diag)
}

func (s *Server) handlePrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	l1Items, l1Bytes, l1Hits, l1Misses, l1Evicts := s.cache.L1.Stats()
	l2Hits, l2Misses, l2Writes := s.cache.L2.Stats()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	fmt.Fprintf(w, "# HELP ogcraft_uptime_seconds Total daemon uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE ogcraft_uptime_seconds counter\n")
	fmt.Fprintf(w, "ogcraft_uptime_seconds %d\n\n", int(time.Since(s.startTime).Seconds()))

	fmt.Fprintf(w, "# HELP ogcraft_renders_total Total number of cards rendered\n")
	fmt.Fprintf(w, "# TYPE ogcraft_renders_total counter\n")
	fmt.Fprintf(w, "ogcraft_renders_total %d\n\n", s.rendersTotal.Load())

	fmt.Fprintf(w, "# HELP ogcraft_cache_hits_total Cache hit counter by tier\n")
	fmt.Fprintf(w, "# TYPE ogcraft_cache_hits_total counter\n")
	fmt.Fprintf(w, "ogcraft_cache_hits_total{tier=\"l1\"} %d\n", l1Hits)
	fmt.Fprintf(w, "ogcraft_cache_hits_total{tier=\"l2\"} %d\n\n", l2Hits)

	fmt.Fprintf(w, "# HELP ogcraft_cache_misses_total Cache misses by tier\n")
	fmt.Fprintf(w, "# TYPE ogcraft_cache_misses_total counter\n")
	fmt.Fprintf(w, "ogcraft_cache_misses_total{tier=\"l1\"} %d\n", l1Misses)
	fmt.Fprintf(w, "ogcraft_cache_misses_total{tier=\"l2\"} %d\n\n", l2Misses)

	fmt.Fprintf(w, "# HELP ogcraft_cache_bytes Current byte usage of L1 memory cache\n")
	fmt.Fprintf(w, "# TYPE ogcraft_cache_bytes gauge\n")
	fmt.Fprintf(w, "ogcraft_cache_bytes{tier=\"l1\"} %d\n\n", l1Bytes)

	fmt.Fprintf(w, "# HELP ogcraft_cache_items Number of items cached\n")
	fmt.Fprintf(w, "# TYPE ogcraft_cache_items gauge\n")
	fmt.Fprintf(w, "ogcraft_cache_items{tier=\"l1\"} %d\n", l1Items)
	fmt.Fprintf(w, "ogcraft_cache_evictions_total{tier=\"l1\"} %d\n", l1Evicts)
	fmt.Fprintf(w, "ogcraft_cache_writes_total{tier=\"l2\"} %d\n\n", l2Writes)

	fmt.Fprintf(w, "# HELP ogcraft_memory_alloc_bytes Memory bytes allocated and currently in use\n")
	fmt.Fprintf(w, "# TYPE ogcraft_memory_alloc_bytes gauge\n")
	fmt.Fprintf(w, "ogcraft_memory_alloc_bytes %d\n", m.Alloc)
	fmt.Fprintf(w, "ogcraft_memory_sys_bytes %d\n", m.Sys)
	fmt.Fprintf(w, "ogcraft_goroutines %d\n", runtime.NumGoroutine())
}

func (s *Server) handleFallbackIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>ogcraft daemon</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #090d16; color: #f0f6fc; margin: 0; padding: 40px; display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 80vh; }
		.card { background: #121a29; border: 1px solid #212d45; border-radius: 12px; padding: 32px; max-width: 640px; width: 100%; box-shadow: 0 10px 25px rgba(0,0,0,0.4); }
		h1 { margin-top: 0; color: #58a6ff; font-size: 28px; }
		p { color: #8b949e; line-height: 1.5; font-size: 15px; }
		code { background: #1c2438; color: #79c0ff; padding: 3px 6px; border-radius: 4px; font-family: monospace; font-size: 14px; }
		a { color: #58a6ff; text-decoration: none; }
		a:hover { text-decoration: underline; }
		.btn { display: inline-block; background: #238636; color: #fff; padding: 10px 18px; border-radius: 6px; font-weight: 600; margin-top: 16px; }
	</style>
</head>
<body>
	<div class="card">
		<h1>ogcraft Daemon Active</h1>
		<p>Self-hosted real-time OpenGraph preview card generator in Go and React 19.</p>
		<p>Test rendering an image directly:</p>
		<p><a class="btn" target="_blank" href="/v1/cards/dark_tech.png?title=Sub-15ms+Social+Cards&author=Alexandr+Motologa&tag=Go">Open Test Card (dark_tech.png)</a></p>
		<p>Check health diagnostics: <a href="/health"><code>/health</code></a></p>
		<p>View registered templates: <a href="/v1/templates"><code>/v1/templates</code></a></p>
		<p>Prometheus metrics: <a href="/metrics"><code>/metrics</code></a></p>
	</div>
</body>
</html>`
	_, _ = w.Write([]byte(html))
}
