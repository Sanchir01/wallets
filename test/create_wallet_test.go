package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/Sanchir01/wallets/internal/feature/wallets"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
)

const (
	host              = "localhost:8080"
	createWalletPath  = "/wallets/create"
	getAllWalletsPath = "/wallets"
	getWalletByIDPath = "/wallets/%s"
	sendMoneyPath     = "/wallet"
)

func TestCreateWallet(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/api/v1",
	}
	e := httpexpect.Default(t, u.String())

	e.POST(createWalletPath).
		WithJSON(wallets.CreateWalletRequest{
			Currency: gofakeit.Currency().Short,
		}).
		Expect().
		Status(http.StatusBadRequest)

	e.POST(createWalletPath).
		WithJSON(wallets.CreateWalletRequest{
			Balance:  float64(gofakeit.Number(100, 10000)),
			Currency: gofakeit.Currency().Short,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object().ContainsKey("balance")
}

func TestGetWalletById(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/api/v1",
	}
	e := httpexpect.Default(t, u.String())

	result := e.GET(getAllWalletsPath).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	raw := result.Raw()

	data, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var response wallets.GetAllWalletsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	pathtest := fmt.Sprintf(getWalletByIDPath, response.Wallets[0].ID)
	e.GET(pathtest).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("balance")
}

func SendMoneyTest(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/api/v1",
	}
	e := httpexpect.Default(t, u.String())

	result := e.GET(getAllWalletsPath).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	raw := result.Raw()

	data, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var response wallets.GetAllWalletsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	testCases := []wallets.SendCoinRequest{
		{
			WalletID:       response.Wallets[0].ID,
			SenderWalletID: &response.Wallets[1].ID,
			Type:           wallets.OperationTypeTransfer,
			Amount:         100,
		},
		{
			WalletID:       response.Wallets[1].ID,
			SenderWalletID: nil,
			Type:           wallets.OperationTypeDeposit,
			Amount:         1,
		},
		{
			WalletID:       response.Wallets[0].ID,
			SenderWalletID: nil,
			Type:           wallets.OperationTypeWithdraw,
			Amount:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.Type), func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
				Path:   "/api/v1",
			}
			e := httpexpect.Default(t, u.String())

			e.POST(sendMoneyPath).
				WithJSON(wallets.SendCoinRequest{
					WalletID:       tc.WalletID,
					SenderWalletID: tc.SenderWalletID,
					Type:           tc.Type,
					Amount:         tc.Amount,
				}).
				Expect().
				Status(http.StatusOK).
				JSON().
				Object()

		})
	}
}
