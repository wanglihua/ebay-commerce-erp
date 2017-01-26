package core

import (
	"bytes"
	"fmt"
)

type SqlBuilder struct {
	select_clause   string
	from_clause     string
	where_clause    string
	order_by_clause string
	limit_clause    string
}

func NewSqlBuilder() *SqlBuilder {
	builder := new(SqlBuilder)
	builder.select_clause = ""
	builder.from_clause = ""
	builder.where_clause = ""
	builder.limit_clause = ""
	builder.order_by_clause = ""

	return builder
}

// .Select("*")
func (builder *SqlBuilder) Select(sql_select string) *SqlBuilder {
	//开始和结束都是换行符
	sql := `
SELECT %s
`
	builder.select_clause = fmt.Sprintf(sql, sql_select)
	return builder
}

// .From("table1")
func (builder *SqlBuilder) From(from string) *SqlBuilder {
	//开始和结束都是换行符
	sql := `
FROM %s
`
	builder.from_clause = fmt.Sprintf(sql, from)
	return builder
}

// .Where("id = 1")
func (builder *SqlBuilder) Where(condition string) *SqlBuilder {
	//开始和结束都是换行符
	sql := `
WHERE %s
`
	builder.where_clause = fmt.Sprintf(sql, condition)
	return builder
}

// .OrderBy("id ASC")
func (builder *SqlBuilder) OrderBy(order string) *SqlBuilder {
	//开始和结束都是换行符
	sql := `
ORDER BY %s
`
	builder.order_by_clause = fmt.Sprintf(sql, order)
	return builder
}

func (builder *SqlBuilder) Limit(skip int64, take int64) *SqlBuilder {
	//开始和结束都是换行符
	sql := `
LIMIT %d, %d
`
	builder.limit_clause = fmt.Sprintf(sql, skip, take)
	return builder
}

func (builder *SqlBuilder) ToSql() string {
	var buffer bytes.Buffer

	if builder.select_clause == "" {
		panic("sql select 语句不能为空！")
	} else {
		buffer.WriteString(builder.select_clause)
	}

	if builder.from_clause == "" {
		panic("sql from 语句不能为空！")
	} else {
		buffer.WriteString(builder.from_clause)
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
