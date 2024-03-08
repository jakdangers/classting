package domain

import (
	"context"
	"github.com/gin-gonic/gin"
)

type SchoolRepository interface {
	CreateSchool(ctx context.Context, school School) (int, error)
	FindSchoolByNameAndRegion(ctx context.Context, params FindSchoolByNameAndRegionParams) (*School, error)
}

type SchoolService interface {
	CreateSchool(ctx context.Context, req CreateSchoolRequest) error
}

type SchoolController interface {
	CreateSchool(c *gin.Context)
}

type School struct {
	Base
	UserID int
	Name   string
	Region string
}

type FindSchoolByNameAndRegionParams struct {
	Name   string
	Region string
}
