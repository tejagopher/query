package query

import (
	"bytes"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ToSql interface {
	ToSql() string
}

type FromQuery struct {
	Table string
}

func From(table string) *FromQuery {
	return &FromQuery{Table: table}
}

func (fq *FromQuery) Select(columns ...string) *Select {
	return &Select{Table: fq.Table, Cols: columns}
}

type Select struct {
	Table     string
	Cols      []string
	Condition string
	Order     []string
}

func (sq *Select) Select(columns ...string) *Select {
	sq.Cols = append(sq.Cols, columns...)
	return sq
}

func (sq *Select) Where(where string) *Select {
	if len(sq.Condition) == 0 {
		sq.Condition = where
	} else {
		sq.Condition += " AND " + where
	}
	return sq
}

func (sq *Select) OrderBy(order string) *Select {
	sq.Order = append(sq.Order, order)
	return sq
}

func (sq *Select) ToSql() string {
	buffer := bytes.Buffer{}

	buffer.WriteString("SELECT ")
	buffer.WriteString(strings.Join(sq.Cols, ", "))
	buffer.WriteString(" FROM " + sq.Table)
	if len(sq.Condition) != 0 {
		buffer.WriteString(" WHERE " + sq.Condition)
	}
	if len(sq.Order) != 0 {
		buffer.WriteString(" ORDER BY " + strings.Join(sq.Order, ","))
	}

	return buffer.String()
}

func (sq *Select) GetOne(conn Querier, result interface{}, args interface{}) error {
	query, args, err := compileQuery(conn, sq, args)
	if err != nil {
		return err
	}
	return conn.Get(result, query, args)
}

func (sq *Select) GetMany(conn Querier, result interface{}, args interface{}) error {
	query, args, err := compileQuery(conn, sq, args)
	if err != nil {
		return err
	}
	return conn.Select(result, query, args)
}

func compileQuery(conn Querier, st ToSql, namedArgs interface{}) (query string,
	positionalArgs interface{}, err error) {
	query, positionalArgs, err = sqlx.Named(st.ToSql(), namedArgs)
	if err != nil {
		return "", nil, err
	}
	query, positionalArgs, err = sqlx.In(query, positionalArgs)
	if err != nil {
		return "", nil, err
	}
	query = conn.Rebind(query)
	return query, positionalArgs, nil
}
