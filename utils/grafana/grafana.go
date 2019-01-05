package grafana

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const URL_DATASOURCE = "/api/datasources"
const DS_JSON = `{
  "name": "goflow",
  "type": "postgres",
  "access": "proxy",
  "url":"%v",
  "secureJsonData": {
    "password": "%v"
  },
  "user":"remoteuser",
  "database":"%v",
  "jsonData":{"postgresVersion":1000,"sslmode":"disable","timescaledb":true}
}`

type Grafana struct {
	Server string
	Key    string
	Log    []string
}

func (g *Grafana) AddDataSource(n string, config map[string]string) string {
	g.Log = append(g.Log, fmt.Sprintf("Adding DS %v\n", n))

	jsonString := []byte(fmt.Sprintf(DS_JSON, config["SQL_SERVER"], os.Getenv("SQL_PASSWORD"), config["SQL_DB"]))
	//fmt.Printf(jsonString)
	//j, err := json.Marshal(jsonString)
	//if err != err {
	//	panic(err)
	//}
	req, err := http.NewRequest("POST", "http://"+g.Server+URL_DATASOURCE, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", g.Key))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
