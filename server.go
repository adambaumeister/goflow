package main

import (
	"fmt"
	"github.com/adamb/goflow/config"
	"time"
)

func main() {
	gc := config.Read("config.yml")

	fmt.Printf("Starting threads...")
	fe := gc.GetFrontends()
	for _, f := range fe {
		go f.Start()
	}
	time.Sleep(10 * time.Second)
}
