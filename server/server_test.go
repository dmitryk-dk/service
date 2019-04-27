package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dmitryk-dk/service/storage"
)

var testNewConfig = map[string]Server{
	"first": Server{
		Port: 20,
		Host: "",
	},
	"second": Server{
		Port: 0,
		Host: "127.0.0.1",
	},
	"third": Server{
		Port: 2000,
		Host: "127.0.0.1",
	},
	"forth": Server{
		Port: -5,
		Host: "some_value",
	},
}

type Request struct {
	Method string
	Key    string
	Value  string
}
type Response struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
	Key    string `json:"key"`
	Error  string `json:"error,omitempty"`
	Result string `json:"result,omitempty"`
}

var testSetRequest = []storage.Storage{
	storage.Storage{Method: "set", Key: "my_new_key_toooolong_111111111", Value: "value"},
	storage.Storage{Method: "set", Key: "", Value: "my_value"},
	storage.Storage{Method: "set", Key: "my_new_key_1", Value: "my_value_1"},
	storage.Storage{Method: "set", Key: "my_new_key_2", Value: ""},
}
var testSetResponse = []Response{
	Response{Method: "set", Key: "my_new_key_toooolong_111111111", Error: "key too long"},
	Response{Method: "set", Key: "", Value: "", Error: "key not set"},
	Response{Method: "set", Key: "my_new_key_1", Result: "success"},
	Response{Method: "set", Key: "my_new_key_2", Value: "", Result: "success"},
}

var testGetRequest = []storage.Storage{
	storage.Storage{Method: "get", Key: "my_new_key_toooolong_111111111"},
	storage.Storage{Method: "get", Key: ""},
	storage.Storage{Method: "get", Key: "my_new_key_1"},
	storage.Storage{Method: "get", Key: "my_new_key_3"},
}
var testGetResponse = []storage.Storage{
	storage.Storage{Method: "get", Key: "my_new_key_toooolong_111111111", Error: "key too long"},
	storage.Storage{Method: "get", Key: "", Error: "key not set"},
	storage.Storage{Method: "get", Key: "my_new_key_1", Value: "my_value_1"},
	storage.Storage{Method: "get", Key: "my_new_key_3", Error: "value not found"},
}

var testDeleteRequest = []storage.Storage{
	storage.Storage{Method: "delete", Key: "my_new_key_toooolong_111111111"},
	storage.Storage{Method: "delete", Key: ""},
	storage.Storage{Method: "delete", Key: "my_new_key_1"},
}

var testDeleteResponse = []storage.Storage{
	storage.Storage{Method: "delete", Key: "my_new_key_toooolong_111111111", Error: "key too long"},
	storage.Storage{Method: "delete", Key: "", Error: "key not set"},
	storage.Storage{Method: "delete", Key: "my_new_key_1", Result: "success"},
}

var testExistRequest = []storage.Storage{
	storage.Storage{Method: "exist", Key: "my_new_key_toooolong_111111111"},
	storage.Storage{Method: "exist", Key: ""},
	storage.Storage{Method: "exist", Key: "my_new_key_1"},
	storage.Storage{Method: "exist", Key: "my_new_key_3"},
}

var testExistResponse = []storage.Storage{
	storage.Storage{Method: "exist", Key: "my_new_key_toooolong_111111111", Error: "key too long"},
	storage.Storage{Method: "exist", Key: "", Error: "key not set"},
	storage.Storage{Method: "exist", Key: "my_new_key_1", Result: "success"},
	storage.Storage{Method: "exist", Key: "my_new_key_3", Error: "value not found"},
}

var testHandleRequest = []storage.Storage{
	storage.Storage{Method: "set", Key: "my_new_key_toooolong_111111111", Value: "my_new_value_toooolong"},
	storage.Storage{Method: "set", Key: "my_new_key_1", Value: "my_new_value_1"},
	storage.Storage{Method: "get", Key: "my_new_key_1"},
	storage.Storage{Method: "get", Key: "my_new_key_3"},
	storage.Storage{Method: "exist", Key: "my_new_key_1"},
	storage.Storage{Method: "exist", Key: "my_new_key_3"},
	storage.Storage{Method: "delete", Key: ""},
	storage.Storage{Method: "delete", Key: "my_new_key_1"},
}

var testHandleResponse = []storage.Storage{
	storage.Storage{Method: "set", Key: "my_new_key_toooolong_111111111", Error: "key too long"},
	storage.Storage{Method: "set", Key: "my_new_key_1", Result: "success"},
	storage.Storage{Method: "get", Key: "my_new_key_1", Value: "my_new_value_1"},
	storage.Storage{Method: "get", Key: "my_new_key_3", Error: "value not found"},
	storage.Storage{Method: "exist", Key: "my_new_key_1", Result: "success"},
	storage.Storage{Method: "exist", Key: "my_new_key_3", Error: "value not found"},
	storage.Storage{Method: "delete", Key: "", Error: "key not set"},
	storage.Storage{Method: "delete", Key: "my_new_key_1", Result: "success"},
}

