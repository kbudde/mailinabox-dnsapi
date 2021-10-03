package main

import (
	"log"
	"net/http"
)

var (
	Revision  = ""
	Branch    = ""
	BuildDate = ""
)

func main() {
	log.Printf("Starting application. branch=%s revision=%s builddate=%s",
		Branch, Revision, BuildDate)
	a, err := initApp()
	if err != nil {
		log.Fatalf("invalid configuration: %s\n", err.Error())
	}

	http.HandleFunc("/cleanup", a.RequestHandler())
	http.HandleFunc("/present", a.RequestHandler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
