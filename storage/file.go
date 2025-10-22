package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	filepath string
	data     map[string]string
	mu       sync.RWMutex
}

func NewFileStorage(filepath string) (*FileStorage, error) {
	fs := &FileStorage{
		filepath: filepath,
		data:     make(map[string]string),
	}
	
	if err := fs.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	
	return fs, nil
}

func (s *FileStorage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &s.data)
}

func (s *FileStorage) save() error {
	dir := filepath.Dir(s.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(s.filepath, data, 0644)
}

func (s *FileStorage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	value, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}
	
	return value, nil
}

func (s *FileStorage) Set(key string, value string) error {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
	
	return s.save()
}

func (s *FileStorage) Delete(key string) error {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
	
	return s.save()
}