func buildUrl(scheme, host string, port int) string {
	addr := fmt.Sprintf("%s:%d", host, port)
	return scheme + "://" + addr
}

func TestNewConfig(t *testing.T) {
	checkNewConfig := func(t *testing.T, data Server) {
		config := NewConfig(&data.Host, &data.Port)
		if !reflect.DeepEqual(config, &data) {
			t.Error("somthing wrong")
		}
	}
	t.Run("Run test for NewConfig Method", func(t *testing.T) {
		for _, v := range testNewConfig {
			checkNewConfig(t, v)
		}
	})
}

func TestSetMethod(t *testing.T) {
	checkSet := func(t *testing.T, data storage.Storage, idx int) {
		response, _ := json.Marshal(testSetResponse[idx])

		_, err := http.NewRequest("GET", buildUrl("http", "127.0.0.1", 8000), nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()

		Set(res, &data)
		if res.Code != http.StatusOK {
			t.Errorf("Response code was %v; want 200", res.Code)
		}

		if bytes.Compare(response, res.Body.Bytes()) != 0 {
			t.Errorf("Response body was '%v'; want '%v';", res.Body, string(response))
		}
	}
	t.Run("Run test for Set Method", func(t *testing.T) {
		for idx, v := range testSetRequest {
			checkSet(t, v, idx)
		}
	})
}

func TestGetMethod(t *testing.T) {
	checkGet := func(t *testing.T, data storage.Storage, idx int) {
		response, _ := json.Marshal(testGetResponse[idx])

		_, err := http.NewRequest("GET", buildUrl("http", "127.0.0.1", 8000), nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()

		Get(res, &data)
		if res.Code != http.StatusOK {
			t.Errorf("Response code was %v; want 200", res.Code)
		}

		if bytes.Compare(response, res.Body.Bytes()) != 0 {
			t.Errorf("Response body was '%v'; want '%v';", res.Body, string(response))
		}
	}
	t.Run("Run test for Get Method", func(t *testing.T) {
		for idx, v := range testGetRequest {
			checkGet(t, v, idx)
		}
	})
}
func TestExistMethod(t *testing.T) {
	checkExist := func(t *testing.T, data storage.Storage, idx int) {
		response, _ := json.Marshal(testExistResponse[idx])

		_, err := http.NewRequest("GET", buildUrl("http", "127.0.0.1", 8000), nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()

		Delete(res, &data)
		if res.Code != http.StatusOK {
			t.Errorf("Response code was %v; want 200", res.Code)
		}

		if bytes.Compare(response, res.Body.Bytes()) != 0 {
			t.Errorf("Response body was '%v'; want '%v';", res.Body, string(response))
		}
	}
	t.Run("Run test for Exist Method", func(t *testing.T) {
		for idx, v := range testExistRequest {
			checkExist(t, v, idx)
		}
	})
}

func TestDeleteMethod(t *testing.T) {
	checkDelete := func(t *testing.T, data storage.Storage, idx int) {
		response, _ := json.Marshal(testDeleteResponse[idx])

		_, err := http.NewRequest("GET", buildUrl("http", "127.0.0.1", 8000), nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()

		Delete(res, &data)
		if res.Code != http.StatusOK {
			t.Errorf("Response code was %v; want 200", res.Code)
		}

		if bytes.Compare(response, res.Body.Bytes()) != 0 {
			t.Errorf("Response body was '%v'; want '%v';", res.Body, string(response))
		}
	}
	t.Run("Run test for Delete Method", func(t *testing.T) {
		for idx, v := range testDeleteRequest {
			checkDelete(t, v, idx)
		}
	})
}

func TestHandleMethodRequest(t *testing.T) {
	checkHandle := func(t *testing.T, data storage.Storage, idx int) {
		response, _ := json.Marshal(testHandleResponse[idx])

		_, err := http.NewRequest("GET", buildUrl("http", "127.0.0.1", 8000), nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		HandleMethodRequest(res, &data)
		if res.Code != http.StatusOK {
			t.Errorf("Response code was %v; want 200", res.Code)
		}

		if bytes.Compare(response, res.Body.Bytes()) != 0 {
			t.Errorf("Response body was '%v'; want '%v';", res.Body, string(response))
		}
	}
	t.Run("Run test for Handle Method", func(t *testing.T) {
		for idx, v := range testHandleRequest {
			checkHandle(t, v, idx)
		}
	})
}
