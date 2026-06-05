package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mannan-Ali/Backend/api"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}
	mainRouter := chi.NewRouter()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	//establishing connection with mysql
	cons, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	//this apiconfig now can be passed into different function and they will have access to our database
	//apiconfig allows us to use the functions we created using sql queires
	db := internaldb.New(cons)
	// We use the exported APIConfig from the api package!
	apiCfg := &api.APIConfig{
		DB: db,
	}
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Mount the v1 API router
	mainRouter.Mount("/v1", api.NewV1Router(apiCfg))

	mainRouter.Get("/healthcheck", handlerReadiness)
	mainRouter.Get("/err", handlerError)
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: mainRouter,
	}
	fmt.Println("Starting server on PORT:", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
