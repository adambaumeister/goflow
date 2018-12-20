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

type Column interface {
	InsertValue(value fields.Value) string
	GetName() string
	GetType() string
	GetOptions() string
	Init() string
}

/*
#
INTCOLMUMN
#
*/
type IntColumn struct {
	Name    string
	Options string
	Type    string
	Wrap    string
	Field   uint16
}

func (c *IntColumn) Init() string {
	return fmt.Sprintf("%v %v %v", c.Name, c.Type, c.Options)
}
func (c *IntColumn) GetName() string {
	return c.Name
}
func (c *IntColumn) GetType() string {
	return c.Type
}
func (c *IntColumn) GetOptions() string {
	return c.Options
}

/*
Insert
Retrieves the string required to insert a value in a column using a pre-prepared statement
*/
func (c *IntColumn) insert() string {
	if len(c.Wrap) > 0 {
		return fmt.Sprintf(c.Wrap, "?")
	}
	return "?"
}

/*
InsertValue
Retrieves the string required to insert a value in a column using normal insert statement
*/
func (c *IntColumn) InsertValue(v fields.Value) string {
	if len(c.Wrap) > 0 {
		return fmt.Sprintf(c.Wrap, v.ToInt())
	}
	return fmt.Sprintf("%v", v.ToInt())
}

type Schema struct {
	columns     []Column
	columnIndex map[uint16]Column
}

func (s *Schema) AddColumn(f uint16, n string, t string, o string) *IntColumn {
	c := IntColumn{
		Name:    n,
		Options: o,
		Type:    t,
		Field:   f,
	}
	s.columns = append(s.columns, &c)
	s.columnIndex[f] = &c
	return &c
}
func (s *Schema) GetColumnStrings(t string) string {
	var qs []string
	for _, col := range s.columns {
		qs = append(qs, col.Init())
	}
	return fmt.Sprintf(t, strings.Join(qs, ", "))
}

func (s *Schema) GetColumn(c string) Column {
	for _, col := range s.columns {
		if col.GetName() == c {
			return col
		}
	}
	return nil
}

func (s *Schema) InsertQuery(t string, v map[uint16]fields.Value) string {
	var cols []string
	var vals []string
	for f, val := range v {
		var insertColumn string
		var insertValue string

		// Only add fields that we have configured the schema for
		if col, ok := s.columnIndex[f]; ok {
			insertColumn = col.GetName()
			insertValue = col.InsertValue(val)

			cols = append(cols, insertColumn)
			vals = append(vals, insertValue)
		}
	}
	return fmt.Sprintf(t, strings.Join(cols, ", "), strings.Join(vals, ", "))
}

func (s *Schema) InsertQueryFields() string {
	var qs []string
	for _, col := range s.columns {
		qs = append(qs, col.GetName())
	}
	return strings.Join(qs, ", ")
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
	s := Schema{
		columnIndex: make(map[uint16]Column),
	}
	datetimec := s.AddColumn(fields.TIMESTAMP, "last_switched", "DATETIME", "")
	datetimec.Wrap = "FROM_UNIXTIME(%v)"
	s.AddColumn(fields.IPV4_SRC_ADDR, "src_ip", "INT(4)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.L4_SRC_PORT, "src_port", "INT(2)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.IPV4_DST_ADDR, "dst_ip", "INT(4)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.L4_DST_PORT, "dst_port", "INT(2)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.IN_BYTES, "in_bytes", "INT(8)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.IN_PKTS, "in_pkts", "INT(8)", "UNSIGNED NOT NULL")
	s.AddColumn(fields.PROTOCOL, "protocol", "INT(1)", "UNSIGNED NOT NULL DEFAULT 6")
	InitQuery := s.GetColumnStrings(INIT_TEMPLATE)

	b.schema = &s

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
	b.CheckSchema()
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
	//var ac []IntColumn

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
			dc = append(dc, c.GetName())
		}
	}

	for _, col := range b.schema.columns {
		if !ec[col.GetName()] {
			fmt.Printf("Adding Missing col %v to schema\n", col.GetName())
			_, err := db.Query(fmt.Sprintf(ALTER_QUERY, col.GetName(), col.GetType(), col.GetOptions()))
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func (b *Mysql) Add(values map[uint16]fields.Value) {
	db := b.db
	InsertQuery := b.schema.InsertQuery(INSERT_TEMPLATE, values)
	fmt.Printf("query: 	%v", InsertQuery)
	_, err := db.Exec(InsertQuery)
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
