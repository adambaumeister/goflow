package backends

import (
	"database/sql"
	"fmt"
	"github.com/adamb/goflow/fields"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

/*
MySQL Backend
*/
const USE_QUERY = "USE testgoflow"
const INIT_QUERY = `CREATE TABLE IF NOT EXISTS goflow_records 
(last_switched DATETIME,
 src_ip INT(4) UNSIGNED NOT NULL,
 src_port INT(2) UNSIGNED NOT NULL, 
 dst_ip INT(4) UNSIGNED NOT NULL,
 dst_port INT(2) UNSIGNED NOT NULL )`
const INSERT_QUERY = `INSERT INTO goflow_records (last_switched, src_ip, src_port, dst_ip, dst_port) VALUES ( FROM_UNIXTIME(?), ?, ?, ?, ? )`

type Mysql struct {
	Dbname string
	Dbpass string
	Dbuser string
	Server string
	db     *sql.DB
}

func (b *Mysql) Init() {
	b.Dbname = os.Getenv("SQL_DATABASE")
	b.Dbpass = os.Getenv("SQL_PASSWORD")
	b.Dbuser = os.Getenv("SQL_USERNAME")
	b.Server = os.Getenv("SQL_SERVER")

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", b.Dbuser, b.Dbpass, b.Server, b.Dbname))

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	// Try and init the database
	_, err = db.Exec(USE_QUERY)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(INIT_QUERY)
	if err != nil {
		panic(err.Error())
	}
	b.db = db
}

func (b *Mysql) Test() {
	err := b.db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func (b *Mysql) Add(values map[uint16]fields.Value) {
	db := b.db
	s, err := db.Prepare(INSERT_QUERY)
	_, err = s.Exec(
		values[fields.TIMESTAMP].ToInt(),
		values[fields.IPV4_SRC_ADDR].ToInt(),
		values[fields.L4_SRC_PORT].ToInt(),
		values[fields.IPV4_DST_ADDR].ToInt(),
		values[fields.L4_DST_PORT].ToInt(),
	)
	if err != nil {
		panic(err.Error())
	}
}
