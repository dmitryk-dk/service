package storage

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

type Storage struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

var DbStorage = map[string]string{}

// Set function save data to storage
func (s *Storage) Set() string {
	var errStr errorString
	errStr.s = CheckKeyLenght(s.Key)
	if errStr.s != "" {
		return errStr.Error()
	}
	DbStorage[s.Key] = s.Value
	return ""
}

// Get function recieve data from storage
func (s *Storage) Get() (bool, string) {
	var errStr errorString
	errStr.s = CheckKeyLenght(s.Key)
	if errStr.s != "" {
		return false, errStr.Error()
	}
	if _, ok := DbStorage[s.Key]; ok {
		return true, DbStorage[s.Key]
	}
	return false, ""
}

// Delete function remove data from storage
func (s *Storage) Delete() (bool, string) {
	var errStr errorString
	errStr.s = CheckKeyLenght(s.Key)
	if errStr.s != "" {
		return false, errStr.Error()
	}
	if _, ok := DbStorage[s.Key]; ok {
		delete(DbStorage, s.Key)
		return true, "success"
	}
	errStr.s = "value not found"
	return false, errStr.Error()
}

// Exist check is key present in map
func (s *Storage) Exist() (bool, string) {
	var errStr errorString
	errStr.s = CheckKeyLenght(s.Key)
	if errStr.s != "" {
		return false, errStr.Error()
	}
	_, ok := DbStorage[s.Key]
	if !ok {
		errStr.s = "value not found"
		return false, errStr.Error()
	}
	return ok, ""
}

func CheckKeyLenght(key string) string {
	if len(key) > 16 {
		return "key too long"
	}
	if len(key) == 0 {
		return "key not set"
	}
	return ""
}

// CheckDbLenght return value of keys in DbStorage
func CheckDbLenght() int {
	return len(DbStorage)
}
