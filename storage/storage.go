package storage

var (
	ErrNotFound     = "value not found"
	ErrKeyTooLong   = "key too long"
	ErrKeyNotSet    = "key not set"
	ErrValueTooLong = "value too long"
	Success         = "success"
)

type Storage struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

var DbStorage = map[string]string{}

// Set function save data to storage
func (s *Storage) Set() *Storage {
	errStr := CheckKeyLenght(s.Key)
	if errStr != "" {
		s.Value = ""
		s.Error = errStr
		return s
	}
	DbStorage[s.Key] = s.Value
	s.Value = ""
	s.Result = Success
	return s
}

// Get function recieve data from storage
func (s *Storage) Get() *Storage {
	errStr := CheckKeyLenght(s.Key)
	if errStr != "" {
		s.Value = ""
		s.Error = errStr
		return s
	}
	if _, ok := DbStorage[s.Key]; ok {
		s.Value = DbStorage[s.Key]
		return s
	}
	s.Error = ErrNotFound
	return s
}

// Delete function remove data from storage
func (s *Storage) Delete() *Storage {
	errStr := CheckKeyLenght(s.Key)
	if errStr != "" {
		s.Error = errStr
		return s
	}
	if _, ok := DbStorage[s.Key]; ok {
		delete(DbStorage, s.Key)
		s.Result = Success
		return s
	}
	s.Error = ErrNotFound
	return s
}

// Exist check is key present in map
func (s *Storage) Exist() *Storage {
	errStr := CheckKeyLenght(s.Key)
	if errStr != "" {
		s.Error = errStr
		return s
	}
	_, ok := DbStorage[s.Key]
	if !ok {
		s.Error = ErrNotFound
		return s
	}
	s.Result = Success
	return s
}

// CheckKeyLenght function check length of incoming key
func CheckKeyLenght(key string) string {
	if len(key) > 16 {
		return ErrKeyTooLong
	}
	if len(key) == 0 {
		return ErrKeyNotSet
	}
	return ""
}

// CheckDbLenght return value of keys in DbStorage
func (s *Storage) CheckDbLenght() int {
	return len(DbStorage)
}

// CheckValueLength define value lenght and if it more than 512 byte return error
func (s *Storage) CheckValueLength(length int) bool {
	return len(s.Value) > length
}
