package wabotapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/malfanmh/wabotapi/model"
	"io"
	"net/http"
)

type UseCase interface {
	VerifyWebhook(ctx context.Context, token string) (bool, error)
	Webhook(ctx context.Context, payload []byte) error
	PaymentCallback(ctx context.Context, callback model.FinpayCallback) error
	ExpiryLink(ctx context.Context) error
}

type handler struct {
	uc UseCase
}

func NewHandler(uc UseCase) *handler {
	return &handler{
		uc: uc,
	}
}

func (h *handler) VerifyWebhook(c echo.Context) error {
	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challenge := c.QueryParam("hub.challenge")
	if mode == "subscribe" {
		ok, err := h.uc.VerifyWebhook(c.Request().Context(), token)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return c.String(http.StatusOK, challenge)
	}

	return c.JSON(http.StatusUnauthorized, nil)
}

func (h *handler) Webhook(c echo.Context) error {
	//var payload map[string]interface{}
	//err := c.Bind(&payload)
	//if err != nil {
	//	return c.String(http.StatusBadRequest, err.Error())
	//}

	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error reading body")
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Println("payload", string(payload))
	if err = h.uc.Webhook(c.Request().Context(), payload); err != nil {
		fmt.Println("webhook error", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *handler) FinpayCallback(c echo.Context) error {
	//Merchant ID : FIES994
	//Merchant Key : WqSG4ViP0v6jQA9Zp3b8ofr2LkdIYmBx
	var payload model.FinpayCallback
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	fmt.Println("FinpayCallback payload", payload)
	err = h.uc.PaymentCallback(c.Request().Context(), payload)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *handler) ExpiryLink(ctx context.Context) error {
	return h.uc.ExpiryLink(ctx)
}
