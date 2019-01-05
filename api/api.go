package api

import (
	"encoding/json"
	"fmt"
	"github.com/adambaumeister/goflow/config"
	"github.com/adambaumeister/goflow/utils/grafana"
	"io/ioutil"
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

type JsonGrafana struct {
	Server string
	ApiKey string
}

func Start(gc *config.GlobalConfig) {
	a := API{}
	a.config = gc

	http.HandleFunc("/", a.getHandler)
	http.HandleFunc("/status", a.Test)
	http.HandleFunc("/grafana", a.Grafana)
	log.Fatal(http.ListenAndServe(a.config.Api, nil))

}

func (a *API) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API works!")
}

func (a *API) Grafana(w http.ResponseWriter, r *http.Request) {
	jg := JsonGrafana{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &jg)
	g := grafana.Grafana{}

	s := fmt.Sprintf("You got this far. %v %v", jg.Server, jg.ApiKey)
	g.Server = jg.Server
	g.Key = jg.ApiKey
	jm := JsonMessage{
		Msg: s,
	}
	j, err := json.Marshal(jm)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(j)
}

func (a *API) Test(w http.ResponseWriter, r *http.Request) {
	var s string
	b := a.config.GetBackends()
	for _, be := range b {
		s = s + be.Status() + "\n"
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
