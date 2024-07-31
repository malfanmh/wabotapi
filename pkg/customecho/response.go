package customecho

import (
	"github.com/RekeningkuDev/portfolio-analyst/api/config"
	"github.com/RekeningkuDev/portfolio-analyst/api/model"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
)

type Response struct {
	Code       codes.Code        `json:"code"`
	Message    string            `json:"message"`
	Data       any               `json:"data"`
	Pagination *model.Pagination `json:"meta,omitempty"`
}

type responseOK struct {
	Message    string            `json:"message"`
	Data       any               `json:"data"`
	Pagination *model.Pagination `json:"meta,omitempty"`
}

func (r *Response) JSON(c echo.Context) error {
	status, ok := config.HTTPStatus[r.Code]
	if !ok {
		status = runtime.HTTPStatusFromCode(r.Code)
	}
	if r.Code == config.CodeSuccess {
		return c.JSON(status, responseOK{
			Message:    r.Message,
			Data:       r.Data,
			Pagination: r.Pagination,
		})
	}
	return c.JSON(status, r)
}
