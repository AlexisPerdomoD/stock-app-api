package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

const GET_USER_QUERY = `
	SELECT 
		id, 
		username, 
		firstname, 
		lastname, 
		active, 
		created_at
	FROM users`
const GET_USER_PASSWORD_INCLUDED_QUERY = `
	SELECT 
		id, 
		username, 
		firstname, 
		lastname, 
		password,
		active, 
		created_at
	FROM users`
const INSERT_USER_QUERY = `
	INSERT INTO users (
		username, 
		firstname, 
		lastname, 
		password, 
		active
	) VALUES ($1, $2, $3, $4, $5) 
	RETURNING 
		id, 
		username, 
		firstname, 
		lastname, 
		active, 
		created_at`

const HAS_USER_STOCK_QUERY = `SELECT 1 FROM stock_users WHERE user_id=$1 AND stock_id=$2`
const INSERT_USER_STOCK_QUERY = `INSERT INTO stock_users (user_id, stock_id) VALUES ($1, $2)`
const DELETE_USER_STOCK_QUERY = `DELETE FROM stock_users WHERE user_id=$1 AND stock_id=$2`

type UserRepository struct {
	db sqlx.ExtContext
}

func (r *UserRepository) get(
	ctx context.Context,
	userID uint64,
	includePassword bool,
) (*domain.User, error) {
	record := &userRecord{}

	var q string
	if includePassword {
		q = GET_USER_PASSWORD_INCLUDED_QUERY + " WHERE id=$1"
	} else {
		q = GET_USER_QUERY + " WHERE id=$1"
	}

	if err := r.db.QueryRowxContext(ctx, q, userID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *UserRepository) getByUsername(
	ctx context.Context,
	username string,
	includePassword bool,
) (*domain.User, error) {
	record := &userRecord{}

	var q string
	if includePassword {
		q = GET_USER_PASSWORD_INCLUDED_QUERY + " WHERE username=$1 LIMIT 1"
	} else {
		q = GET_USER_QUERY + " WHERE username=$1 LIMIT 1"
	}

	if err := r.db.QueryRowxContext(ctx, q, username).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID uint64) (*domain.User, error) {
	return r.get(ctx, userID, false)
}

func (r *UserRepository) GetByIDWithPassword(ctx context.Context, userID uint64) (*domain.User, error) {
	return r.get(ctx, userID, true)
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.getByUsername(ctx, username, false)
}

func (r *UserRepository) GetByUsernameWithPassword(ctx context.Context, username string) (*domain.User, error) {
	return r.getByUsername(ctx, username, true)
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	if user == nil {
		return nil // no-op
	}

	record := &userRecord{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Password:  user.Password,
		Active:    user.Active,
	}
	q := INSERT_USER_QUERY
	if err := r.db.QueryRowxContext(ctx, q,
		record.Username,
		record.Firstname,
		record.Lastname,
		record.Password,
		record.Active).
		StructScan(record); err != nil {
		return err
	}

	record.MapDomain(user)
	return nil
}

func (r *UserRepository) HasUserStock(
	ctx context.Context,
	userID uint64,
	stockID uint64,
) (bool, error) {
	var one int
	if err := r.db.QueryRowxContext(
		ctx,
		HAS_USER_STOCK_QUERY,
		userID,
		stockID,
	).Scan(&one); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *UserRepository) RegisterUserStock(
	ctx context.Context,
	userID uint64,
	stockID uint64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		INSERT_USER_STOCK_QUERY,
		userID,
		stockID,
	)
	return err
}

func (r *UserRepository) RemoveUserStock(ctx context.Context, userID uint64, stockID uint64) error {
	_, err := r.db.ExecContext(
		ctx,
		DELETE_USER_STOCK_QUERY,
		userID,
		stockID,
	)
	return err
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}
	return &UserRepository{db}
}
