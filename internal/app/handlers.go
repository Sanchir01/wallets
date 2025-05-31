package app

import (
	"log/slog"

	"github.com/Sanchir01/wallets/internal/feature/wallets"
)

type Handlers struct {
	WalletHandler *wallets.Handler
}

func NewHandlers(services *Services, lg *slog.Logger) *Handlers {
	return &Handlers{
		WalletHandler: wallets.NewHandler(services.WalletService, lg),
	}
}
