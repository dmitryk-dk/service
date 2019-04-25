package storage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type testStorageStruct struct {
	Method string
	Value  string
	Key    string
	Error  string
}

var testSetData = []testStorageStruct{
	{"set", "some_value", "qwert", "invalid UUID (got 5 bytes)"},
	{"set", "some_value", "", "invalid UUID (got 0 bytes)"},
	{"set", "", "qwertyuiopasdfga", ""},
	{"set", "some_value", "qwertyuiopasdfgh", ""},
	{"set", "", "qwertyuiopasdfghaa", "length of key too long"},
	{"set", "some_value", "qwertyuiopasdfghaa", "length of key too long"},
}

var testGetData = []testStorageStruct{
	{"set", "some_value", "qwert", "invalid UUID (got 5 bytes)"},
	{"set", "some_value", "", "invalid UUID (got 0 bytes)"},
	{"set", "", "qwertyuiopasdfga", ""},
	{"set", "some_value", "qwertyuiopasdfgh", ""},
	{"set", "", "qwertyuiopasdfghaa", "invalid UUID (got 18 bytes)"},
	{"set", "some_value", "qwertyuiopasdfghaa", "invalid UUID (got 18 bytes)"},
}

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
func TestSet(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		err := store.Set()
		if !ErrorContains(err, data.Error) {
			t.Error("Invalid keys length set")
		}
		if err == nil {
			key, _ := uuid.FromBytes([]byte(data.Key))
			if DbStorage[key] != data.Value {
				t.Error("Data not set in DbStorage")
			}
		}
	}
	for _, v := range testSetData {
		t.Run("set", func(t *testing.T) {
			checkSet(t, v)
		})
	}
	if CheckDbLenght() != 2 {
		t.Error("Db length not equal")
	}
	fmt.Printf("** TestSet - ALL PASSED (number of test cases: %d)**\n", len(testSetData))
}

func TestGet(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		value, err := store.Get()
		if err != nil {
			if !ErrorContains(err, data.Error) {
				t.Errorf("%s", err)
			}
		}
		if err == nil {
			if value != data.Value {
				t.Error("Value from db not equal")
			}
		}
	}
	for _, v := range testGetData {
		t.Run("set", func(t *testing.T) {
			checkSet(t, v)
		})
	}
	fmt.Printf("** TestGet - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}

func TestExist(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		_, err := store.Exist()
		if err != nil {
			if !ErrorContains(err, data.Error) {
				t.Errorf("%s", err)
			}
		}
	}
	for _, v := range testGetData {
		t.Run("set", func(t *testing.T) {
			checkSet(t, v)
		})
	}
	fmt.Printf("** TestDelete - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}

func TestDelete(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		err := store.Delete()
		if err != nil {
			if !ErrorContains(err, data.Error) {
				t.Errorf("%s", err)
			}
		}
		if err == nil {
			key, _ := uuid.FromBytes([]byte(data.Key))
			if _, ok := DbStorage[key]; ok {
				t.Error("Method delete not remove data from Db")
			}
		}
	}
	for _, v := range testGetData {
		t.Run("set", func(t *testing.T) {
			checkSet(t, v)
		})
	}
	fmt.Printf("** TestDelete - ALL PASSED (number of test cases: %d)**\n", len(testGetData))
}
