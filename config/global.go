package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GlobalConfig struct {
	Frontends []FrontendConfig
}

type FrontendConfig struct {
	Type   string
	Config map[string]string
}

func Read(filename string) {
	gc := GlobalConfig{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &gc)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v", gc.Frontends[0].Config["bindaddr"])
}
