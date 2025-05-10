package main

import (
	"net/http"
	"os"

	auth0adapter "github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/auth"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/handler-http"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

func main() {
	useLive := os.Getenv("USE_LIVE") == "true"
	n := notifier.New()
	auth := auth0adapter.New(useLive, n)

	notifierFn := func() ports.Notifier {
		return n
	}

	h := handler.New(auth, useLive, handler.IsAdmin, notifierFn)

	http.HandleFunc("/auth/register", h.RegisterHandler)
	http.HandleFunc("/auth/login", h.LoginHandler)

	port := ":8080"
	n.Event("Listening on " + port)
	http.ListenAndServe(port, nil)
}
