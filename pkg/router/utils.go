package router

import (
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"github.com/gin-gonic/gin"
)

func GetUserInfoFromContext(c *gin.Context) (int, domain.UserType, error) {
	const op cerrors.Op = "router/GetUserIDFromContext"

	userID, ok := c.Get("userID")
	if !ok {
		return 0, "", cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return 0, "", cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	userType, ok := c.Get("userType")
	if !ok {
		return 0, "", cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	userTypeEnum, ok := userType.(domain.UserType)
	if !ok {
		return 0, "", cerrors.E(op, cerrors.Internal, "서버에 문제가 발생했습니다.")
	}

	return userIDInt, userTypeEnum, nil
}
