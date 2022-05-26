package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	serverIP   = flag.String("s", "", "Server IP")
	serverPort = flag.String("p", "8081", "Server Port")
)

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./words/")))
	http.ListenAndServe(*serverIP+":"+*serverPort, r)
}
