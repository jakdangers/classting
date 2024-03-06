package cerrors

import (
	"errors"
	"net/http"
)

type SentinelAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewSentinelAPIError(statusCode int, message string) (int, SentinelAPIError) {
	return statusCode, SentinelAPIError{
		Code:    statusCode,
		Message: message,
	}
}

func (e SentinelAPIError) Error() string {
	return e.Message
}

func ToSentinelAPIError(err error) (int, SentinelAPIError) {
	var cErr *Error

	if errors.As(err, &cErr) {
		switch cErr.Kind {
		case Other, Internal:
			return NewSentinelAPIError(http.StatusInternalServerError, cErr.ServiceMessage)
		case Invalid, Permission:
			return NewSentinelAPIError(http.StatusBadRequest, cErr.ServiceMessage)
		case Auth:
			return NewSentinelAPIError(http.StatusUnauthorized, cErr.ServiceMessage)
		case Exist:
			return NewSentinelAPIError(http.StatusConflict, cErr.ServiceMessage)
		case NotExist:
			return NewSentinelAPIError(http.StatusNotFound, cErr.ServiceMessage)
		default:
			return NewSentinelAPIError(http.StatusInternalServerError, cErr.ServiceMessage)
		}
	}

	return http.StatusBadRequest, SentinelAPIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}
