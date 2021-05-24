package main

import (
	"fmt"
	"go-mysql/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running"))
	})

	router.HandleFunc("/users", server.StoreUser).Methods(http.MethodPost)
	router.HandleFunc("/users", server.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", server.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", server.DeleteUser).Methods(http.MethodDelete)

	fmt.Println("Server listen on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
