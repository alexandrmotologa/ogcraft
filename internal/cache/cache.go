package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// CacheSource identifies where an image was retrieved from.
const (
	SourceL1   = "HIT-L1"
	SourceL2   = "HIT-L2"
	SourceMiss = "MISS"
)

// CacheResult contains retrieved cache data and metadata.
type CacheResult struct {
	Data   []byte
	Source string
	ETag   string
	Key    string
}

// TwoTierCache orchestrates L1 in-memory LRU and L2 persistent disk caching.
type TwoTierCache struct {
	L1 *L1MemoryCache
	L2 *L2DiskCache
}

// NewTwoTierCache initializes the dual caching system.
func NewTwoTierCache(l1MaxBytes int64, l2Dir string) (*TwoTierCache, error) {
	l1 := NewL1MemoryCache(l1MaxBytes)
	l2, err := NewL2DiskCache(l2Dir)
	if err != nil {
		return nil, err
	}

	return &TwoTierCache{
		L1: l1,
		L2: l2,
	}, nil
}

// ComputeKey returns a deterministic SHA-256 cache key and RFC 9111 ETag.
func ComputeKey(identifier string, params map[string]string) (key string, etag string) {
	var keys []string
	for k := range params {
		if k == "sig" || k == "_refresh" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(identifier)
	sb.WriteString("?")
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
		sb.WriteString("&")
	}

	h := sha256.Sum256([]byte(sb.String()))
	key = hex.EncodeToString(h[:])
	etag = fmt.Sprintf("W/\"%s\"", key[:16])
	return key, etag
}

// Get checks L1 memory first, falling back to L2 disk.
func (c *TwoTierCache) Get(key string) (*CacheResult, bool) {
	etag := fmt.Sprintf("W/\"%s\"", key[:16])

	// 1. Check L1 Memory LRU
	if data, hit := c.L1.Get(key); hit {
		return &CacheResult{
			Data:   data,
			Source: SourceL1,
			ETag:   etag,
			Key:    key,
		}, true
	}

	// 2. Check L2 Persistent Disk
	if data, hit := c.L2.Get(key); hit {
		c.L1.Set(key, data)
		return &CacheResult{
			Data:   data,
			Source: SourceL2,
			ETag:   etag,
			Key:    key,
		}, true
	}

	return nil, false
}

// Set stores image data into both L1 and L2 concurrently.
func (c *TwoTierCache) Set(key string, data []byte) {
	c.L1.Set(key, data)
	_ = c.L2.Set(key, data)
}

// Clear purges all cached assets from both memory and disk.
func (c *TwoTierCache) Clear() error {
	c.L1.Clear()
	return c.L2.Clear()
}
