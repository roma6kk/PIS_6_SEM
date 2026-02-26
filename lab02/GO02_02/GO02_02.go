package main

import (
	"fmt"
	"net/http"

	go02_02lib "GO02_02/GO02_02lib"
)

var A01 int = 3

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "A01 = %d \nA02 = %t \nA03 = %s", A01, A02, go02_02lib.A03)
	})

	http.ListenAndServe(":4000", nil)
}
