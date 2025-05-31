package wallets

import (
	"context"

	"errors"

	sq "github.com/Masterminds/squirrel"
	api "github.com/Sanchir01/wallets/pkg/lib/api/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllWallets(ctx context.Context) ([]*Wallet, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id,currency,balance,created_at,updated_at").
		From("wallets").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, api.ErrQueryString
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*Wallet
	for rows.Next() {
		var wallet Wallet
		err := rows.Scan(&wallet.ID, &wallet.Currency, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func (r *Repository) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("balance").
		From("wallets").
		Where(sq.Eq{"id": walletID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, api.ErrQueryString
	}

	var balance float64
	err = conn.QueryRow(ctx, query, args...).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *Repository) UpdateBalance(ctx context.Context, walletID uuid.UUID, balance float64, tx pgx.Tx) error {
	query, args, err := sq.
		Update("wallets").
		Set("balance", balance).
		Where(sq.Eq{"id": walletID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return api.ErrQueryString
	}
	cmdTag, err := tx.Exec(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return api.ErrNotFoundById
	}
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return api.ErrNotFoundById
	}
	return nil
}

func (r *Repository) CreateWallet(ctx context.Context, currency string, balance float64, tx pgx.Tx) error {
	query, args, err := sq.
		Insert("wallets").
		Columns("currency", "balance").
		Values(currency, balance).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return api.ErrQueryString
	}
	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return err
	}
	return nil
}

func (r *Repository) WriteTransaction(ctx context.Context, senderID, receiverID uuid.UUID, typemoney OperationType, amount float64, tx pgx.Tx) error {
	query, args, err := sq.
		Insert("transactions").
		Columns("wallet_id", "sender_wallet_id", "type", "amount").
		Values(receiverID, senderID, typemoney, amount).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return api.ErrQueryString
	}

	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return err
	}

	return nil
}
