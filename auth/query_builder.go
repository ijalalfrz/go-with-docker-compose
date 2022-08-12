package auth

import "fmt"

type query struct {
	query  string
	params []interface{}
}

func NewAuthQuery() *query {
	return &query{}
}

func (q *query) BaseQueryGetAll() *query {
	q.query = `SELECT 
			u.id,
			u.name ,
			u.user_name ,
			u.password as user_password,
			u.created_at as user_created_at,
			u.created_by as user_created_by,
			u.updated_at as user_updated_at,
			u.updated_by as user_updated_by
		FROM Users u
		`
	return q
}

func (q *query) AddWhereClause() *query {
	q.query += " WHERE"
	return q
}
func (q *query) AddFilter(column string, value interface{}) *query {
	q.query += fmt.Sprintf(" %s=?", column)
	q.params = append(q.params, value)
	return q
}

func (q *query) AddAndClause() *query {
	q.query += " AND "
	return q
}

func (q *query) AddRangeFilter(column string) *query {
	q.query += fmt.Sprintf(" %s >= ? AND %s <= ?", column, column)
	return q
}

func (q *query) AddOrderBy(column string, order string) *query {
	q.query += fmt.Sprintf(" ORDER BY %s %s", column, order)
	return q
}
