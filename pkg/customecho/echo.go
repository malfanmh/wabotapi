package customecho

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func SetupEcho(debugMode bool) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{Validator: validator.New()}
	e.Debug = debugMode
	ErrorHandler(e)
	return e
}
