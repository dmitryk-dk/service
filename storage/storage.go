package storage

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("value not found")
	ErrKeyTooLong   = errors.New("key too long")
	ErrKeyNotSet    = errors.New("key not set")
	ErrValueTooLong = errors.New("value too long")
	ErrDbLentgh     = errors.New("Database is full")
	Success         = "success"
)

// Storage describe object of data for work with db
type Storage struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

type errorString struct {
	s string
}

func (e errorString) String() string {
	return e.s
}

var DbStorage = map[string]string{}

// Set function save data to storage
func (s *Storage) Set() *Storage {
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Value = ""
		s.Error = err.Error()
		return s
	}
	DbStorage[s.Key] = s.Value
	s.Value = ""
	s.Result = Success
	return s
}

// Get function recieve data from storage
func (s *Storage) Get() *Storage {
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Value = ""
		s.Error = err.Error()
		return s
	}
	if _, ok := DbStorage[s.Key]; ok {
		s.Value = DbStorage[s.Key]
	} else {
		s.Error = ErrNotFound.Error()
	}
	return s
}

// Delete function remove data from storage
func (s *Storage) Delete() *Storage {
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Value = ""
		s.Error = err.Error()
		return s
	}
	if _, ok := DbStorage[s.Key]; ok {
		fmt.Printf("is ok ->> %v", ok)
		delete(DbStorage, s.Key)
		s.Result = Success
		return s
	}
	s.Error = ErrNotFound.Error()
	fmt.Printf("delete ->> %v \n", s)
	return s
}

// Exist check is key present in map
func (s *Storage) Exist() *Storage {
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Value = ""
		s.Error = err.Error()
		return s
	}
	if _, ok := DbStorage[s.Key]; ok {
		s.Result = Success
	} else {
		s.Error = ErrNotFound.Error()
	}
	return s
}

// RequestMethod return method which income in request
func (s Storage) RequestMethod() string {
	return s.Method
}

// CheckDbLenght return value of keys in DbStorage
func (s Storage) CheckDbLenght() error {
	if len(DbStorage) > 1024 {
		return ErrDbLentgh
	}
	return nil
}

// CheckKeyLenght function check length of incoming key
func CheckKeyLenght(key string) error {
	if len(key) > 16 {
		return ErrKeyTooLong
	}
	if len(key) == 0 {
		return ErrKeyNotSet
	}
	return nil
}

// CheckValueLength define storage value length and if it more than 512 byte return error
func CheckValueLength(value string, length int) error {
	if len(value) > length {
		return ErrValueTooLong
	}
	return nil
}

// CheckDbErrors check income key and value for length and return error if
// key or value has incorrect length
func CheckDbErrors(key, value string) error {
	err := CheckKeyLenght(key)
	if err != nil {
		return err
	}
	err = CheckValueLength(value, 512)
	if err != nil {
		return err
	}
	return nil
}
