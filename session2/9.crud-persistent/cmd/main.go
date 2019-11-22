package main

import (
	"github.com/gorilla/mux"
	"github.com/wgarunap/devfest/session1/9.crud-persistent/pkg/handlers"
	// "github.com/wgarunap/devfest/session1/9.crud-persistent/pkg/models"

	"net/http"
	"time"
)

func main() {

	// handlers.PersonMap = make(map[int]models.Person, 0)

	router := mux.NewRouter()

	server := http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	router.Handle("/person", handlers.HandlerPost{}).
		Methods(http.MethodPost).
		Name("create-person-info").
		Headers("content-type", "application/json")

	router.Handle("/person/{id}", handlers.HandlerGet{}).
		Methods(http.MethodGet).
		Name("get-person-info").
		Headers("content-type", "application/json")

	router.Handle("/person", handlers.HandlerGetAll{}).
		Methods(http.MethodGet).
		Name("get-person-info").
		Headers("content-type", "application/json")

	router.Handle("/person/{id}", handlers.HandlerPut{}).
		Methods(http.MethodPut).
		Name("update-person-info").
		Headers("content-type", "application/json")

	router.Handle("/person/{id}", handlers.HandlerDelete{}).
		Methods(http.MethodDelete).
		Name("delete-person-info").
		Headers("content-type", "application/json")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
