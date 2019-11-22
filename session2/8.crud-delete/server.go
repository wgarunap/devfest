package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {

	router := mux.NewRouter()

	server := http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	router.Handle("/hello", handler{}).
		Methods(http.MethodGet).
		Name("get-hello").
		Headers("content-type", "application/json")

	router.Handle("/hello", handler{}).
		Methods(http.MethodPost).
		Name("post-hello")

	//router.Handle("/hello2", handler{})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type handler struct {
}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connection received")
	_, _ = fmt.Fprintf(w, "%s received for hello endpoint with handler \n", r.Method)
}
