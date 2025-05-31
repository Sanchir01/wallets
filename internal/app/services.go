package app

import "github.com/Sanchir01/wallets/internal/feature/wallets"

type Services struct {
	WalletService *wallets.Service
}

func NewServices(repos *Repositories, db *Database) *Services {
	return &Services{
		WalletService: wallets.NewService(repos.WalletRepository, db.PrimaryDB),
	}
}
