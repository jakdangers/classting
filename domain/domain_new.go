package domain

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

type NewsRepository interface {
	CreateNews(ctx context.Context, news News) (int, error)
	ListNews(ctx context.Context, params ListNewsParams) ([]News, error)
	UpdateNews(ctx context.Context, news News) error
	FindNewsByID(ctx context.Context, newsID int) (*News, error)
	DeleteNews(ctx context.Context, newsID int) error
}

type NewsService interface {
	CreateNews(ctx context.Context, req CreateNewsRequest) error
	ListNews(ctx context.Context, req ListNewsRequest) (ListNewsResponse, error)
	UpdateNews(ctx context.Context, req UpdateNewsRequest) error
	DeleteNews(ctx context.Context, req DeleteNewsRequest) error
}

type NewsController interface {
	CreateNews(c *gin.Context)
	ListNews(c *gin.Context)
	UpdateNews(c *gin.Context)
	DeleteNews(c *gin.Context)
}

type News struct {
	Base
	SchoolID int
	UserID   int
	Title    string
}

type ListNewsParams struct {
	UserID int
	Cursor *int
}

func (lp ListNewsParams) AfterCursor() string {
	if lp.Cursor == nil {
		return ""
	}

	return fmt.Sprintf("AND id > %d", *lp.Cursor)
}
