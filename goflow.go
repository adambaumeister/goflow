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
