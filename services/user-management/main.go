package main

import (
	"log"
	"net/http"
	"os"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables (optional)
	_ = godotenv.Load()

	// Register API endpoints
	http.HandleFunc("/auth/register", handler.RegisterHandler)
	http.HandleFunc("/auth/login", handler.LoginHandler)

	// Start HTTP server
	port := getPort()
	log.Printf("User Management Service running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getPort returns PORT from environment or falls back to 8081
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8081"
	}
	return port
}
