package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Commands struct {
	Paths      map[string]Path
	HelpString string
}

type Path interface {
	Get()
	Help() string
}

func (c *Commands) Get() {
	s := "Goflow command line client help\n------------------\n\n"
	for cmd, p := range c.Paths {
		s = s + fmt.Sprintf("%v : %v\n", cmd, p.Help())
	}
	fmt.Printf(s)
}
func (c *Commands) Help() string {
	return "Display this help message."
}

func parseArgs(a []string) []string {
	i := 2
	var r []string
	for _, arg := range a {
		if i+1 > len(os.Args) {
			fmt.Printf("Error: Missing required argument: %v\n", arg)
			os.Exit(1)
		}
		r = append(r, os.Args[i])
		i++
	}
	return r
}

/*
Post constructs
Grafana
*/
type GrafanaPath struct {
	Url        string
	Data       []byte
	HelpString string
	Body       []byte

	Args []string
}

func (p *GrafanaPath) Get() {
	jg := JsonGrafana{}
	args := parseArgs(p.Args)
	jg.Server = args[0]
	jg.ApiKey = args[1]
	j, _ := json.Marshal(jg)
	resp, err := http.Post("http://127.0.0.1:8880"+p.Url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		panic(err)
	}

	jm := JsonMessage{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jm)
	fmt.Printf("%v\n", jm.Msg)
}
func (p *GrafanaPath) Help() string {
	return fmt.Sprintf("Nothing yet")
}

/*
HTTP Path do a call to the API to present their results

No POST data, just a generic HTTP get.
*/
type GenericPath struct {
	Url        string
	Data       []byte
	HelpString string
	Body       []byte

	Args []string
}

func (p *GenericPath) Get() {
	resp, err := http.Get("http://127.0.0.1:8880" + p.Url)
	if err != nil {
		panic(err)
	}
	jm := JsonMessage{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jm)
	fmt.Printf("%v\n", jm.Msg)
}

func (p *GenericPath) Help() string {
	return p.HelpString
}

func (c *Commands) Parse() {
	c.Paths = make(map[string]Path)

	// Setup the routes
	testPath := GenericPath{
		Url:        "/status",
		HelpString: "Displays the status of configured Backends",
	}
	grafanaPath := GrafanaPath{
		Url:        "/grafana",
		HelpString: "Configures a Grafana instance with Goflow compatible dashboards",
		Args:       []string{"Grafana server", "Grafana API Key"},
	}

	c.Paths["status"] = &testPath
	c.Paths["configure-grafana"] = &grafanaPath
	c.Paths["help"] = c
	if len(os.Args) > 1 {
		c.Paths[os.Args[1]].Get()
		os.Exit(0)
	}
}
