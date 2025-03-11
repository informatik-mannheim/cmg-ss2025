package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handler_http "gitty.informatik.hs-mannheim.de/steger/cmg-ws202425/services/entity/adapters/handler-http"
	repo "gitty.informatik.hs-mannheim.de/steger/cmg-ws202425/services/entity/adapters/repo-in-memory"
	"gitty.informatik.hs-mannheim.de/steger/cmg-ws202425/services/entity/core"
)

func main() {

	core := core.NewEntityService(repo.NewRepo(), nil)

	srv := &http.Server{Addr: ":8080"}

	h := handler_http.NewHandler(core)
	http.Handle("/", h)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Print("The service is shutting down...")
		srv.Shutdown(context.Background())
	}()

	log.Print("listening...")
	srv.ListenAndServe()
	log.Print("Done")
}
