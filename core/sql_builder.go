package core

import (
	"bytes"
	"fmt"
	"strings"
)

type SqlBuilder struct {
	select_clause   string
	count_clause    string
	from_clause     string
	where_clause    string
	group_by_clause string
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
	//结束是换行符
	sql := `SELECT %s
`
	builder.select_clause = fmt.Sprintf(sql, sql_select)
	return builder
}

// .Count("count(*) as user_count")
func (builder *SqlBuilder) Count(sql_count string) *SqlBuilder {
	//结束是换行符
	sql := `SELECT %s
`
	builder.count_clause = fmt.Sprintf(sql, sql_count)
	return builder
}

// .From("table1")
func (builder *SqlBuilder) From(from string) *SqlBuilder {
	//结束是换行符
	sql := `FROM %s
`
	builder.from_clause = fmt.Sprintf(sql, from)
	return builder
}

// .Where("id = 1")
func (builder *SqlBuilder) Where(condition string) *SqlBuilder {
	//结束是换行符
	sql := `WHERE %s
`
	builder.where_clause = fmt.Sprintf(sql, condition)
	return builder
}

// .GroupBy("type")
func (builder *SqlBuilder) GroupBy(group_by string) *SqlBuilder {
	//结束是换行符
	sql := `GROUP BY %s
`
	builder.group_by_clause = fmt.Sprintf(sql, group_by)
	return builder
}

// .OrderBy("id ASC")
func (builder *SqlBuilder) OrderBy(order string) *SqlBuilder {
	//结束是换行符
	sql := `ORDER BY %s
`
	builder.order_by_clause = fmt.Sprintf(sql, order)
	return builder
}

func (builder *SqlBuilder) Limit(skip int64, take int64) *SqlBuilder {
	//结束是换行符
	sql := `LIMIT %d, %d
`
	builder.limit_clause = fmt.Sprintf(sql, skip, take)
	return builder
}

func (builder *SqlBuilder) CloneWithoutLimit() *SqlBuilder {
	var new_sql_builder SqlBuilder = *builder
	new_sql_builder.limit_clause = ""
	return &new_sql_builder
}

type sql_type int64

const (
	sql_type_get   sql_type = 1
	sql_type_list  sql_type = 2
	sql_type_count sql_type = 3
)

func (builder *SqlBuilder) to_sql(sql_type sql_type) string {
	var buffer bytes.Buffer

	if sql_type == sql_type_count && builder.group_by_clause != "" {

		if builder.count_clause == "" {
			panic("sql count 语句不能为空！")
		} else {
			buffer.WriteString(builder.count_clause)
		}

		buffer.WriteString(" FROM ( ")

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
		buffer.WriteString(builder.group_by_clause)
		buffer.WriteString(" ) count_alias")
	} else {
		if sql_type == sql_type_list || sql_type == sql_type_get {
			if builder.select_clause == "" {
				panic("sql select 语句不能为空！")
			} else {
				buffer.WriteString(builder.select_clause)
			}
		} else if sql_type == sql_type_count {
			if builder.count_clause == "" {
				panic("sql count 语句不能为空！")
			} else {
				buffer.WriteString(builder.count_clause)
			}
		}

		if builder.from_clause == "" {
			panic("sql from 语句不能为空！")
		} else {
			buffer.WriteString(builder.from_clause)
		}

		if builder.where_clause != "" {
			buffer.WriteString(builder.where_clause)
		}

		if builder.group_by_clause != "" {
			buffer.WriteString(builder.group_by_clause)
		}

		if sql_type == sql_type_list || sql_type == sql_type_get {
			if builder.order_by_clause != "" {
				buffer.WriteString(builder.order_by_clause)
			}
		}

		if sql_type == sql_type_list {
			if builder.limit_clause != "" {
				buffer.WriteString(builder.limit_clause)
			}
		} else if sql_type == sql_type_get {
			sql_limit_one := `LIMIT 0, 1
`
			buffer.WriteString(sql_limit_one)
		}
	}

	return buffer.String()
}

func (builder *SqlBuilder) SqlGet() string {
	return builder.to_sql(sql_type_get)
}

func (builder *SqlBuilder) SqlList() string {
	return builder.to_sql(sql_type_list)
}

func (builder *SqlBuilder) SqlCount() string {
	return builder.to_sql(sql_type_count)
}

func SqlAnd(exps ...string) string {
	//前后表达式都需要 ( )
	exp_list := make([]string, 0, len(exps))
	for _, exp := range exps {
		exp_list = append(exp_list, fmt.Sprintf("(%s)", exp))
	}

	return strings.Join(exp_list, " AND ")
}

func SqlOr(exps ...string) string {
	//前后表达式都需要 ( )
	exp_list := make([]string, 0, len(exps))
	for _, exp := range exps {
		exp_list = append(exp_list, fmt.Sprintf("(%s)", exp))
	}

	return strings.Join(exp_list, " OR ")
}
