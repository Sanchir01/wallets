package wallets

import (
	"encoding/json"
	"log/slog"
	"net/http"

	api "github.com/Sanchir01/wallets/pkg/lib/api/response"
	"github.com/Sanchir01/wallets/pkg/lib/logger/sl"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
	Log     *slog.Logger
}

func NewHandler(service *Service, lg *slog.Logger) *Handler {
	return &Handler{
		service: service,
		Log:     lg,
	}
}

// @Tags wallet
// @Description get all wallets
// @Accept json
// @Produce json
// @Success 200 {object}  GetAllWalletsResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /wallets [get]
func (h *Handler) GetAllWallets(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.GetAllWallets"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	wallets, err := h.service.GetAllWallets(r.Context())
	if err != nil {
		log.Error("fail get all users", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("get all wallets success")
	render.JSON(w, r, GetAllWalletsResponse{
		Response: api.OK(),
		Wallets:  wallets,
	})
}

// @Tags wallet
// @Description get balance wallet by id
// @Param id path string true "user id"
// @Accept json
// @Produce json
// @Success 200 {object}  GetBalanceResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /wallets/{id} [get]
func (h *Handler) GetBalanceWallets(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.GetBalanceWallets"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	id := chi.URLParam(r, "id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		log.Error("invalid UUID format", sl.Err(err))
		render.JSON(w, r, api.Error("invalid id format"))
		return
	}
	balance, err := h.service.GetBalance(r.Context(), uuidID)
	if err != nil {
		log.Error("fail get all users", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("get all wallets success")
	render.JSON(w, r, GetBalanceResponse{
		Response: api.OK(),
		Balance:  balance,
	})
}

// @Tags wallet
// @Description get balance wallet by id
// @Accept json
// @Produce json
// @Param input body CreateWalletRequest true "create wallet"
// @Success 200 {object}  CreateWalletResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /wallets/create [post]
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.CreateWallet"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации тела"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	err := h.service.CreateWallet(r.Context(), req.Currency, req.Balance)
	if err != nil {
		log.Error("fail get all users", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("create wallet success")
	render.JSON(w, r, CreateWalletResponse{
		Response: api.OK(),
		OK:       "create wallet success",
	})
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.SendMoney"
	log := h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.JSON(w, r, api.Error("Ошибка при валидации тела"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	err := h.service.SendMoney(r.Context(), req.SenderWalletID, req.WalletID, req.Amount, req.Type)
	if err != nil {
		log.Error("send money failed", sl.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	log.Info("send money success")
	render.JSON(w, r, SendCoinResponse{
		Response: api.OK(),
		OK:       "send money success",
	})
}
