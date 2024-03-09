package router

import (
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

func parseJWTToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func extractInfoFromToken(token *jwt.Token) (int, domain.UserType, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("Invalid token claims")
	}

	userIDClaim, ok := claims["userID"]
	if !ok {
		return 0, "", fmt.Errorf("UserID claim not found")
	}

	userTypeClaim, ok := claims["userType"]
	if !ok {
		return 0, "", fmt.Errorf("Type claim not found")
	}

	userID, ok := userIDClaim.(float64) // 여기서 적절한 타입으로 형변환 필요
	if !ok {
		return 0, "", fmt.Errorf("Invalid UserID type")
	}

	userType, ok := userTypeClaim.(string) // 여기서 적절한 타입으로 형변환 필요
	if !ok {
		return 0, "", fmt.Errorf("Invalid Type type")
	}

	return int(userID), domain.UserType(userType), nil
}

func JWTMiddleware(secret string, userTypes []domain.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(cerrors.NewSentinelAPIError(http.StatusUnauthorized, "로그인 후 이용해주세요"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(cerrors.NewSentinelAPIError(http.StatusUnauthorized, "로그인 후 이용해주세요"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := parseJWTToken(tokenString, secret)
		if err != nil {
			c.JSON(cerrors.NewSentinelAPIError(http.StatusUnauthorized, "올바르지 않은 토큰입니다"))
			c.Abort()
			return
		}

		userID, usertype, err := extractInfoFromToken(token)
		if err != nil {
			c.JSON(cerrors.NewSentinelAPIError(http.StatusUnauthorized, "올바르지 않은 토큰입니다"))
			c.Abort()
			return
		}

		if !lo.Contains(userTypes, usertype) {
			c.JSON(cerrors.NewSentinelAPIError(http.StatusUnauthorized, "권한이 없습니다."))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("userType", usertype)
		c.Set("tokenString", tokenString)

		c.Next()
	}
}
