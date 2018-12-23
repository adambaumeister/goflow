package mysql

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

type Mysql struct {
	Dbname string
	Dbpass string
	Dbuser string
	Server string
	db     *sql.DB

	schema *Schema

	CheckQuery    string
	AlterQuery    string
	InitQuery     string
	AddIndexQuery string
	InsertQuery   string
	DropQuery     string
	AlterColQuery string
}

type Column interface {
	InsertValue(value fields.Value) string
	GetName() string
	GetType() string
	GetOptions() string
	getFieldString() string
	Init() string
}

/*
#
INTCOLMUMN
Integer-based column and methods
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
func (c *IntColumn) getFieldString() string {
	return fmt.Sprintf("%v %v %v", c.GetName(), c.GetType(), c.GetOptions())
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

/*
#
Binary
Binary columns for binary data
#
*/
type BinaryColumn struct {
	Name    string
	Options string
	Type    string
	Wrap    string
	Field   uint16
}

func (c *BinaryColumn) Init() string {
	return fmt.Sprintf("%v %v %v", c.Name, c.Type, c.Options)
}
func (c *BinaryColumn) GetName() string {
	return c.Name
}
func (c *BinaryColumn) GetType() string {
	return c.Type
}
func (c *BinaryColumn) GetOptions() string {
	return c.Options
}

/*
InsertValue
Retrieves the string required to insert a value in a column using normal insert statement
*/
func (c *BinaryColumn) InsertValue(v fields.Value) string {
	if len(c.Wrap) > 0 {
		return fmt.Sprintf(c.Wrap, v.ToString())
	}
	return fmt.Sprintf("UNHEX(\"%v\")", v.ToString())
}

type Schema struct {
	columns     []Column
	columnIndex map[uint16]Column
}

func (c *BinaryColumn) getFieldString() string {
	return fmt.Sprintf("%v %v %v", c.GetName(), c.GetType(), c.GetOptions())
}

func (s *Schema) AddIntColumn(f uint16, n string, t string, o string) *IntColumn {
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
func (s *Schema) AddBinaryColumn(f uint16, n string, t string, o string) *BinaryColumn {
	c := BinaryColumn{
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

	b.CheckQuery = "SHOW COLUMNS IN goflow_records;"
	b.AlterQuery = "ALTER TABLE goflow_records ADD COLUMN %v %v %v;"
	b.InitQuery = "CREATE TABLE IF NOT EXISTS goflow_records (%v, INDEX last_switched_idx (last_switched));"
	b.AddIndexQuery = "ALTER TABLE goflow_records ADD INDEX last_switched_idx (last_switched)"
	b.InsertQuery = "INSERT INTO goflow_records (%v) VALUES (%v);"
	b.DropQuery = "DROP TABLE goflow_records"
	b.AlterColQuery = "ALTER TABLE goflow_records MODIFY COLUMN %v"

	b.Dbpass = os.Getenv("SQL_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", b.Dbuser, b.Dbpass, b.Server, b.Dbname))
	b.db = db
	s := Schema{
		columnIndex: make(map[uint16]Column),
	}
	datetimec := s.AddIntColumn(fields.TIMESTAMP, "last_switched", "datetime", "NOT NULL")
	datetimec.Wrap = "FROM_UNIXTIME(%v)"
	s.AddIntColumn(fields.IPV4_SRC_ADDR, "src_ip", "int(4)", "unsigned DEFAULT NULL")
	s.AddIntColumn(fields.L4_SRC_PORT, "src_port", "int(2)", "unsigned NOT NULL")
	s.AddIntColumn(fields.IPV4_DST_ADDR, "dst_ip", "int(4)", "unsigned DEFAULT NULL")
	s.AddIntColumn(fields.L4_DST_PORT, "dst_port", "int(2)", "unsigned NOT NULL")
	s.AddIntColumn(fields.IN_BYTES, "in_bytes", "int(8)", "unsigned NOT NULL")
	s.AddIntColumn(fields.IN_PKTS, "in_pkts", "int(8)", "unsigned NOT NULL")
	s.AddIntColumn(fields.PROTOCOL, "protocol", "int(1)", "unsigned NOT NULL")
	s.AddBinaryColumn(fields.IPV6_SRC_ADDR, "src_ipv6", "varbinary(16)", "DEFAULT NULL")
	s.AddBinaryColumn(fields.IPV6_DST_ADDR, "dst_ipv6", "varbinary(16)", "DEFAULT NULL")
	InitQuery := s.GetColumnStrings(b.InitQuery)

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
	ec := make(map[string]string)
	// Columns to delete
	var dc []string
	// Columns to add
	//var ac []IntColumn

	db := b.db
	rows, err := db.Query(b.CheckQuery)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&field, &FieldType, &null, &key, &HasDefault, &extra)
		if err != nil {
			panic(err.Error())
		}
		var fieldString string
		if null.String == "YES" {
			fieldString = fmt.Sprintf("%v %v DEFAULT NULL", field.String, FieldType.String)
		} else {
			fieldString = fmt.Sprintf("%v %v NOT NULL", field.String, FieldType.String)
		}

		ec[field.String] = fieldString

		c := b.schema.GetColumn(field.String)
		if c == nil {
			dc = append(dc, c.GetName())
		}

		if field.String == "last_switched" {
			if key.String != "MUL" {
				fmt.Printf("Bad index %v: %v\n", field.String, extra.String)
				_, err := db.Query(b.AddIndexQuery)
				if err != nil {
					panic(err.Error())
				}
			}

		}
	}

	fmt.Print("Undergoing schema check. If changes are found, this may take a while...\n")
	for _, col := range b.schema.columns {
		// Check column exists
		if fs, ok := ec[col.GetName()]; ok {
			// If it does, check it matches what it's sposed ta be
			if fs != col.getFieldString() {
				fmt.Printf("Field mismatch in schema: %v (%v should be %v)\n", col.GetName(), fs, col.getFieldString())
				_, err := db.Query(fmt.Sprintf(b.AlterColQuery, col.getFieldString()))
				if err != nil {
					panic(err.Error())
				}

			}
		} else {
			fmt.Printf("Adding Missing col %v to schema\n", col.GetName())
			_, err := db.Query(fmt.Sprintf(b.AlterQuery, col.GetName(), col.GetType(), col.GetOptions()))
			if err != nil {
				panic(err.Error())
			}
		}
	}

	fmt.Print("Schema check done!\n")
}

func (b *Mysql) Add(values map[uint16]fields.Value) {
	db := b.db
	InsertQuery := b.schema.InsertQuery(b.InsertQuery, values)
	//fmt.Printf("query: 	%v\n", InsertQuery)
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
	_, err = db.Exec(b.DropQuery)
	if err != nil {
		panic(err.Error())
	}
	InitQuery := b.schema.GetColumnStrings(b.InitQuery)
	_, err = db.Exec(InitQuery)
	if err != nil {
		panic(err.Error())
	}
}
