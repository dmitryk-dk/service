package storage

import (
	"errors"
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

var DbStorage = map[string]string{}

// SetDefaultValues set values to default state
func (s *Storage) SetDefaultValues() {
	s.Error = ""
	s.Result = ""
}

// Set function save data to storage
func (s *Storage) Set() error {
	s.SetDefaultValues()
	err := CheckDbLenght()
	if err != nil {
		return err
	}
	err = CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Error = err.Error()
		return err
	}
	DbStorage[s.Key] = s.Value
	s.Result = Success
	return nil
}

// Get function recieve data from storage
func (s *Storage) Get() error {
	s.SetDefaultValues()
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Error = err.Error()
		return err
	}
	if _, ok := DbStorage[s.Key]; ok {
		s.Value = DbStorage[s.Key]
	} else {
		s.Error = ErrNotFound.Error()
	}
	return nil
}

// Delete function remove data from storage
func (s *Storage) Delete() error {
	s.SetDefaultValues()
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Error = err.Error()
		return err
	}
	if _, ok := DbStorage[s.Key]; ok {
		delete(DbStorage, s.Key)
		s.Result = Success
		return nil
	}
	s.Error = ErrNotFound.Error()
	return nil
}

// Exist check is key present in map
func (s *Storage) Exist() error {
	s.SetDefaultValues()
	err := CheckDbErrors(s.Key, s.Value)
	if err != nil {
		s.Error = err.Error()
		return err
	}
	if _, ok := DbStorage[s.Key]; ok {
		s.Result = Success
	} else {
		s.Error = ErrNotFound.Error()
	}
	return nil
}

// RequestMethod return method which income in request
func (s Storage) RequestMethod() string {
	return s.Method
}

// CheckDbLenght return value of keys in DbStorage
func CheckDbLenght() error {
	if len(DbStorage) >= 1024 {
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
