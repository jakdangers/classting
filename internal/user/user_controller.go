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
		api.POST("/login", controller.LoginUser)
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
// @Summary 회원가입 [테스트 추가 API]
// @Description 관리자, 학생의 역할로 회원가입 요청 (관리자의 경우 Type = ADMIN, 학생의 경우 Type = STUDENT)
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

// LoginUser
// @Tags User
// @Summary 로그인 [테스트 추가 API]
// @Description 관리자 계정 classting_admin_1, classting_admin_2, classting_admin_3, empty_classting_admin
// @Description 학생 계정 classting_student_1, empty_classting_student
// @Description 비밀번호는 모두 1234 입니다.
// @Accept json
// @Produce json
// @Param LoginUserRequest body domain.LoginUserRequest true "로그인 요청"
// @Success 200 {object} domain.LoginUserResponse
// @Router /users/login [post]
func (u userController) LoginUser(c *gin.Context) {
	var req domain.LoginUserRequest

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

	res, err := u.service.LoginUser(ctx, req)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.JSON(domain.ClasstingResponseFrom(http.StatusOK, res))
}
