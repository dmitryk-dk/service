package main

import (
	"flag"

	"github.com/dmitryk-dk/service/server"
	"github.com/dmitryk-dk/service/storage"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "host")
	var port = flag.Int("port", 8000, "port")
	flag.Parse()
	var store server.Storager = new(storage.Storage)

	server.Routes(store)
	server.PrepareShutdown()
	srv := server.NewConfig(host, port)
	srv.ListenAndServ()
}
