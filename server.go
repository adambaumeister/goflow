package main

import (
	"github.com/adamb/goflow/config"
)

func main() {
	gc := config.Read("config.yml")

	fe := gc.GetFrontends()
	for _, f := range fe {
		f.Start()
	}
}
