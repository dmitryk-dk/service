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
)

type Server struct {
	Host string
	Port int
}

type Storager interface {
	Get() error
	Set() error
	Delete() error
	Exist() error
	RequestMethod() string
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
func Routes(store Storager) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleRootReuest(w, r, store)
	})
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
func HandleRootReuest(w http.ResponseWriter, r *http.Request, store Storager) {
	checkRequest(w, r)
	//var store Storager = new(storage.Storage)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, store); err != nil {
		http.Error(w, "Error json format", http.StatusInternalServerError)
		return
	}

	HandleMethodRequest(w, store)
}

// HandleMethodRequest handle methods from request body
func HandleMethodRequest(w http.ResponseWriter, store Storager) {
	switch store.RequestMethod() {
	case "get":
		err := store.Get()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		SendData(w, store)
	case "set":
		err := store.Set()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		SendData(w, store)
	case "delete":
		err := store.Delete()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		SendData(w, store)
	case "exist":
		err := store.Exist()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		SendData(w, store)
	}
}

// Get function write response on get method
func SendData(w http.ResponseWriter, store Storager) {
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
