package school

import (
	"classting/config"
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"classting/pkg/router"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes(e *gin.Engine, controller domain.SchoolController, cfg *config.Config) {
	api := e.Group("/schools")
	{
		api.POST("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin}), controller.CreateSchool)
	}
}

type schoolController struct {
	service domain.SchoolService
}

func NewSchoolController(service domain.SchoolService, cfg *config.Config) *schoolController {
	return &schoolController{
		service: service,
	}
}

var _ domain.SchoolController = (*schoolController)(nil)

// CreateSchool
// @Tags School
// @Summary 학교 생성
// @Description 지역, 학교명으로 학교를 생성합니다. (지역, 학교명이 중복되지 않아야 합니다.)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateSchoolRequest body domain.CreateSchoolRequest true "회원가입 요청"
// @Success 204
// @Router /schools [post]
func (u schoolController) CreateSchool(c *gin.Context) {
	var req domain.CreateSchoolRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	userID, err := router.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}
	req.UserID = userID

	if err := req.Validate(); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := u.service.CreateSchool(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
