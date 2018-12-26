package api

import (
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
		s = s + fmt.Sprintf("%v			: %v\n", cmd, p.Help())
	}
	fmt.Printf(s)
}
func (c *Commands) Help() string {
	return "Display this help message."
}

/*
HTTP Path do a call to the API to present their results
*/
type HttpPath struct {
	Url        string
	Data       []byte
	HelpString string
}

func (p *HttpPath) Get() {
	resp, err := http.Get("http://127.0.0.1:8880" + p.Url)
	if err != nil {
		panic(err)
	}
	jm := JsonMessage{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jm)
	fmt.Printf("%v\n", jm.Msg)
}

func (p *HttpPath) Help() string {
	return p.HelpString
}

func (c *Commands) Parse() {
	c.Paths = make(map[string]Path)

	// Setup the routes
	testPath := HttpPath{
		Url:        "/status",
		HelpString: "Displays the status of configured Backends",
	}

	c.Paths["status"] = &testPath
	c.Paths["help"] = c
	if len(os.Args) > 1 {
		c.Paths[os.Args[1]].Get()
		os.Exit(0)
	}
}
