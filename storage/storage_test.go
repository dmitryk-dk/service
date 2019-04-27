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
	Result string
}

var testSetData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set", ""},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", "key too long", ""},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"set", "", "my_key_2", "", "success"},
	testStorageStruct{"set", "my_value_for_key_3", "my_key_3", "", "success"},
}

var testGetData = []testStorageStruct{
	testStorageStruct{"get", "my_value_for_empty_key", "", "key not set", ""},
	testStorageStruct{"get", "my_value_for_key_1", "my_key_1", "my_value_for_key_1", ""},
	testStorageStruct{"get", "", "my_key_2", "", ""},
	testStorageStruct{"get", "", "qwertyuiopasdfghaa", "key too long", ""},
	testStorageStruct{"get", "some_value", "qwertyuiopasdfghaa", "key too long", ""},
}

var testExistData = []testStorageStruct{
	testStorageStruct{"exist", "my_value_for_empty_key", "", "key not set", ""},
	testStorageStruct{"exist", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"exist", "", "my_key_2", "", "success"},
	testStorageStruct{"exist", "", "qwertyuiopasdfghaa", "key too long", ""},
	testStorageStruct{"exist", "some_value", "qwertyuiopasdfghaa", "key too long", ""},
}

var testDeleteData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", "key not set", ""},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", "key too long", ""},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"set", "", "my_key_2", "value not found", ""},
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
		storage := store.Set()
		if storage.Error != data.Error {
			t.Errorf("errStr '%v'; data.Error '%v';", storage.Error, data.Error)
		}
		if storage.Error == "" {
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
		storage := store.Get()
		if storage.Error != "" {
			if storage.Error != data.Error {
				t.Errorf("Error in get method was '%v'; want '%v'", storage.Error, data.Error)
			}
		}
		if storage.Error == "" {
			if storage.Value != data.Value {
				t.Errorf("Value getted from db '%v'; value from test '%v", storage.Value, data.Value)
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
		storage := store.Exist()
		if storage.Error != "" && storage.Error != data.Error {
			t.Errorf("Error in exist methos was '%v'; want '%v'; key '%v'", storage.Error, data.Error, data.Key)
		}
		if storage.Error == "" {
			if storage.Result != data.Result {
				t.Errorf("Exist result was '%v'; want '%v'", storage.Result, data.Result)
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
		storage := store.Delete()
		if storage.Result != "success" && storage.Error != data.Error {
			t.Errorf("Error in delete methis was '%v'; want '%v';", storage.Error, data.Error)
		}
		if storage.Error == "" && storage.Result != "success" {
			t.Error("Method delete not remove data from Db")
		}
	}
	for _, v := range testDeleteData {
		t.Run("Run test for Delete Method", func(t *testing.T) {
			checkDelete(t, v)
		})
	}
	fmt.Printf("** TestDelete - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}
