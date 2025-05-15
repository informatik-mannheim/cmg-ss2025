package main

/*
func main() {

	core := core.NewJobService(repo.NewRepo(), nil)

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
*/
