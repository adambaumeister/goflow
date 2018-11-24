package config

import (
	"github.com/adamb/goflow/backends"
	"github.com/adamb/goflow/frontends"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GlobalConfig struct {
	Frontends map[string]FrontendConfig
	Backends  map[string]BackendConfig
}

type FrontendConfig struct {
	Type    string
	Config  map[string]string
	Backend string
}

type BackendConfig struct {
	Type   string
	Config map[string]string
}

func Read(filename string) GlobalConfig {
	gc := GlobalConfig{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &gc)
	if err != nil {
		panic(err.Error())
	}
	return gc
}

func (gc *GlobalConfig) getBackends() map[string]backends.Backend {
	bm := make(map[string]backends.Backend)
	for n, bc := range gc.Backends {
		switch bc.Type {
		case "mysql":
			b := backends.Mysql{}
			b.Configure(gc.Backends[n].Config)
			bm[n] = &b
		case "dump":
			b := backends.Dump{}
			b.Configure(gc.Backends[n].Config)
			bm[n] = &b
		}
	}
	return bm
}

func (gc GlobalConfig) GetFrontends() []frontends.Frontend {
	/*
		Returns all the configured Frontends as frontend.Frontend objects
		Maps frontends to backends in the same run
	*/
	var r []frontends.Frontend
	bm := gc.getBackends()
	for n, fields := range gc.Frontends {
		switch n {
		case "netflow":
			f := frontends.Netflow{}
			f.Configure(gc.Frontends["netflow"].Config, bm[fields.Backend])
			r = append(r, &f)
		}
	}
	return r
}