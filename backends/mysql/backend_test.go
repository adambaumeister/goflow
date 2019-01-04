package mysql

import (
	"fmt"
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

func BenchmarkBackend(t *testing.B) {
	fmt.Printf("Benchmarking TIMESCALE. This will generate a lot of stuff in the database!\n")
	b := Mysql{}
	config := make(map[string]string)
	config["SQL_DB"] = "testgoflow"
	config["SQL_SERVER"] = os.Getenv("SQL_SERVER")
	config["SQL_USERNAME"] = TEST_USER

	b.Configure(config)
	b.Init()

	t.ResetTimer()
	//fmt.Printf(":::: %v  :::", t.N)
	for i := 0; i < t.N; i++ {
		testFlow := backends.GetTestFlowRand(int64(i))
		b.Add(testFlow)
	}
}
