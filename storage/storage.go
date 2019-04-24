package storage

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	Method string         `json:"method"`
	Value  string         `json:"value,omitempty"`
	Key    string         `json:"key"`
	Error  string         `json:"error,omitempty"`
	Result string         `json:"result,omitempty"`
	Mutex  sync.Mutex     `json:"-"`
	WaitGr sync.WaitGroup `json:"-"`
}

var DbStorage = map[uuid.UUID]string{}

// Set function save data to storage
func (s *Storage) Set() error {
	if len(s.Key) > 16 {
		return fmt.Errorf("length of key too long")
	}
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return err
	}
	s.mutex.Lock()
	DbStorage[key] = s.Value
	s.mutex.Unlock()
	return nil
}

// Get function recieve data from storage
func (s *Storage) Get() (string, error) {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return "", err
	}
	s.mutex.Lock()
	value := DbStorage[key]
	s.mutex.Unlock()
	return value, nil
}

// Delete function remove data from storage
func (s *Storage) Delete() error {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return err
	}
	s.mutex.Lock()
	delete(DbStorage, key)
	s.mutex.Unlock()
	return nil
}

// Exist check is key present in map
func (s *Storage) Exist() bool {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return false
	}
	s.mutex.Lock()
	_, ok := DbStorage[key]
	s.mutex.Unlock()
	return ok
}
