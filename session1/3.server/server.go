package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()

	server := http.Server{
		Addr:         ":8001",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "%s received for hello endpoint\n", request.Method)
	})

	mux.Handle("/hello2", handler{})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type handler struct {
}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connection received")
	_, _ = fmt.Fprintf(w, "%s received for hello2 endpoint\n", r.Method)
}