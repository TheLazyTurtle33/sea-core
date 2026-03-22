package store

import (
	"encoding/json"
	"os"
	"sync"
)

const dataPath = "/app/data/store.json"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

var instance *Store

func Get() *Store {
	if instance == nil {
		instance = &Store{data: make(map[string]string)}
		instance.load()
	}
	return instance
}

func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return s.save()
}

func (s *Store) GetKey(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) load() {
	f, err := os.ReadFile(dataPath)
	if err != nil {
		return // file doesn't exist yet, start fresh
	}
	json.Unmarshal(f, &s.data)
}

func (s *Store) save() error {
	f, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataPath, f, 0644)
}
