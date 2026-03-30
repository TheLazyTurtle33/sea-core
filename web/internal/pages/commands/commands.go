package commands

import (
	_ "embed"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

// --- Cache ---

type cache struct {
	mu        sync.Mutex
	data      []byte
	fetchedAt time.Time
	ttl       time.Duration
}

var commandCache = &cache{ttl: 5 * time.Minute}

func (c *cache) get() ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil || time.Since(c.fetchedAt) > c.ttl {
		return nil, false
	}
	return c.data, true
}

func (c *cache) set(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = data
	c.fetchedAt = time.Now()
}

// --- Handler ---

//go:embed page.html
var pageHTML []byte

func Page(w http.ResponseWriter, r *http.Request) {
	commandsJSON, ok := commandCache.get()
	if !ok {
		var err error
		commandsJSON, err = botclient.GetCommandsJson()
		if err != nil {
			log.Printf("Error getting commands JSON: %v", err)
			commandsJSON = []byte("[]")
		} else {
			commandCache.set(commandsJSON)
		}
	}

	log.Println(r)

	if r.URL.Query().Get("json") == "true" {
		w.Header().Set("Content-Type", "application/json")
		w.Write(commandsJSON)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(pageHTML)
}
