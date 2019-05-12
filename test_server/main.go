package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Set() {
	for i := 0; i < 1026; i++ {
		key := RandStringBytes(16)
		values := map[string]string{"method": "set", "key": key, "value": key}
		jsonData, _ := json.Marshal(values)
		req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", strings.NewReader(string(jsonData)))
		client := &http.Client{}
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func Get() {
	for i := 0; i <= 3; i++ {
		values := map[string]string{"method": "get", "key": "EdScSmEWxIlSDgPu"}
		jsonData, _ := json.Marshal(values)
		req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", strings.NewReader(string(jsonData)))
		client := &http.Client{}
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func Del() {
	for i := 0; i <= 3; i++ {
		values := map[string]string{"method": "delete", "key": "JbSAAyRIiFjsQTNU"}
		jsonData, _ := json.Marshal(values)
		req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", strings.NewReader(string(jsonData)))
		client := &http.Client{}
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func main() {
	Del()
}
