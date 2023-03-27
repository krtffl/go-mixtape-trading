package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	log.Printf("received %s request at %s from %s", r.Method, r.URL, r.RemoteAddr)
	fmt.Fprintf(w, "Server is up and running!")
}
