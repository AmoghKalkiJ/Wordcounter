package main

import (
	"Wordcounter/pkg/counter"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/counter", counter.Portal)
	r.HandleFunc("/api/wordcounter", counter.Wordcounter)
	fmt.Println("Started server on  localhost:8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		panic(err)
	}

}
