package ui

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageCache struct {
	mu      sync.RWMutex
	cache   map[string]*ebiten.Image
	loading map[string]bool
	client  *http.Client
}

func newImageCache() *imageCache {
	return &imageCache{
		cache:   make(map[string]*ebiten.Image),
		loading: make(map[string]bool),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Get devuelve la imagen si ya está en cache
func (c *imageCache) Get(url string) (*ebiten.Image, bool) {
	if url == "" {
		return nil, false
	}
	c.mu.RLock()
	im, ok := c.cache[url]
	c.mu.RUnlock()
	return im, ok
}

// Load descarga la URL en background si no está en cache
func (c *imageCache) Load(url string) {
	if url == "" {
		return
	}

	c.mu.Lock()
	if c.loading[url] || c.cache[url] != nil {
		c.mu.Unlock()
		return
	}
	c.loading[url] = true
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			delete(c.loading, url)
			c.mu.Unlock()
		}()

		resp, err := c.client.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// leer contenido
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, resp.Body); err != nil {
			return
		}

		img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return
		}

		ebImg := ebiten.NewImageFromImage(img)

		c.mu.Lock()
		c.cache[url] = ebImg
		c.mu.Unlock()
	}()
}
