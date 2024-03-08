package router

import (
	cerrors "classting/pkg/cerrors"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (int, error) {
	const op cerrors.Op = "router/GetUserIDFromContext"

	userID, ok := c.Get("userID")
	if !ok {
		return 0, cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return 0, cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	return userIDInt, nil
}
