package handlers

import (
	"fmt"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running!")
}
