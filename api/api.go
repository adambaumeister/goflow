package api

import (
	"fmt"
	"log"
	"net/http"
)

type API struct {
	c chan string
}

func Start() {
	a := API{}
	http.HandleFunc("/", a.getHandler)
	http.HandleFunc("/test", a.Test)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (a *API) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API works!")
}

func (a *API) Test(w http.ResponseWriter, r *http.Request) {
	TestNFFrontend()
}
