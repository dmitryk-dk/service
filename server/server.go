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
	if store.Exist() {
		switch store.Method {
		case "GET":
			store.Get()
		case "SET":
			store.Set()
		case "DELETE":
			store.Delete()
		default:
			store.Exist()
		}
	}
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
