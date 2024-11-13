package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akim-malyshchyk/fantasy-backend/internal/routes"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName),
	)
	if err != nil {
		log.Fatalf("Failed to prepare database driver: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := routes.NewRouter(db)
	port := 3000
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Started at http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err.Error())
	}
}
