package main

import (
	"fmt"
	"github.com/adambaumeister/goflow/api"
	"github.com/adambaumeister/goflow/config"
	"os"
)

func main() {

	c := api.Commands{}
	c.Parse()

	var gc config.GlobalConfig
	configLoc := os.Getenv("GOFLOW_CONFIG")
	if len(configLoc) != 0 {
		gc = config.Read(configLoc)
	} else {
		gc = config.Read("config.yml")
	}

	fmt.Printf("Starting Frontends...")
	fe := gc.GetFrontends()
	for _, f := range fe {
		go f.Start()
	}
	fmt.Printf("[ OK ]\n")
	fmt.Printf("Starting utilities...")
	utilities := gc.GetUtilities()
	for _, u := range utilities {
		go u.Run()
	}
	fmt.Printf("[ OK ]\n")
	api.Start(&gc)
}
