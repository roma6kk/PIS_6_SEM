package main

import (
	go02_01lib "GO02_01/go02_01lib"
	"fmt"
	"net/http"
)

const C01 = 3.14

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "C01: %f,\nC02: %e,\nC03: %f", C01, C02, go02_01lib.C03)
	})
	http.ListenAndServe(":8081", nil)
}
