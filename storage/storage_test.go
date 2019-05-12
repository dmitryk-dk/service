package storage

import (
	"fmt"
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
	testStorageStruct{"set", "my_value_for_empty_key", "", ErrKeyNotSet.Error(), ""},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", ErrKeyTooLong.Error(), ""},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"set", "", "my_key_2", "", "success"},
	testStorageStruct{"set", "my_value_for_key_3", "my_key_3", "", "success"},
}

var testGetData = []testStorageStruct{
	testStorageStruct{"get", "my_value_for_empty_key", "", ErrKeyNotSet.Error(), ""},
	testStorageStruct{"get", "my_value_for_key_1", "my_key_1", "", "my_value_for_key_1"},
	testStorageStruct{"get", "", "my_key_2", "", ""},
	testStorageStruct{"get", "", "qwertyuiopasdfghaa", ErrKeyTooLong.Error(), ""},
	testStorageStruct{"get", "some_value", "qwertyuiopasdfghaa", ErrKeyTooLong.Error(), ""},
}

var testExistData = []testStorageStruct{
	testStorageStruct{"exist", "my_value_for_empty_key", "", ErrKeyNotSet.Error(), ""},
	testStorageStruct{"exist", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"exist", "", "my_key_2", "", "success"},
	testStorageStruct{"exist", "", "qwertyuiopasdfghaa", ErrKeyTooLong.Error(), ""},
	testStorageStruct{"exist", "some_value", "qwertyuiopasdfghaa", ErrKeyTooLong.Error(), ""},
}

var testDeleteData = []testStorageStruct{
	testStorageStruct{"set", "my_value_for_empty_key", "", ErrKeyNotSet.Error(), ""},
	testStorageStruct{"set", "my_value_for_long_key", "my_very_long_key_10000000", ErrKeyTooLong.Error(), ""},
	testStorageStruct{"set", "my_value_for_key_1", "my_key_1", "", "success"},
	testStorageStruct{"set", "", "my_key_2", ErrNotFound.Error(), ""},
}

var long512String = `asdhgahsdjahsdjahdahsdahsghgdsjfhasdjfhgajhsdgfj
kahsdgfjahsgdfjhasgdfjhgasjfhgashdgakshgdakjhsgdjahdsgf
jhagsdfjhgajskhdgfkajshdgfkahjsgdkjfhagsdkjhfgkashdgfkah
sgdfkjhagsdkfhgajhdgjjahgsdfjhgasdjfhgasjhdgfajhsgdfkjhasgdas
fjhgaksdhjgfashdfglSJFGALSDHJFGAJSHFGasjhdgfkahjsgdfkhagsdkj
fhgakshjdfgkahsgdfkhagskdfhgakjshdgfkjahgsdkfhagskdfhg
akjshdgfkjahfasdfjhagsdfkjhgaksjdhgfkajhsgdfkjhasgdfkjhgask
jdhfgkjashgdfkjahsgdfkjhgaskdjfhgkajshdgfkjahsgdfkjhagjsdfgasdfash
dfgkjasghdfjkhagskdjfhgakjshdgfkjhagsdkfjhgak`

var longString = `asdhgahsdjahsdjahdahsdahsghgdsjfhasdjfhgajhsdgfj
kahsdgfjahsgdfjhasgdfjhgasjfhgashdgakshgdakjhsgdjahdsgf
jhagsdfjhgajskhdgfkajshdgfkahjsgdkjfhagsdkjhfgkashdgfkah
sgdfkjhagsdkfhgajhdgjjahgsdfjhgasdjfhgasjhdgfajhsgdfkjhasgdas
fjhgaksdhjgfashdfglSJFGALSDHJFGAJSHFGasjhdgfkahjsgdfkhagsdkj
fhgakshjdfgkahsgdfkhagskdfhgakjshdgfkjahgsdkfhagskdfhg
akjshdgfkjahfasdfjhagsdfkjhgaksjdhgfkajhsgdfkjhasgdfkjhgask
jdhfgkjashgdfkjahsgdfkjhgaskdjfhgkajshdgfkjahsgdfkjhagjsdfgasdfash
dfgkjasghdfjkhagskdjfhgakjshdgfkjhagsdkfjhgakjshdfgkj_`

func TestSet(t *testing.T) {
	checkSet := func(t *testing.T, data testStorageStruct) {
		t.Helper()
		store := &Storage{
			Method: data.Method,
			Value:  data.Value,
			Key:    data.Key,
		}
		err := store.Set()
		if err != nil {
			if store.Error != data.Error {
				t.Errorf("errStr '%v'; data.Error '%v';", store.Error, data.Error)
			}
		}
		if err == nil {
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
		err := store.Get()
		if err != nil {
			if store.Error != data.Error {
				t.Errorf("Error in get method was '%v'; want '%v'", store.Error, data.Error)
			}
		}
		if err == nil {
			if store.Value != data.Value {
				t.Errorf("Value getted from db '%v'; value from test '%v", store.Value, data.Value)
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
		err := store.Exist()
		if err != nil {
			if store.Error != "" && store.Error != data.Error {
				t.Errorf("Error in exist methos was '%v'; want '%v'; key '%v'", store.Error, data.Error, data.Key)
			}
		}
		if err == nil {
			if store.Result != data.Result {
				t.Errorf("Exist result was '%v'; want '%v'", store.Result, data.Result)
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
		err := store.Delete()
		if err != nil {
			if store.Error != data.Error {
				t.Errorf("Error in delete methis was '%v'; want '%v';", store.Error, data.Error)
			}
		}
		if err == nil {
			if store.Result != "success" {
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

func TestCheckValueLength(t *testing.T) {
	t.Run("Run test for CheckValueLength Method", func(t *testing.T) {
		store := Storage{}

		store.Value = "askdhaskdjhkashdkahsdkjahsdkjhad"
		if err := CheckValueLength(store.Value, 512); err != nil {
			t.Error("CheckValueLength must return true")
		}

		store.Value = long512String
		if err := CheckValueLength(store.Value, 512); err != nil {
			t.Errorf("%v", len(store.Value))
			t.Error("CheckValueLength must return false")
		}

		store.Value = longString
		if err := CheckValueLength(store.Value, 512); err == nil {
			t.Error("CheckValueLength must return true")
		}
	})
}
