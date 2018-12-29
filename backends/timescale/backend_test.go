package timescale

import (
	"os"
	"testing"
)

const TEST_USER = "remoteuser"

func TestBackend(t *testing.T) {
	b := Tsdb{}
	config := make(map[string]string)
	config["SQL_DB"] = "testgoflow"
	config["SQL_SERVER"] = os.Getenv("SQL_SERVER")
	config["SQL_USERNAME"] = TEST_USER

	b.Configure(config)
	db := b.connect()
	err := db.Ping()
	if err != nil {
		t.Error(err.Error())
	}
}
