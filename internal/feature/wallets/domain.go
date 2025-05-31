package wallets

import (
	"time"

	api "github.com/Sanchir01/wallets/pkg/lib/api/response"
	"github.com/google/uuid"
)

type OperationType string

const (
	OperationTypeDeposit  OperationType = "DEPOSIT"
	OperationTypeWithdraw OperationType = "WITHDRAW"
	OperationTypeTransfer OperationType = "TRANSFER"
)

type Wallet struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Balance   float64   `db:"balance" json:"balance"`
	Currency  string    `db:"currency" json:"currency"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type TransactionWallet struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	WalletID       uuid.UUID  `db:"wallet_id" json:"wallet_id"`
	SenderWalletID *uuid.UUID `db:"sender_wallet_id" json:"sender_wallet_id"`
	Amount         float64    `db:"amount" json:"amount"`
	Type           string     `db:"type" json:"type"`
	Description    *string    `db:"description" json:"description"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

type GetAllWalletsResponse struct {
	Response api.Response `json:"response"`
	Wallets  []*Wallet    `json:"wallets"`
}
type GetBalanceResponse struct {
	Response api.Response `json:"response"`
	Balance  float64      `json:"balance"`
}

type CreateWalletRequest struct {
	Balance  float64 `json:"balance" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type CreateWalletResponse struct {
	Response api.Response `json:"response"`
	OK       string       `json:"balance"`
}

type SendCoinRequest struct {
	WalletID       uuid.UUID     `json:"wallet_id" validate:"required"`
	SenderWalletID *uuid.UUID    `json:"sender_id" validate:"omitempty"`
	Type           OperationType `json:"type" validate:"required"`
	Amount         float64       `json:"amount" validate:"required"`
}

type SendCoinResponse struct {
	Response api.Response `json:"response"`
	OK       string       `json:"balance"`
}
