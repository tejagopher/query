package query

import (
	"bytes"
	"database/sql"
	"strings"
)

type InsertQuery struct {
	Table   string
	Columns []string
	Values  []string
	suffix  string
}

func Insert(table string) *InsertQuery {
	return &InsertQuery{Table: table}
}

func (iq *InsertQuery) Set(values map[string]string) *InsertQuery {
	for k, v := range values {
		iq.Columns = append(iq.Columns, k)
		iq.Values = append(iq.Values, v)
	}
	return iq
}

func (iq *InsertQuery) SetCol(column, value string) *InsertQuery {
	iq.Columns = append(iq.Columns, column)
	iq.Values = append(iq.Values, value)
	return iq
}

func (iq *InsertQuery) Suffix(value string) *InsertQuery {
	iq.suffix = value
	return iq
}

func (iq *InsertQuery) ToSql() string {
	buffer := bytes.Buffer{}

	buffer.WriteString("INSERT INTO " + iq.Table + " (")
	buffer.WriteString(strings.Join(iq.Columns, ", "))
	buffer.WriteString(") VALUES (")
	buffer.WriteString(strings.Join(iq.Values, ", "))
	buffer.WriteString(") ")
	if len(iq.suffix) != 0 {
		buffer.WriteString(iq.suffix)
	}

	return buffer.String()
}

func (iq *InsertQuery) Exec(conn Querier, args interface{}) (sql.Result, error) {
	query, args, err := compileQuery(conn, iq, args)
	if err != nil {
		return nil, err
	}
	return conn.Exec(query, args)
}

func (iq *InsertQuery) GetOne(conn Querier, result interface{}, args interface{}) error {
	query, args, err := compileQuery(conn, iq, args)
	if err != nil {
		return err
	}
	return conn.Get(result, query, args)
}
