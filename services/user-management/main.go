package main

import (
	"log"
	"net/http"
	"os"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
)

func main() {

	http.HandleFunc("/auth/register", handler.RegisterHandler)
	http.HandleFunc("/auth/login", handler.LoginHandler)

	port := getPort()
	log.Printf("User Management Service running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8081"
	}
	return port
}
