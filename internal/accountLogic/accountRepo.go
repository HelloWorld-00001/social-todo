package accountLogic

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
)

type AccountRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (r *AccountRepository) CreateAccount(a *Account) error {
	query, args, err := r.qb.
		Insert("account").
		Columns("username", "hashpassword", "salt").
		Values(a.Username, a.HashPassword, a.Salt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *AccountRepository) GetByUsername(username string) (*Account, error) {
	query, args, err := r.qb.
		Select("id", "username", "hashpassword", "salt").
		From("account").
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(query, args...)

	acc := &Account{}
	err = row.Scan(&acc.ID, &acc.Username, &acc.HashPassword, &acc.Salt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return acc, nil
}
