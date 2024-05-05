package main

import (
	"fmt"
	"net/http"

	"github.com/akim-malyshchyk/fantasy-backend/internal/routes"
)

func main() {
	router := routes.NewRouter()
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Started at http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err.Error())
	}
}
