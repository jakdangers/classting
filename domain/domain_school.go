package domain

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SchoolRepository interface {
	CreateSchool(ctx context.Context, school School) (int, error)
	ListSchools(ctx context.Context, params ListSchoolsParams) ([]School, error)
	FindSchoolByNameAndRegion(ctx context.Context, params FindSchoolByNameAndRegionParams) (*School, error)
	FindSchoolByID(ctx context.Context, schoolID int) (*School, error)
}

type SchoolService interface {
	CreateSchool(ctx context.Context, req CreateSchoolRequest) error
	ListSchools(ctx context.Context, req ListSchoolsRequest) (ListSchoolsResponse, error)
}

type SchoolController interface {
	CreateSchool(c *gin.Context)
	ListSchools(c *gin.Context)
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

type ListSchoolsParams struct {
	UserID *int
	Cursor *int
}

func (lp ListSchoolsParams) AndUserID() string {
	if lp.UserID == nil {
		return ""
	}

	return fmt.Sprintf("AND user_id = %d", *lp.UserID)
}

func (lp ListSchoolsParams) AfterCursor() string {
	if lp.Cursor == nil {
		return ""
	}

	return fmt.Sprintf("AND id > %d", *lp.Cursor)
}
