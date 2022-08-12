package auth

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/ijalalfrz/majoo-test-1/entity"
	"github.com/ijalalfrz/majoo-test-1/exception"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (user entity.User, err error)
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

// NewAuthRepository is a constructor
func NewAuthRepository(logger *logrus.Logger, db *sql.DB) Repository {
	return &repository{
		logger: logger,
		db:     db,
	}
}

func (r repository) FindByUsername(ctx context.Context, username string) (user entity.User, err error) {
	var cmd sqlCommand = r.db
	q := NewAuthQuery().
		BaseQueryGetAll().
		AddWhereClause().
		AddFilter("u.user_name", username)

	bunchOfUsers, err := r.query(ctx, cmd, q.query, q.params...)
	if err != nil {
		err = wrapError(err)
		return
	}

	length := len(bunchOfUsers)
	if length < 1 {
		err = exception.ErrNotFound
		return
	}
	user = bunchOfUsers[length-1]
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

func (r repository) query(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfUser []entity.User, err error) {
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
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.UserName,
			&user.Password,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.UpdatedAt,
			&user.UpdatedBy,
		)
		if err != nil {
			r.logger.Error(query, err)
			return
		}

		bunchOfUser = append(bunchOfUser, user)
	}
	return
}
