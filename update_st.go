package query

import (
	"bytes"
	"database/sql"
	"strings"
)

type UpdateQuery struct {
	Table     string
	Columns   []string
	Values    []string
	Condition string
}

func Update(table string) *UpdateQuery {
	return &UpdateQuery{Table: table}
}

func (uq *UpdateQuery) Set(values map[string]string) *UpdateQuery {
	for k, v := range values {
		uq.Columns = append(uq.Columns, k)
		uq.Values = append(uq.Values, v)
	}
	return uq
}

func (uq *UpdateQuery) SetCol(column, value string) *UpdateQuery {
	uq.Columns = append(uq.Columns, column)
	uq.Values = append(uq.Values, value)
	return uq
}

func (uq *UpdateQuery) Where(where string) *UpdateQuery {
	if len(uq.Condition) == 0 {
		uq.Condition = where
	} else {
		uq.Condition += " AND " + where
	}
	return uq
}

func (uq *UpdateQuery) ToSql() string {
	buffer := bytes.Buffer{}

	buffer.WriteString("UPDATE " + uq.Table + " SET(")
	buffer.WriteString(strings.Join(uq.Columns, ", "))
	buffer.WriteString(") VALUES (")
	buffer.WriteString(strings.Join(uq.Values, ", "))
	buffer.WriteString(") ")
	if len(uq.Condition) != 0 {
		buffer.WriteString(" WHERE " + uq.Condition)
	}

	return buffer.String()
}

func (uq *UpdateQuery) Exec(conn Querier, args interface{}) (sql.Result, error) {
	query, args, err := compileQuery(conn, uq, args)
	if err != nil {
		return nil, err
	}
	return conn.Exec(query, args)
}
