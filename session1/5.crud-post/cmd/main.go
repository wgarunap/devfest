package cmd

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {

	//person = make(map[int]Person, 0)

	router := mux.NewRouter()

	server := http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	router.Handle("/person", handler{}).
		Methods(http.MethodPost).
		Name("add-post").
		Headers("content-type", "application/json")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

