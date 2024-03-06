package domain

import (
	"net/http"
	"time"
)

type ClasstingResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type BaseDTO struct {
	ID         int       `json:"id" validate:"required" example:"1"`
	CreateDate time.Time `json:"createDate" validate:"required" example:"2024-02-28T15:04:05Z"`
	UpdateDate time.Time `json:"updateDate" validate:"required" example:"2024-02-28T15:04:05Z"`
}

func ClasstingResponseFrom(code int, data any) (int, ClasstingResponse) {
	return code, ClasstingResponse{
		Code:    code,
		Message: http.StatusText(code),
		Data:    data,
	}
}
