package api_test

import (
	"bytes"
	"encoding/json"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexandrmotologa/ogcraft/internal/api"
	"github.com/alexandrmotologa/ogcraft/internal/cache"
	"github.com/alexandrmotologa/ogcraft/internal/config"
	"github.com/alexandrmotologa/ogcraft/internal/fonts"
	"github.com/alexandrmotologa/ogcraft/internal/templates"
)

func setupTestServer(t *testing.T, secret string) (*api.Server, string) {
	_ = fonts.Init()
	tempDir := filepath.Join(os.TempDir(), "ogcraft_api_test")
	_ = os.RemoveAll(tempDir)

	customDir := filepath.Join(tempDir, "custom_templates")
	cacheDir := filepath.Join(tempDir, "cache")

	cfg := config.DefaultConfig()
	cfg.Secret = secret
	cfg.CacheDir = cacheDir
	cfg.TemplatesDir = customDir

	tmplMgr, err := templates.NewManager(customDir)
	if err != nil {
		t.Fatalf("failed to create template manager: %v", err)
	}

	cacheMgr, err := cache.NewTwoTierCache(cfg.L1MaxBytes, cacheDir)
	if err != nil {
		t.Fatalf("failed to create cache manager: %v", err)
	}

	server := api.NewServer(cfg, tmplMgr, cacheMgr, nil)
	return server, tempDir
}

func TestHealthEndpoint(t *testing.T) {
	server, tempDir := setupTestServer(t, "")
	defer os.RemoveAll(tempDir)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if res["status"] != "healthy" {
		t.Fatalf("expected status healthy, got %v", res["status"])
	}
}

func TestRenderCardWithTwoTierCacheAndETag(t *testing.T) {
	server, tempDir := setupTestServer(t, "")
	defer os.RemoveAll(tempDir)

	cardURL := "/v1/cards/dark_tech.png?title=Testing+Two-Tier+Cache&author=Alexandr"

	// 1. Initial Request (MISS)
	req1 := httptest.NewRequest("GET", cardURL, nil)
	rec1 := httptest.NewRecorder()
	server.Router().ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("initial request failed with code %d: %s", rec1.Code, rec1.Body.String())
	}
	if xCache := rec1.Header().Get("X-Cache"); xCache != "MISS" {
		t.Errorf("expected X-Cache MISS, got %s", xCache)
	}
	etag := rec1.Header().Get("ETag")
	if etag == "" {
		t.Fatal("expected non-empty ETag header")
	}

	// Verify valid PNG
	_, err := png.DecodeConfig(bytes.NewReader(rec1.Body.Bytes()))
	if err != nil {
		t.Fatalf("invalid PNG returned: %v", err)
	}

	// 2. Repeat Request (HIT-L1)
	req2 := httptest.NewRequest("GET", cardURL, nil)
	rec2 := httptest.NewRecorder()
	server.Router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second request failed with code %d", rec2.Code)
	}
	if xCache := rec2.Header().Get("X-Cache"); xCache != "HIT-L1" {
		t.Errorf("expected X-Cache HIT-L1, got %s", xCache)
	}

	// 3. RFC 9111 ETag 304 Not Modified
	req3 := httptest.NewRequest("GET", cardURL, nil)
	req3.Header.Set("If-None-Match", etag)
	rec3 := httptest.NewRecorder()
	server.Router().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusNotModified {
		t.Fatalf("expected 304 Not Modified, got %d", rec3.Code)
	}
}

func TestHMACSignatureVerification(t *testing.T) {
	secret := "super-secret-key-1234"
	server, tempDir := setupTestServer(t, secret)
	defer os.RemoveAll(tempDir)

	template := "minimal_light"
	params := map[string]string{
		"title":  "Signed Request Card",
		"author": "Motologa",
	}

	// 1. Request without signature should fail 401
	unsignedURL := "/v1/cards/minimal_light.png?title=Signed+Request+Card&author=Motologa"
	req1 := httptest.NewRequest("GET", unsignedURL, nil)
	rec1 := httptest.NewRecorder()
	server.Router().ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for unsigned request, got %d", rec1.Code)
	}

	// 2. Request with valid signature
	sig := api.GenerateSignature(secret, template, params)
	signedURL := unsignedURL + "&sig=" + sig
	req2 := httptest.NewRequest("GET", signedURL, nil)
	rec2 := httptest.NewRecorder()
	server.Router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for valid signature, got %d: %s", rec2.Code, rec2.Body.String())
	}
}

func TestTemplatesCRUD(t *testing.T) {
	server, tempDir := setupTestServer(t, "")
	defer os.RemoveAll(tempDir)

	// List builtin
	reqList := httptest.NewRequest("GET", "/v1/templates", nil)
	recList := httptest.NewRecorder()
	server.Router().ServeHTTP(recList, reqList)

	var list []templates.TemplateSummary
	_ = json.NewDecoder(recList.Body).Decode(&list)
	if len(list) < 6 {
		t.Fatalf("expected at least 6 templates, got %d", len(list))
	}

	// Create custom template
	custom := templates.Template{
		ID:          "custom_brand",
		Name:        "Custom Brand Card",
		Description: "A test brand template",
		Width:       1200,
		Height:      630,
		Background: templates.BackgroundSpec{
			Type:  "solid",
			Color: "#1e1b4b",
		},
		Elements: []templates.Element{
			{
				Type:     "text",
				Field:    "title",
				Default:  "Custom Title",
				Font:     "Inter-Bold",
				FontSize: 48,
				Color:    "#ffffff",
				X:        100,
				Y:        200,
			},
		},
	}

	body, _ := json.Marshal(custom)
	reqCreate := httptest.NewRequest("POST", "/v1/templates", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	recCreate := httptest.NewRecorder()
	server.Router().ServeHTTP(recCreate, reqCreate)

	if recCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", recCreate.Code, recCreate.Body.String())
	}

	// Verify custom template can render
	reqRender := httptest.NewRequest("GET", "/v1/cards/custom_brand.png?title=Hello+Custom", nil)
	recRender := httptest.NewRecorder()
	server.Router().ServeHTTP(recRender, reqRender)

	if recRender.Code != http.StatusOK {
		t.Fatalf("rendering custom template failed: %d", recRender.Code)
	}

	// Delete custom template
	reqDel := httptest.NewRequest("DELETE", "/v1/templates/custom_brand", nil)
	recDel := httptest.NewRecorder()
	server.Router().ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusOK {
		t.Fatalf("expected 200 OK deleting custom template, got %d", recDel.Code)
	}
}
