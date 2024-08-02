package wabotapi

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type UseCase interface {
	VerifyWebhook(ctx context.Context, hash, token string) (bool, error)
	Webhook(ctx context.Context, hash string, payload []byte) error
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
		hash := c.Param("clientHash")
		ok, err := h.uc.VerifyWebhook(c.Request().Context(), hash, token)
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
	if err = h.uc.Webhook(c.Request().Context(), c.Param("clientHash"), payload); err != nil {
		fmt.Println("webhook error", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func sendMessage(client *resty.Client, msg string) (string, error) {
	//msgTemplate1 := `{"messaging_product":"whatsapp","to":"6281286312989","type":"template","template":{"name":"qna_template_test_1","language":{"code":"en"}}}`
	//msgHello := `{ "messaging_product": "whatsapp", "to": "6281286312989", "type": "template", "template": { "name": "hello_world", "language": { "code": "en_US" } } }`
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetBody(msg).
		Post("385924484596973/messages")
	if err != nil {
		return "", err
	}
	b := resp.Body()
	fmt.Println(string(b), resp.StatusCode())
	return string(b), nil
}
