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

type Storager interface {
	Get() *storage.Storage
	Set() *storage.Storage
	Delete() *storage.Storage
	Exist() *storage.Storage
	RequestMethod() string
	CheckDbLenght() error
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

// NewConfig server constructor
func NewConfig(host *string, port *int) *Server {
	return &Server{
		Host: *host,
		Port: *port,
	}
}

// Routes handle routes requests
func Routes() {
	http.HandleFunc("/", HandleRootReuest)
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
func HandleRootReuest(w http.ResponseWriter, r *http.Request) {
	checkRequest(w, r)
	var store Storager = new(storage.Storage)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &store); err != nil {
		http.Error(w, "Error json format", http.StatusInternalServerError)
		return
	}

	HandleMethodRequest(w, store)
}

// HandleMethodRequest handle methods from request body
func HandleMethodRequest(w http.ResponseWriter, store Storager) {
	switch store.RequestMethod() {
	case "get":
		Get(w, store)
	case "set":
		Set(w, store)
	case "delete":
		Delete(w, store)
	case "exist":
		Exist(w, store)
	}
}

// Get function write response on get method
func Get(w http.ResponseWriter, store Storager) {
	sendingData := store.Get()
	data, err := json.Marshal(sendingData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Set function set data to Db and write response
func Set(w http.ResponseWriter, store Storager) {
	err := store.CheckDbLenght()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	sendingData := store.Set()
	data, err := json.Marshal(sendingData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Delete function delete data from Db
func Delete(w http.ResponseWriter, store Storager) {
	sendingData := store.Delete()
	data, err := json.Marshal(sendingData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Exist function check if key present in db
func Exist(w http.ResponseWriter, store Storager) {
	sendingData := store.Exist()
	data, err := json.Marshal(sendingData)
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
