package subscription

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

func RegisterRoutes(e *gin.Engine, controller domain.SubscriptionController, cfg *config.Config) {
	api := e.Group("/subscriptions")
	{
		api.POST("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeStudent}), controller.CreateSubscription)
		api.GET("", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeStudent}), controller.ListSubscriptionSchools)
		api.GET("/news/:schoolID", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeStudent}), controller.ListSubscriptionSchoolNews)
		api.DELETE("/:schoolID", router.JWTMiddleware(cfg.Auth.Secret, []domain.UserType{domain.UserUseTypeStudent}), controller.DeleteSubscription)
	}
}

type subscriptionController struct {
	service domain.SubscriptionService
}

func NewSubscriptionController(service domain.SubscriptionService) *subscriptionController {
	return &subscriptionController{
		service: service,
	}
}

var _ domain.SubscriptionController = (*subscriptionController)(nil)

// CreateSubscription
// @Tags Subscription
// @Summary 구독 생성 [필수 구현] 권한 - 학생
// @Description 학교ID로 구독을 생성합니다.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateSubscriptionRequest body domain.CreateSubscriptionRequest true "구독 생성 요청"
// @Success 204
// @Router /subscriptions [post]
func (n subscriptionController) CreateSubscription(c *gin.Context) {
	var req domain.CreateSubscriptionRequest

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

	if err := n.service.CreateSubscription(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSubscriptionSchools
// @Summary 구독 중인 학교 목록 조회 [필수 구현] 권한 - 학생
// @Description 구독 학교 목록을 20개씩 조회합니다
// @Tags Subscription
// @Produce json
// @Security BearerAuth
// @Param cursor query int false "커서"
// @Success 200 {object} domain.ListSubscriptionSchoolsResponse "구독 학교 목록"
// @Router /subscriptions [get]
func (n subscriptionController) ListSubscriptionSchools(c *gin.Context) {
	var req domain.ListSubscriptionSchoolsRequest

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

	res, err := n.service.ListSubscriptionSchools(ctx, req)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.JSON(domain.ClasstingResponseFrom(http.StatusOK, res))
}

// DeleteSubscription
// @Summary 구독 취소 [필수 구현] 권한 - 학생
// @Description 학교 ID로 구독을 취소합니다.
// @Tags Subscription
// @Produce json
// @Param schoolID path int true "학교 ID"
// @Security BearerAuth
// @Success 204
// @Router /subscriptions/{schoolID} [delete]
func (n subscriptionController) DeleteSubscription(c *gin.Context) {
	var req domain.DeleteSubscriptionRequest

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

	if err := n.service.DeleteSubscription(ctx, req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSubscriptionSchoolNews
// @Summary 구독 중인 학교 페이지별 소식 조회 [필수 구현] 권한 - 학생
// @Description 구독 중인 학교 페이지별 소식을 20개씩 조회합니다
// @Tags Subscription
// @Produce json
// @Security BearerAuth
// @Param cursor query int false "커서"
// @Param schoolID path int true "학교 ID"
// @Success 200 {object} domain.ListSubscriptionSchoolNewsResponse "구독 중인 학교 페이지별 소식 조회"
// @Router /subscriptions/news/{schoolID} [get]
func (n subscriptionController) ListSubscriptionSchoolNews(c *gin.Context) {
	var req domain.ListSubscriptionSchoolNewsRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	if err := c.ShouldBindUri(&req); err != nil {
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

	res, err := n.service.ListSubscriptionSchoolNews(ctx, req)
	if err != nil {
		c.JSON(cerrors.ToSentinelAPIError(err))
		return
	}

	c.JSON(domain.ClasstingResponseFrom(http.StatusOK, res))
}
