package httphandlers

import (
	"net/http"

	_ "github.com/Sanchir01/wallets/docs"
	"github.com/Sanchir01/wallets/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func StartHTTTPHandlers(handlers *app.Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer)
	router.Use(PrometheusMiddleware)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/wallets", handlers.WalletHandler.GetAllWallets)
		r.Get("/wallets/{id}", handlers.WalletHandler.GetBalanceWallets)
		r.Post("/wallets/create", handlers.WalletHandler.CreateWallet)
		r.Post("/wallet", handlers.WalletHandler.SendMoney)
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return router
}
func StartPrometheusHandlers() http.Handler {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return router
}
