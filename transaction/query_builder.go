package transaction

import "fmt"

type query struct {
	query string
}

func NewTransactionQuery() *query {
	return &query{}
}

func (q *query) BaseQueryGetGroupOmzetPerOutlet() *query {
	q.query = `
	WITH t as (
		SELECT 
			t.outlet_id, 
			COUNT(t.outlet_id) as number_of_transaction ,
			SUM(t.bill_total) as omzet,
			DAY(t.created_at) as transaction_date,
			MONTH(t.created_at) as transaction_month, 
			YEAR(t.created_at) as transaction_year,
			DATE(t.created_at) as full_date
		FROM Transactions t
		INNER JOIN Outlets o ON outlet_id = o.id
		GROUP BY t.outlet_id, o.merchant_id,DAY(t.created_at), MONTH(t.created_at), YEAR(t.created_at)
	)
	SELECT 
		t.outlet_id,
		m.merchant_name,
		o.outlet_name,
		t.number_of_transaction,
		t.omzet,
		t.transaction_date,
		t.transaction_month,
		t.transaction_year,
		t.full_date
	FROM t
	INNER JOIN Outlets o ON t.outlet_id = o.id
	INNER JOIN Merchants m ON o.merchant_id = m.id
	INNER JOIN Users u ON m.user_id = u.id
	`
	return q
}

func (q *query) BaseQueryGetGroupOmzetPerMerchant() *query {
	q.query = `
	WITH t as (
		SELECT
			o.merchant_id ,
			COUNT(t.outlet_id) as number_of_transaction ,
			SUM(t.bill_total) as omzet,
			DAY(t.created_at) as transaction_date,
			MONTH(t.created_at) as transaction_month, 
			YEAR(t.created_at) as transaction_year,
			DATE(t.created_at) as full_date
		FROM Transactions t
		INNER JOIN Outlets o ON outlet_id = o.id
		GROUP BY o.merchant_id, DAY(t.created_at), MONTH(t.created_at), YEAR(t.created_at)
	)
	SELECT 
		t.merchant_id,
		m.merchant_name,
		t.number_of_transaction,
		t.omzet,
		t.transaction_date,
		t.transaction_month,
		t.transaction_year,
		t.full_date
	FROM t
	INNER JOIN Merchants m ON t.merchant_id = m.id
	INNER JOIN Users u ON m.user_id = u.id
	`
	return q
}

func (q *query) BaseQueryGetAll() *query {
	q.query = `SELECT 
			u.id as user_id,
			u.name ,
			u.user_name ,
			u.password as user_password,
			u.created_at as user_created_at,
			u.created_by as user_created_by,
			u.updated_at as user_updated_at,
			u.updated_by as user_updated_by,
			m.id as merchant_id ,
			m.user_id as merchant_user_id ,
			m.merchant_name,
			m.created_at as merchant_created_at ,
			m.created_by as merchant_created_by ,
			m.updated_at as merchant_updated_at ,
			m.updated_by as merchant_updated_by ,
			o.id as outlet_id,
			o.merchant_id as outlet_merchant_id,
			o.outlet_name,
			o.created_at as outlet_created_at,
			o.created_by as outlet_created_by,
			o.updated_at as outlet_updated_at,
			o.updated_by as outlet_updated_by,
			t.id, 
			t.outlet_id, 
			t.bill_total, 
			t.created_at, 
			t.created_by, 
			t.updated_at, 
			t.updated_by
		FROM Transactions t
		INNER JOIN Outlets o ON t.outlet_id = o.id
		INNER JOIN Merchants m ON o.merchant_id = m.id
		INNER JOIN Users u ON m.user_id = u.id
		`
	return q
}

func (q *query) AddWhereClause() *query {
	q.query += " WHERE"
	return q
}
func (q *query) AddFilter(column string) *query {
	q.query += fmt.Sprintf(" %s=?", column)
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

func (q *query) Build() string {
	return q.query
}
