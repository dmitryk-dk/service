package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitryk-dk/service/storage"
)

type Server struct {
	Host string
	Port int
}

// NewConfig server constructor
func NewConfig(host *string, port *int) *Server {
	return &Server{
		Host: *host,
		Port: *port,
	}
}

// Routes handle routes requests
func Routes() {
	http.HandleFunc("/", Values)
}

func checkRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Values handle request to root path
func Values(w http.ResponseWriter, r *http.Request) {
	checkRequest(w, r)
	var store storage.Storage
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &store)
	if err != nil {
		http.Error(w, "Error json format", http.StatusInternalServerError)
		return
	}
	err = CheckValueLength(w, &store)
	if err != nil {
		store.Error = "Value too long"
		store.Value = ""
		data, err := json.Marshal(store)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	switch store.Method {
	case "get":
		Get(w, &store)
	case "set":
		Set(w, &store)
	case "delete":
		Delete(w, &store)
	case "exist":
		Get(w, &store)
	}
}

// CheckValueLength define value lenght and if it more than 512 byte return error
func CheckValueLength(w http.ResponseWriter, store *storage.Storage) error {
	if len(store.Value) > 512 {
		return fmt.Errorf("value length too much")
	}
	return nil
}

// Get function write response on get method
func Get(w http.ResponseWriter, store *storage.Storage) {
	val, err := store.Get()
	if err != nil || val == "" {
		store.Error = "Value not found"
	} else {
		if store.Method != "exist" {
			store.Value = val
		}
		store.Result = "success"
	}
	data, err := json.Marshal(store)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Set function set data to Db and write response
func Set(w http.ResponseWriter, store *storage.Storage) {
	err := store.Set()
	if err != nil {
		store.Error = "key to long"
	} else {
		store.Result = "success"
	}
	data, err := json.Marshal(store)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Delete function delete data from Db
func Delete(w http.ResponseWriter, store *storage.Storage) {
	err := store.Delete()
	if err != nil {
		store.Error = "Value not found"
	} else {
		store.Result = "success"
	}
	data, err := json.Marshal(store)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// PrepareShutdown prepare server to shutdown
func PrepareShutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println()
		log.Printf("Got signal %d", <-sig)
		os.Exit(0)
	}()
}

// ListenAndServ function run http server
func (s *Server) ListenAndServ() {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	fmt.Printf("Server start in addr: %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
