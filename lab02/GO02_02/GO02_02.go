package main

import (
	go02_02lib "GO02_02/GO02_02lib"
	"fmt"
	"net/http"
)

var A01 int = 3

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "A01 = %d \nA02 = %t \nA03 = %s", A01, A02, go02_02lib.A03)
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed", r.Method)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed", r.Method)
		}
	})

	http.ListenAndServe(":4000", nil)
}
