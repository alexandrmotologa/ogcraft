package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// L2DiskCache manages persistent, content-addressed disk caching.
type L2DiskCache struct {
	baseDir string
	hits    atomic.Uint64
	misses  atomic.Uint64
	writes  atomic.Uint64
}

// NewL2DiskCache creates a disk cache manager at the specified directory.
func NewL2DiskCache(baseDir string) (*L2DiskCache, error) {
	if baseDir == "" {
		baseDir = "./data/cache"
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory %s: %w", baseDir, err)
	}

	return &L2DiskCache{
		baseDir: baseDir,
	}, nil
}

func (d *L2DiskCache) getFilePath(key string) string {
	if len(key) >= 4 {
		sub1 := key[:2]
		sub2 := key[2:4]
		return filepath.Join(d.baseDir, sub1, sub2, key+".cache")
	}
	return filepath.Join(d.baseDir, key+".cache")
}

// Get reads cached image bytes from disk.
func (d *L2DiskCache) Get(key string) ([]byte, bool) {
	filePath := d.getFilePath(key)

	data, err := os.ReadFile(filePath)
	if err != nil {
		d.misses.Add(1)
		return nil, false
	}

	d.hits.Add(1)
	return data, true
}

// Set writes image bytes to disk atomically via temp file and rename.
func (d *L2DiskCache) Set(key string, data []byte) error {
	filePath := d.getFilePath(key)
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tempPath := fmt.Sprintf("%s.tmp.%d", filePath, time.Now().UnixNano())
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	d.writes.Add(1)
	return nil
}

// Clear removes all cached files from disk.
func (d *L2DiskCache) Clear() error {
	entries, err := os.ReadDir(d.baseDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		_ = os.RemoveAll(filepath.Join(d.baseDir, entry.Name()))
	}
	return nil
}

// Stats returns operational metrics of the L2 disk cache.
func (d *L2DiskCache) Stats() (hits uint64, misses uint64, writes uint64) {
	return d.hits.Load(), d.misses.Load(), d.writes.Load()
}
