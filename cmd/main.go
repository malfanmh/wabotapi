package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/malfanmh/wabotapi"
	"github.com/malfanmh/wabotapi/pkg/customecho"
	"github.com/malfanmh/wabotapi/repository"
	"github.com/malfanmh/wabotapi/usecase"
	"net/http"
)

func main() {
	env := wabotapi.New()

	e := customecho.SetupEcho(true)

	db := wabotapi.OpenMysqlDB(env.MysqlDSN)
	repo := repository.NewMysql(db)
	wa := repository.NewWhatsappAPI(env.WABaseURL, env.WABusinessAccountID, env.WAAccessToken, &http.Client{})
	uc := usecase.New(repo, wa)
	handler := wabotapi.NewHandler(uc)

	// setup Router
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong !!")
	})
	e.GET("/webhook/:clientHash", handler.VerifyWebhook)
	e.POST("/webhook/:clientHash", handler.Webhook)
	//e.GET("/send", func(c echo.Context) error {
	//	resp, err := sendMessage(client, msgTemplate1)
	//	if err != nil {
	//		return c.String(http.StatusInternalServerError, err.Error())
	//	}
	//	return c.JSON(http.StatusOK, resp)
	//})

	fmt.Println(e.Start(":" + env.APPPort))
}
