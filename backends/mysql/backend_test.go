package mysql

import (
	"github.com/adambaumeister/goflow/backends"
	"os"
	"testing"
)

const TEST_USER = "remoteuser"

/*
Test this backend using the dummy set of data

Requires a running instance and the following env variables exported:
	- SQL_SERVER
	- SQL_PASSWORD
*/
func TestBackend(t *testing.T) {
	b := Mysql{}
	config := make(map[string]string)
	config["SQL_DB"] = "testgoflow"
	config["SQL_SERVER"] = os.Getenv("SQL_SERVER")
	config["SQL_USERNAME"] = TEST_USER

	b.Configure(config)
	b.Init()

	testFlow := backends.GetTestFlow()
	b.Add(testFlow)
}
