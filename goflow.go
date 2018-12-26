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

	fmt.Printf("Starting threads...")
	fe := gc.GetFrontends()
	for _, f := range fe {
		go f.Start()
	}
	api.Start(&gc)
}
