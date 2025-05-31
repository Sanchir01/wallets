package app

import "github.com/Sanchir01/wallets/internal/feature/wallets"

type Repositories struct {
	WalletRepository *wallets.Repository
}

func NewRepositories(databases *Database) *Repositories {
	return &Repositories{
		WalletRepository: wallets.NewRepository(databases.PrimaryDB),
	}
}
