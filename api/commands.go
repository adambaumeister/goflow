package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Commands struct {
	Paths map[string]*Path
}

type Path struct {
	Url  string
	Data []byte
}

func (p *Path) Get() {
	resp, err := http.Get("http://127.0.0.1:8880" + p.Url)
	if err != nil {
		panic(err)
	}
	jm := JsonMessage{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jm)
	fmt.Printf("%v\n", jm.Msg)
}

func (c *Commands) Parse() {
	c.Paths = make(map[string]*Path)

	// Setup the routes
	testPath := Path{
		Url: "/test",
	}
	c.Paths["test"] = &testPath
	if len(os.Args) > 1 {
		c.Paths[os.Args[1]].Get()
		os.Exit(0)
	}
}
