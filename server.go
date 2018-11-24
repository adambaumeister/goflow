package main

import (
	"github.com/adamb/goflow/backends"
	"github.com/adamb/goflow/config"
)

func main() {
	gc := config.Read("config.yml")
	b := backends.Mysql{}
	b.Init()
	fe := gc.GetFrontends()
	for _, f := range fe {
		f.Start(&b)
	}

	/*
		nf := frontends.Netflow{}

		b := backends.Mysql{}
		b.Init()
		nf.Start(&b)
	*/
}
