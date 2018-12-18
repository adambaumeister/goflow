package backends

import (
	"database/sql"
	"fmt"
	"github.com/adambaumeister/goflow/fields"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

/*
MySQL Backend
*/
const USE_QUERY = "USE %v"
const CHECK_QUERY = "SHOW COLUMNS IN goflow_records;"
const ALTER_QUERY = "ALTER TABLE goflow_records ADD COLUMN %v %v %v;"
const INIT_TEMPLATE = `CREATE TABLE IF NOT EXISTS goflow_records (%v);`
const INSERT_TEMPLATE = `INSERT INTO goflow_records (%v) VALUES (%v);`
const DROP_QUERY = "DROP TABLE goflow_records"

type Mysql struct {
	Dbname string
	Dbpass string
	Dbuser string
	Server string
	db     *sql.DB

	schema *Schema
}

type Column struct {
	Name    string
	Options string
	Type    string
	Wrap    string
}

type Schema struct {
	columns []*Column
}

func (s *Schema) AddColumn(n string, t string, o string) *Column {
	c := Column{
		Name:    n,
		Options: o,
		Type:    t,
	}
	s.columns = append(s.columns, &c)
	return &c
}
func (s *Schema) GetColumnStrings(t string) string {
	var qs []string
	for _, col := range s.columns {
		qs = append(qs, col.init())
	}
	return fmt.Sprintf(t, strings.Join(qs, ", "))
}

func (s *Schema) GetColumn(c string) *Column {
	for _, col := range s.columns {
		if col.Name == c {
			return col
		}
	}
	return nil
}

func (s *Schema) InsertQueryValues() string {
	var qs []string
	for _, col := range s.columns {
		qs = append(qs, col.insert())
	}
	return strings.Join(qs, ", ")
}

func (s *Schema) InsertQuery(t string) string {
	return fmt.Sprintf(t, s.InsertQueryFields(), s.InsertQueryValues())
}

func (s *Schema) InsertQueryFields() string {
	var qs []string
	for _, col := range s.columns {
		qs = append(qs, col.Name)
	}
	return strings.Join(qs, ", ")
}

func (c *Column) init() string {
	return fmt.Sprintf("%v %v %v", c.Name, c.Type, c.Options)
}
func (c *Column) insert() string {
	if len(c.Wrap) > 0 {
		return fmt.Sprintf(c.Wrap, "?")
	}
	return "?"
}

func (b *Mysql) Configure(config map[string]string) {
	b.Dbname = config["SQL_DB"]
	b.Dbpass = os.Getenv("SQL_PASSWORD")
	b.Dbuser = config["SQL_USERNAME"]
	b.Server = config["SQL_SERVER"]
}

func (b *Mysql) Init() {
	b.Dbpass = os.Getenv("SQL_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", b.Dbuser, b.Dbpass, b.Server, b.Dbname))
	b.db = db
	s := Schema{}
	datetimec := s.AddColumn("last_switched", "DATETIME", "")
	datetimec.Wrap = "FROM_UNIXTIME(%v)"
	s.AddColumn("src_ip", "INT(4)", "UNSIGNED NOT NULL")
	s.AddColumn("src_port", "INT(2)", "UNSIGNED NOT NULL")
	s.AddColumn("dst_ip", "INT(4)", "UNSIGNED NOT NULL")
	s.AddColumn("dst_port", "INT(2)", "UNSIGNED NOT NULL")
	s.AddColumn("in_bytes", "INT(8)", "UNSIGNED NOT NULL")
	s.AddColumn("in_pkts", "INT(8)", "UNSIGNED NOT NULL")
	s.AddColumn("protocol", "INT(1)", "UNSIGNED NOT NULL DEFAULT 6")
	InitQuery := s.GetColumnStrings(INIT_TEMPLATE)

	b.schema = &s
	b.CheckSchema()
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	// Try and init the database
	_, err = db.Exec(fmt.Sprintf(USE_QUERY, b.Dbname))
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(InitQuery)
	if err != nil {
		panic(err.Error())
	}
}

func (b *Mysql) Test() {
	err := b.db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func (b *Mysql) CheckSchema() {
	/*
	   Validates the SQL Schema matches
	*/
	var (
		field      sql.NullString
		FieldType  sql.NullString
		null       sql.NullString
		key        sql.NullString
		HasDefault sql.NullString
		extra      sql.NullString
	)
	// Existing columns
	ec := make(map[string]bool)
	// Columns to delete
	var dc []string
	// Columns to add
	//var ac []Column

	db := b.db
	rows, err := db.Query(CHECK_QUERY)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&field, &FieldType, &null, &key, &HasDefault, &extra)
		ec[field.String] = true
		if err != nil {
			panic(err.Error())
		}
		c := b.schema.GetColumn(field.String)
		if c == nil {
			dc = append(dc, c.Name)
		}
	}

	for _, col := range b.schema.columns {
		if !ec[col.Name] {
			fmt.Printf("Adding Missing col %v to schema\n", col.Name)
			_, err := db.Query(fmt.Sprintf(ALTER_QUERY, col.Name, col.Type, col.Options))
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func (b *Mysql) Add(values map[uint16]fields.Value) {
	db := b.db
	InsertQuery := b.schema.InsertQuery(INSERT_TEMPLATE)
	s, err := db.Prepare(InsertQuery)
	_, err = s.Exec(
		values[fields.TIMESTAMP].ToInt(),
		values[fields.IPV4_SRC_ADDR].ToInt(),
		values[fields.L4_SRC_PORT].ToInt(),
		values[fields.IPV4_DST_ADDR].ToInt(),
		values[fields.L4_DST_PORT].ToInt(),
		values[fields.IN_BYTES].ToInt(),
		values[fields.IN_PKTS].ToInt(),
		values[fields.PROTOCOL].ToInt(),
	)
	if err != nil {
		panic(err.Error())
	}
}

/*
Re-initilize the database by dropping, and then re-adding, the schema
This will remove all data within the DB.
*/
func (b *Mysql) Reinit() {
	db := b.db
	_, err := db.Exec(fmt.Sprintf(USE_QUERY, b.Dbname))
	_, err = db.Exec(DROP_QUERY)
	if err != nil {
		panic(err.Error())
	}
	InitQuery := b.schema.GetColumnStrings(INIT_TEMPLATE)
	_, err = db.Exec(InitQuery)
	if err != nil {
		panic(err.Error())
	}
}
