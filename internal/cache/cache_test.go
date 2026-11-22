package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexandrmotologa/ogcraft/internal/cache"
)

func TestTwoTierCacheWorkflow(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "ogcraft_cache_test")
	_ = os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	// Set small 1KB L1 limit to test eviction
	c, err := cache.NewTwoTierCache(1024, tempDir)
	if err != nil {
		t.Fatalf("failed to create two tier cache: %v", err)
	}

	key1, etag1 := cache.ComputeKey("dark_tech", map[string]string{"title": "Card 1"})
	key2, _ := cache.ComputeKey("dark_tech", map[string]string{"title": "Card 2"})

	if etag1 == "" {
		t.Fatal("expected non-empty etag")
	}

	// 1. Initial Get should MISS
	_, hit := c.Get(key1)
	if hit {
		t.Fatal("expected cache MISS")
	}

	// 2. Put 600 bytes
	data1 := make([]byte, 600)
	data1[0] = 0xAA
	c.Set(key1, data1)

	// 3. Get should HIT-L1
	res1, hit := c.Get(key1)
	if !hit || res1.Source != cache.SourceL1 {
		t.Fatalf("expected HIT-L1, got hit=%v source=%s", hit, res1.Source)
	}

	// 4. Put another 600 bytes, exceeding 1KB max in L1, evicting key1 from L1
	data2 := make([]byte, 600)
	data2[0] = 0xBB
	c.Set(key2, data2)

	// Key2 should be in L1
	res2, hit := c.Get(key2)
	if !hit || res2.Source != cache.SourceL1 {
		t.Fatalf("expected key2 HIT-L1, got source=%s", res2.Source)
	}

	// Key1 should now be evicted from L1, but found in L2!
	res1L2, hit := c.Get(key1)
	if !hit || res1L2.Source != cache.SourceL2 {
		t.Fatalf("expected key1 HIT-L2 from disk, got source=%s", res1L2.Source)
	}

	// Once read from L2, key1 is promoted back to L1!
	res1Promoted, hit := c.Get(key1)
	if !hit || res1Promoted.Source != cache.SourceL1 {
		t.Fatalf("expected key1 promoted to HIT-L1, got source=%s", res1Promoted.Source)
	}
}
