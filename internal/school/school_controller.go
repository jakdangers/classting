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
		api.GET("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin, domain.UserUseTypeStudent}), controller.ListSchools)
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
// @Tags Schools
// @Summary 학교 생성 [필수 구현] 권한 - 관리자
// @Description 지역, 학교명으로 학교를 생성합니다. (지역, 학교명이 중복되지 않아야 합니다.)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateSchoolRequest body domain.CreateSchoolRequest true "학교 생성 요청"
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

// ListSchools
// @Summary 학교 목록 조회 [테스트 추가 API] 권한 - 관리자, 학생
// @Description 학교 목록을 조회합니다  userID를 쿼리스트링에 입력시 해당 유저의 학교 목록을 조회합니다. 디폴트 20개씩 조회
// @Tags Schools
// @Produce json
// @Security BearerAuth
// @Param cursor query int false "커서"
// @Param userID query int false "유저 아이디"
// @Success 200 {object} domain.ListSchoolsResponse "학교 목록"
// @Router /schools [get]
func (u schoolController) ListSchools(c *gin.Context) {
	var req domain.ListSchoolsRequest

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

	res, err := u.service.ListSchools(ctx, req)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.JSON(domain.ClasstingResponseFrom(http.StatusOK, res))
}
