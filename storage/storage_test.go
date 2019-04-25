package storage

import (
	"fmt"
	"strings"
	"testing"
)

type testStorageStruct struct {
	Method string
	Value  string
	Key    string
	Error  string
}

var testSetData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set"},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", "key too long"},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", ""},
	testStorageStruct{"set", "", "my_key_2", ""},
}

var testGetData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set"},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", "my_value_for_key_1"},
	testStorageStruct{"set", "", "my_key_2", "value not found"},
	testStorageStruct{"set", "", "qwertyuiopasdfghaa", "key too long"},
	testStorageStruct{"set", "some_value", "qwertyuiopasdfghaa", "key too long"},
}

var testExistData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set"},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", ""},
	testStorageStruct{"set", "", "my_key_2", "value not found"},
	testStorageStruct{"set", "", "qwertyuiopasdfghaa", "key too long"},
	//testStorageStruct{"set", "some_value", "qwertyuiopasdfghaa", "key too long"},
}

var testDeleteData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set"},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", "key too long"},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", ""},
	testStorageStruct{"set", "", "my_key_2", "value not found"},
}

func ErrorContains(out string, want string) bool {
	if out == "" {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out, want)
}
func TestSet(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		errStr := store.Set()
		if errStr != data.Error {
			t.Errorf("errStr '%v'; data.Error '%v';", errStr, data.Error)
		}
		if errStr == "" {
			if DbStorage[data.Key] != data.Value {
				t.Errorf("DbStorage value '%v'; test data value '%v'", DbStorage[data.Key], data.Value)
			}
		}
	}
	for _, v := range testSetData {
		t.Run("Run test for Set method", func(t *testing.T) {
			checkSet(t, v)
		})
	}
	t.Run("run test for check DBLength", func(t *testing.T) {
		if CheckDbLenght() != 2 {
			t.Errorf("Db length was: '%v'; want: '%d'", CheckDbLenght(), 2)
		}
	})
	fmt.Printf("** TestSet - ALL PASSED (number of test cases: %d)**\n", len(testSetData))
}

func TestGet(t *testing.T) {
	checkGet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		_, value := store.Get()
		if value != "" {
			if !ErrorContains(value, data.Error) {
				t.Errorf("%s", value)
			}
		}
		if value == "" {
			if value != data.Value {
				t.Errorf("Value getted from db '%v'; value from test '%v", value, data.Value)
			}
		}
	}
	for _, v := range testGetData {
		t.Run("Run test for Get method", func(t *testing.T) {
			checkGet(t, v)
		})
	}
	fmt.Printf("** TestGet - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}

func TestExist(t *testing.T) {
	checkExist := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		ok, errStr := store.Exist()
		if errStr != "" && errStr != data.Error {
			t.Errorf("errStr '%v'; data.Error '%v'; key '%v'", errStr, data.Error, data.Key)
		}
		if errStr == "" {
			if !ok {
				t.Errorf("Value not found in db check method exist '%v'", ok)
			}
		}
	}
	for _, v := range testExistData {
		t.Run("Run tests for Exist Method", func(t *testing.T) {
			checkExist(t, v)
		})
	}
	fmt.Printf("** TestExist - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}

func TestDelete(t *testing.T) {
	checkDelete := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		_, errStr := store.Delete()
		if errStr != "success" && errStr != data.Error {
			t.Errorf("errStr '%v'; data.Error '%v';", errStr, data.Error)
		}
		if errStr == "" {
			if _, ok := DbStorage[data.Key]; ok {
				t.Error("Method delete not remove data from Db")
			}
		}
	}
	for _, v := range testDeleteData {
		t.Run("Run test for Delete Method", func(t *testing.T) {
			checkDelete(t, v)
		})
	}
	fmt.Printf("** TestDelete - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}
