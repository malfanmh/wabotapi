package customecho

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
)

func ErrorHandler(c *echo.Echo) {
	c.HTTPErrorHandler = func(err error, c echo.Context) {
		resp := Response{}
		if !c.Response().Committed {
			var he *echo.HTTPError
			if errors.As(err, &he) {
				err = he
			}
			// Validation Error
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				err = errs
			}

			// gRPC Error
			if st, ok := status.FromError(err); ok {
				resp.Code = st.Code()
				resp.Message = st.Message()
			}
			_ = resp.JSON(c)
			return
		}
	}
}
