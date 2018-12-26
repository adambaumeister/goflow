package main

import (
	"fmt"
	"github.com/adambaumeister/goflow/api"
	"github.com/adambaumeister/goflow/config"
)

func main() {

	c := api.Commands{}
	c.Parse()

	gc := config.Read("config.yml")

	fmt.Printf("Starting Frontends...\n")
	fe := gc.GetFrontends()
	for _, f := range fe {
		go f.Start()
	}
	fmt.Printf("Starting utilities...\n")
	utilities := gc.GetUtilities()
	for _, u := range utilities {
		go u.Run()
	}
	api.Start(&gc)
}
