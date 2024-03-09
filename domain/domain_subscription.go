package domain

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription Subscription) (int, error)
	ListSubscriptionSchools(ctx context.Context, params ListSubscriptionSchoolsParams) ([]SubscriptionSchool, error)
	DeleteSubscription(ctx context.Context, subscriptionID int) error
	FindSubscriptionByUserIDAndSchoolID(ctx context.Context, params FindSubscriptionByUserIDAndSchoolIDParams) (*Subscription, error)
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) error
	ListSubscriptionSchools(ctx context.Context, req ListSubscriptionSchoolsRequest) (ListSubscriptionSchoolsResponse, error)
	ListSubscriptionSchoolNews(ctx context.Context, req ListSubscriptionSchoolNewsRequest) (ListSubscriptionSchoolNewsResponse, error)
	DeleteSubscription(ctx context.Context, req DeleteSubscriptionRequest) error
}

type SubscriptionController interface {
	CreateSubscription(c *gin.Context)
	ListSubscriptionSchools(c *gin.Context)
	ListSubscriptionSchoolNews(c *gin.Context)
	DeleteSubscription(c *gin.Context)
}

type Subscription struct {
	Base
	UserID   int
	SchoolID int
}

type SubscriptionSchool struct {
	Base
	SchoolID int
	Name     string
	Region   string
}

type FindSubscriptionByUserIDAndSchoolIDParams struct {
	UserID   int
	SchoolID int
}

type ListSubscriptionSchoolsParams struct {
	UserID int
	Cursor *int
}

func (lp ListSubscriptionSchoolsParams) AfterCursor() string {
	if lp.Cursor == nil {
		return ""
	}

	return fmt.Sprintf("AND subscriptions.id < %d", *lp.Cursor)
}
