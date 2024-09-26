package customecho

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
)

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	Total       int `json:"total"`
	Size        int `json:"size"`
	From        int `json:"from"`
	To          int `json:"to"`
}

type Response struct {
	Code       codes.Code  `json:"code"`
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"meta,omitempty"`
}

type responseOK struct {
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"meta,omitempty"`
}

func (r *Response) JSON(c echo.Context) error {
	status := runtime.HTTPStatusFromCode(r.Code)
	if r.Code == 200 {
		return c.JSON(status, responseOK{
			Message:    r.Message,
			Data:       r.Data,
			Pagination: r.Pagination,
		})
	}
	return c.JSON(status, r)
}
