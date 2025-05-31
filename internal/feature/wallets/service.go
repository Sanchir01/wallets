package wallets

import (
	"context"
	"errors"

	api "github.com/Sanchir01/wallets/pkg/lib/api/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo      *Repository
	primaryDB *pgxpool.Pool
}

func NewService(repo *Repository, primaryDB *pgxpool.Pool) *Service {
	return &Service{
		repo:      repo,
		primaryDB: primaryDB,
	}
}
func (s *Service) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	balance, err := s.repo.GetBalance(ctx, walletID)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
func (s *Service) GetAllWallets(ctx context.Context) ([]*Wallet, error) {
	wallets, err := s.repo.GetAllWallets(ctx)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (s *Service) CreateWallet(ctx context.Context, currency string, balance float64) error {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()
	err = s.repo.CreateWallet(ctx, currency, balance, tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) SendMoney(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64, typeoperation OperationType) error {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()
	senderbalance, err := s.repo.GetBalance(ctx, senderID)
	if err != nil {
		return err
	}
	if (senderbalance - amount) < 0 {
		return api.ErrInsufficientBalance
	}
	if err := s.repo.UpdateBalance(ctx, senderID, senderbalance-amount, tx); err != nil {
		return err
	}

	receiverbalance, err := s.repo.GetBalance(ctx, receiverID)
	if err := s.repo.UpdateBalance(ctx, receiverID, receiverbalance+amount, tx); err != nil {
		return err
	}
	if err := s.repo.WriteTransaction(ctx, senderID, receiverID, typeoperation, amount, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
