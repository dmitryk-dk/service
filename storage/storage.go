package storage

import (
	"fmt"

	"github.com/google/uuid"
)

type Storage struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
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
	DbStorage[key] = s.Value
	return nil
}

// Get function recieve data from storage
func (s *Storage) Get() (string, error) {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return "", err
	}
	value := DbStorage[key]
	return value, nil
}

// Delete function remove data from storage
func (s *Storage) Delete() error {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return err
	}
	if _, ok := DbStorage[key]; ok {
		delete(DbStorage, key)
	} else {
		return fmt.Errorf("value not found")
	}
	return fmt.Errorf("value not deleted")
}

// Exist check is key present in map
func (s *Storage) Exist() bool {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return false
	}
	_, ok := DbStorage[key]
	return ok
}

func CheckDbLenght() int {
	return len(DbStorage)
}
