package api

import (
	"encoding/json"
	"fmt"
	"github.com/adambaumeister/goflow/config"
	"log"
	"net/http"
)

type API struct {
	c      chan string
	config *config.GlobalConfig
}

type JsonMessage struct {
	Msg string
}

func Start(gc *config.GlobalConfig) {
	a := API{}
	a.config = gc

	http.HandleFunc("/", a.getHandler)
	http.HandleFunc("/status", a.Test)
	log.Fatal(http.ListenAndServe("127.0.0.1:8880", nil))

}

func (a *API) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API works!")
}

func (a *API) Test(w http.ResponseWriter, r *http.Request) {
	var s string
	b := a.config.GetBackends()
	for _, be := range b {
		s = s + be.Test() + "\n"
	}
	jm := JsonMessage{
		Msg: s,
	}
	j, err := json.Marshal(jm)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(j)
}
