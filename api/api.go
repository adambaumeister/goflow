package api

import (
	"fmt"
	"github.com/adambaumeister/goflow/config"
	"log"
	"net/http"
)

type API struct {
	c      chan string
	config *config.GlobalConfig
}

func Start(gc *config.GlobalConfig) {
	a := API{}
	a.config = gc

	http.HandleFunc("/", a.getHandler)
	http.HandleFunc("/test", a.Test)
	log.Fatal(http.ListenAndServe(":8880", nil))

}

func (a *API) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API works!")
}

func (a *API) Test(w http.ResponseWriter, r *http.Request) {
	var s string
	b := a.config.GetBackends()
	fmt.Fprintf(w, "Status:")
	for _, be := range b {
		s = s + be.Test() + "\n"
	}

	fmt.Fprintf(w, s)
}
