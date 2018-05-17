package main

import (
	"fmt"
	"net/http"
)

func nsaIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "INDEX - Hello, %q", r.URL.Path[1:])
}

func nsaOtherHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OTHER - Hello, %q", r.URL.Path[1:])
}
