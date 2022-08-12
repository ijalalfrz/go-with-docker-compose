package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ijalalfrz/majoo-test-1/entity"
	"github.com/ijalalfrz/majoo-test-1/exception"
	"github.com/ijalalfrz/majoo-test-1/model"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	FindManyByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time, limit int) (transactions []entity.Transaction, err error)
	FindOmzetPerMerchantByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time) (omzets []model.GroupOmzetPerMerchant, err error)
	FindOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time, skip int, limit int) (omzets []model.GroupOmzetPerOutlet, err error)
	CountOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time) (totalData int64, err error)
}

type sqlCommand interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type repository struct {
	logger *logrus.Logger
	db     *sql.DB
}

// NewTransactionRepository is a constructor
func NewTransactionRepository(logger *logrus.Logger, db *sql.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}
func (r repository) CountOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time) (totalData int64, err error) {
	var cmd sqlCommand = r.db
	q := NewTransactionQuery().
		BaseQueryGetCountGroupOmzetPerOutlet().
		AddWhereClause().
		AddFilter("m.user_id").
		AddAndClause().
		AddRangeFilter("t.full_date").
		Build()
	totalData, err = r.queryCountOmzetPerOutlet(ctx, cmd, q, userId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		err = wrapError(err)
		return
	}
	return
}

func (r repository) FindManyByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time, limit int) (transactions []entity.Transaction, err error) {
	var cmd sqlCommand = r.db
	q := NewTransactionQuery().
		BaseQueryGetAll().
		AddWhereClause().
		AddFilter("m.user_id").
		Build()
	bunchOfTransactions, err := r.query(ctx, cmd, q, userId)
	if err != nil {
		err = wrapError(err)
		return
	}

	length := len(bunchOfTransactions)
	if length < 1 {
		err = exception.ErrNotFound
		return
	}
	transactions = bunchOfTransactions
	return
}

func (r repository) FindOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time, skip int, limit int) (omzets []model.GroupOmzetPerOutlet, err error) {

	var cmd sqlCommand = r.db
	q := NewTransactionQuery().
		BaseQueryGetGroupOmzetPerOutlet().
		AddWhereClause().
		AddFilter("m.user_id").
		AddAndClause().
		AddRangeFilter("t.full_date").
		AddOrderBy("t.full_date", "asc").
		AddLimit(limit).
		AddOffset(skip).
		Build()
	bunchOfOmzet, err := r.queryOmzetPerOutlet(ctx, cmd, q, userId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		err = wrapError(err)
		return
	}

	length := len(bunchOfOmzet)
	if length < 1 {
		err = exception.ErrNotFound
		return
	}
	omzets = bunchOfOmzet
	return
}

func (r repository) FindOmzetPerMerchantByUserID(ctx context.Context, userId int64, startDate *time.Time, endDate *time.Time) (omzets []model.GroupOmzetPerMerchant, err error) {

	var cmd sqlCommand = r.db
	q := NewTransactionQuery().
		BaseQueryGetGroupOmzetPerMerchant().
		AddWhereClause().
		AddFilter("m.user_id").
		AddAndClause().
		AddRangeFilter("t.full_date").
		AddOrderBy("t.full_date", "asc").
		Build()
	bunchOfOmzet, err := r.queryOmzetPerMerchant(ctx, cmd, q, userId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		err = wrapError(err)
		return
	}

	length := len(bunchOfOmzet)
	if length < 1 {
		err = exception.ErrNotFound
		return
	}
	omzets = bunchOfOmzet
	return
}

