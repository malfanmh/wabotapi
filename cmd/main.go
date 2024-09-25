package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/malfanmh/wabotapi"
	"github.com/malfanmh/wabotapi/pkg/customecho"
	"github.com/malfanmh/wabotapi/repository"
	"github.com/malfanmh/wabotapi/usecase"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	env := wabotapi.New()

	db := wabotapi.OpenMysqlDB(env.MysqlDSN)
	repo := repository.NewMysql(db)
	wa := repository.NewWhatsappAPI(env.WABaseURL, env.WABusinessAccountID, env.WAAccessToken, &http.Client{})
	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	payment := repository.NewFinpay(httpClient, env.FinpayBaseURL)
	uc := usecase.New(repo, wa, payment, env.WASecret)
	handler := wabotapi.NewHandler(uc)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	// setup Router
	e := customecho.SetupEcho(true)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong !!")
	})
	e.GET("/webhook", handler.VerifyWebhook)
	e.POST("/webhook", handler.Webhook)
	e.POST("/finpay/callback", handler.FinpayCallback)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", env.APPPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start http server\n")
			cancelFunc()
		}
	}()

	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(jakarta)
	log.Printf("starting ExpiryLink scheduler every minutes \n")
	s.Every(time.Minute).Do(func() {
		if err := handler.ExpiryLink(context.Background()); err != nil {
			log.Printf("scheduler ExpiryLink err %v\n", err)
		}
	})
	go func() {
		s.StartBlocking()
		cancelFunc()
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("shutting down http")
	s.Stop()
	log.Printf("shutting down cron")

}
