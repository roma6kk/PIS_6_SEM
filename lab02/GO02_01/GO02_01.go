package main

import (
	go02_01lib "GO02_01/GO02_01lib"
	"fmt"
	"net/http"
)

const C01 = 3.14

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "C01: %f,\nC02: %e,\nC03: %f", C01, C02, go02_01lib.C03)

		case http.MethodPost, http.MethodPut, http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed", r.Method)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed", r.Method)
		}
	})
	http.ListenAndServe(":8081", nil)
}
