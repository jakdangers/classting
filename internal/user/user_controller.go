package user

import (
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes(e *gin.Engine, controller domain.UserController) {
	api := e.Group("/users")
	{
		api.POST("", controller.CreateUser)
	}
}

type userController struct {
	service domain.UserService
}

func NewUserController(service domain.UserService) *userController {
	return &userController{
		service: service,
	}
}

var _ domain.UserController = (*userController)(nil)

// CreateUser
// @Tags User
// @Summary 회원가입
// @Description 관리자, 학생의 역할로 회원가입 요청 (관리자의 경우 UserType = ADMIN, 학생의 경우 UserType = STUDENT)
// @Accept json
// @Produce json
// @Param CreateUserRequest body domain.CreateUserRequest true "회원가입 요청"
// @Success 204
// @Router /users [post]
func (u userController) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := u.service.CreateUser(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
