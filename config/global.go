package config

import (
	"github.com/adamb/goflow/frontends"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GlobalConfig struct {
	Frontends map[string]FrontendConfig
}

type FrontendConfig struct {
	Type    string
	Config  map[string]string
	Backend string
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

func (gc GlobalConfig) GetFrontends() []frontends.Frontend {
	var r []frontends.Frontend
	for n, _ := range gc.Frontends {
		switch n {
		case "netflow":
			f := frontends.Netflow{}
			f.Configure(gc.Frontends["netflow"].Config)
			r = append(r, &f)
		}
	}
	return r
}
