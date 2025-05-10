package main

import (
	"net/http"
	"os"

	auth0adapter "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
)

func main() {
	useLive := os.Getenv("USE_LIVE") == "true"
	auth := auth0adapter.New(useLive)
	service := core.NewService()

	h := handler.New(service, auth, useLive, handler.IsAdmin, notifier.New)

	http.HandleFunc("/auth/register", h.RegisterHandler)
	http.HandleFunc("/auth/login", h.LoginHandler)

	port := ":8080"
	println("Listening on", port)
	http.ListenAndServe(port, nil)
}
