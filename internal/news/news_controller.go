package news

import (
	"classting/config"
	"classting/domain"
	"classting/pkg/cerrors"
	"classting/pkg/router"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes(e *gin.Engine, controller domain.NewsController, cfg *config.Config) {
	api := e.Group("/news")
	{
		api.POST("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin}), controller.CreateNews)
		api.GET("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin}), controller.ListNews)
		api.PUT("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin}), controller.UpdateNews)
		api.DELETE("/:newsID", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeAdmin}), controller.DeleteNews)
	}
}

type newsController struct {
	service domain.NewsService
}

func NewNewsController(service domain.NewsService) *newsController {
	return &newsController{
		service: service,
	}
}

var _ domain.NewsController = (*newsController)(nil)

// CreateNews
// @Tags News
// @Summary 소식 생성 [필수 구현] 권한 - 관리자
// @Description 학교ID, 제목으로 소식을 생성 (자신의 학교에만 소식을 생성할 수 있음, 학교를 여러개 소유 할 수 있으므로 학교 ID 필요)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateNewsRequest body domain.CreateNewsRequest true "소식 생성 요청"
// @Success 204
// @Router /news [post]
func (n newsController) CreateNews(c *gin.Context) {
	var req domain.CreateNewsRequest

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

	if err := n.service.CreateNews(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ListNews
// @Summary 학교 소식 목록 조회 [필수 구현] 권한 - 관리자
// @Description 학교 소식 목록을 20개씩 조회합니다
// @Tags News
// @Produce json
// @Security BearerAuth
// @Param cursor query int false "커서"
// @Success 200 {object} domain.ListNewsResponse "학교 목록"
// @Router /news [get]
func (n newsController) ListNews(c *gin.Context) {
	var req domain.ListNewsRequest

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

	res, err := n.service.ListNews(ctx, req)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.JSON(domain.ClasstingResponseFrom(http.StatusOK, res))
}

// UpdateNews
// @Summary 소식 수정 [필수 구현] 권한 - 관리자
// @Description 학교 소식을 수정 (단 자신의 학교 소식만 수정 가능)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param UpdateNewsRequest body domain.UpdateNewsRequest true "소식 수정 요청"
// @Success 204
// @Router /news [put]
func (n newsController) UpdateNews(c *gin.Context) {
	var req domain.UpdateNewsRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	userID, err := router.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}
	req.UserID = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := n.service.UpdateNews(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteNews
// @Summary 소식 삭제 [필수 구현] 권한 - 관리자
// @Description 소식 ID로 상품을 삭제합니다. (단 자신의 학교 소식만 삭제 가능)
// @Tags News
// @Produce json
// @Param newsID path int true "소식 ID"
// @Security BearerAuth
// @Success 204
// @Router /news/{newsID} [delete]
func (n newsController) DeleteNews(c *gin.Context) {
	var req domain.DeleteNewsRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	userID, err := router.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}
	req.UserID = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := n.service.DeleteNews(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}
