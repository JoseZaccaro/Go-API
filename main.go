package main

import (
	"api/usuarios/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Usuarios API!")
	})

	mux.HandleFunc("/api/register", routes.RegistrarUsuario).Methods("POST")
	mux.HandleFunc("/api/login", routes.LoguearUsuario).Methods("POST")

	//* variables de entorno
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}

	//* servidor
	server := &http.Server{
		Addr:         "127.0.0.1:" + os.Getenv("PORT"),
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on port " + os.Getenv("PORT"))

	log.Fatal(server.ListenAndServe())
}
