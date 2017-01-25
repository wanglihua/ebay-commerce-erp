package core

import (
	"bytes"
	"fmt"
)

type SqlBuilder struct {
	select_from_clause string
	where_clause       string
	limit_clause       string
	order_by_clause    string
}

func NewSqlBuilder() *SqlBuilder {
	builder := new(SqlBuilder)
	builder.select_from_clause = ""
	builder.where_clause = ""
	builder.limit_clause = ""
	builder.order_by_clause = ""

	return builder
}

// .SelectFrom("SELECT * FROM table1")
func (builder *SqlBuilder) SelectFrom(select_from string) *SqlBuilder {
	builder.select_from_clause = select_from
	return builder
}

// .Where("id = 1")
func (builder *SqlBuilder) Where(condition string) *SqlBuilder {
	//开始和结束都是换行符
	where := `
WHERE %s
`
	builder.where_clause = fmt.Sprintf(where, condition)
	return builder
}

// .OrderBy("id ASC")
func (builder *SqlBuilder) OrderBy(order string) *SqlBuilder {
	//开始和结束都是换行符
	order_by := `
ORDER BY %s
`
	builder.order_by_clause = fmt.Sprintf(order_by, order)
	return builder
}

func (builder *SqlBuilder) Limit(skip int64, take int64) *SqlBuilder {
	//开始和结束都是换行符
	limit := `
LIMIT %d, %d
`
	builder.limit_clause = fmt.Sprintf(limit, skip, take)
	return builder
}

func (builder *SqlBuilder) ToSql() string {
	var buffer bytes.Buffer

	if builder.select_from_clause == "" {
		panic("sql select from 语句不能为空！")
	} else {
		buffer.WriteString(builder.select_from_clause)
	}

	if builder.where_clause != "" {
		buffer.WriteString(builder.where_clause)
	}

	if builder.order_by_clause != "" {
		buffer.WriteString(builder.order_by_clause)
	}

	if builder.limit_clause != "" {
		buffer.WriteString(builder.limit_clause)
	}

	return buffer.String()
}

func SqlAnd(exp1, exp2 string) string {
	//前后表达式都需要 ( )
	return fmt.Sprintf("(%s) AND (%s)", exp1, exp2)
}

func SqlOr(exp1, exp2 string) string {
	//前后表达式都需要 ( )
	return fmt.Sprintf("(%s) OR (%s)", exp1, exp2)
}
