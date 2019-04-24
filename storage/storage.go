package storage

import (
	"github.com/google/uuid"
)

type Storage struct {
	Method string
	Value  []byte
	Key    string
	Error  string
	Result string
}

var DbStorage = map[uuid.UUID]string{}

// Set function save data to storage
func (s *Storage) Set() bool {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return false
	}
	DbStorage[key] = string(s.Value)
	return true
}

// Get function recieve data from storage
func (s *Storage) Get() string {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return ""
	}
	return DbStorage[key]
}

// Delete function remove data from storage
func (s *Storage) Delete() bool {
	key, err := uuid.FromBytes([]byte(s.Key))
	if err != nil {
		return false
	}
	delete(DbStorage, key)
	return true
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

// type storageError struct {
// 	errType string
// 	err     error
// }

// func (e *storageError) Error() string {
// 	return fmt.Sprintf("[%s : %s ]", e.errType, e.err.Error())
// }