func (r repository) queryOmzetPerMerchant(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfOmzet []model.GroupOmzetPerMerchant, err error) {
	var rows *sql.Rows

	if rows, err = cmd.QueryContext(ctx, query, args...); err != nil {

		r.logger.Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error(query, err)
		}
	}()

	for rows.Next() {
		var omzet model.GroupOmzetPerMerchant
		err = rows.Scan(
			&omzet.MerchantID,
			&omzet.MerchantName,
			&omzet.NumberOftransaction,
			&omzet.Omzet,
			&omzet.TransactionDate,
			&omzet.TransactionMonth,
			&omzet.TransactionYear,
			&omzet.TransactionFullDate,
		)
		if err != nil {
			r.logger.Error(query, err)
			return
		}

		bunchOfOmzet = append(bunchOfOmzet, omzet)
	}
	return
}

func (r repository) queryOmzetPerOutlet(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfOmzet []model.GroupOmzetPerOutlet, err error) {
	var rows *sql.Rows

	if rows, err = cmd.QueryContext(ctx, query, args...); err != nil {

		r.logger.Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error(query, err)
		}
	}()

	for rows.Next() {
		var omzet model.GroupOmzetPerOutlet
		err = rows.Scan(
			&omzet.OutletID,
			&omzet.OutletName,
			&omzet.MerchantName,
			&omzet.NumberOftransaction,
			&omzet.Omzet,
			&omzet.TransactionDate,
			&omzet.TransactionMonth,
			&omzet.TransactionYear,
			&omzet.TransactionFullDate,
		)
		if err != nil {
			r.logger.Error(query, err)
			return
		}

		bunchOfOmzet = append(bunchOfOmzet, omzet)
	}
	return
}

func (r repository) queryCountOmzetPerOutlet(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (total int64, err error) {
	var rows *sql.Rows

	if rows, err = cmd.QueryContext(ctx, query, args...); err != nil {

		r.logger.Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error(query, err)
		}
	}()

	for rows.Next() {
		err = rows.Scan(
			&total,
		)
		if err != nil {
			r.logger.Error(query, err)
			return
		}

	}
	return
}

func (r repository) query(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfTransactions []entity.Transaction, err error) {
	var rows *sql.Rows

	if rows, err = cmd.QueryContext(ctx, query, args...); err != nil {

		r.logger.Error(query, err)
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error(query, err)
		}
	}()

	for rows.Next() {
		var transaction entity.Transaction
		transaction.User = &entity.User{}
		transaction.Merchant = &entity.Merchant{}
		transaction.Outlet = &entity.Outlet{}
		err = rows.Scan(
			&transaction.User.ID,
			&transaction.User.Name,
			&transaction.User.UserName,
			&transaction.User.Password,
			&transaction.User.CreatedAt,
			&transaction.User.CreatedBy,
			&transaction.User.UpdatedAt,
			&transaction.User.UpdatedBy,
			&transaction.Merchant.ID,
			&transaction.Merchant.UserID,
			&transaction.Merchant.MerchantName,
			&transaction.Merchant.CreatedAt,
			&transaction.Merchant.CreatedBy,
			&transaction.Merchant.UpdatedAt,
			&transaction.Merchant.UpdatedBy,
			&transaction.Outlet.ID,
			&transaction.Outlet.MerchantID,
			&transaction.Outlet.OutletName,
			&transaction.Outlet.CreatedAt,
			&transaction.Outlet.CreatedBy,
			&transaction.Outlet.UpdatedAt,
			&transaction.Outlet.UpdatedBy,
			&transaction.ID,
			&transaction.OutletID,
			&transaction.BillTotal,
			&transaction.CreatedAt,
			&transaction.CreatedBy,
			&transaction.UpdatedAt,
			&transaction.UpdatedBy,
		)
		if err != nil {
			r.logger.Error(query, err)
			return
		}

		bunchOfTransactions = append(bunchOfTransactions, transaction)
	}
	return
}

func wrapError(e error) (err error) {
	if e == sql.ErrNoRows {
		return exception.ErrNotFound
	}
	if driverErr, ok := e.(*mysql.MySQLError); ok {
		if driverErr.Number == 1062 {
			return exception.ErrConflict
		}
	}
	return exception.ErrInternalServer
}
