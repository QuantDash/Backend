package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mannan-Ali/Backend/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Print("Check Mannan")
	godotenv.Load() //to load the env file we use this function

	//to use env go has the follwing functoin
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}
	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Mount the v1 API router
	mainRouter.Mount("/v1", api.NewV1Router())

	mainRouter.Get("/healthcheck", handlerReadiness)
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: mainRouter,
	}
	fmt.Println("Starting server on PORT:", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
