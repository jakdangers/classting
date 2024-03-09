package domain

import (
	"classting/pkg/cerrors"
)

type SubscriptionSchoolDTO struct {
	BaseDTO
	SchoolID int    `json:"schoolID" validate:"required" example:"1"`
	Name     string `json:"name" validate:"required" example:"클래스팅"`
	Region   string `json:"region" validate:"required" example:"서울"`
}

type SubscriptionSchoolNewsDTO struct {
	BaseDTO
	SchoolID int    `json:"schoolID" validate:"required" example:"1"`
	Title    string `json:"title" validate:"required" example:"클래스팅 새소식"`
}

type CreateSubscriptionRequest struct {
	UserID   int `swaggerignore:"true"`
	SchoolID int `json:"schoolID" validate:"required" example:"1"`
}

func (req CreateSubscriptionRequest) Validate() error {
	const op cerrors.Op = "domain/CreateSubscriptionRequest.Validate"

	if req.SchoolID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "학교 ID를 확인해주세요.")
	}

	return nil
}

type ListSubscriptionSchoolsRequest struct {
	UserID int  `swaggerignore:"true"`
	Cursor *int `form:"cursor" example:"1"`
}

func (req ListSubscriptionSchoolsRequest) Validate() error {
	const op cerrors.Op = "domain/ListSubscriptionSchoolsRequest.Validate"

	if req.Cursor != nil && *req.Cursor <= 0 {
		return cerrors.E(op, cerrors.Invalid, "커서를 확인해주세요.")
	}

	return nil
}

type ListSubscriptionSchoolsResponse struct {
	SubscriptionSchools []SubscriptionSchoolDTO `json:"subscriptionSchools"`
	Cursor              *int                    `json:"cursor"`
}

type DeleteSubscriptionRequest struct {
	UserID   int `swaggerignore:"true"`
	SchoolID int `uri:"schoolID" validate:"required" example:"1"`
}

func (req DeleteSubscriptionRequest) Validate() error {
	const op cerrors.Op = "domain/DeleteSubscriptionRequest.Validate"

	if req.SchoolID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "학교 ID를 확인해주세요.")
	}

	return nil
}

type ListSubscriptionSchoolNewsRequest struct {
	UserID   int  `swaggerignore:"true"`
	SchoolID int  `uri:"schoolID" validate:"required" example:"1"`
	Cursor   *int `form:"cursor"`
}

func (req ListSubscriptionSchoolNewsRequest) Validate() error {
	const op cerrors.Op = "domain/ListSubscriptionSchoolNewsRequest.Validate"

	if req.SchoolID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "학교 ID를 확인해주세요.")
	}

	if req.Cursor != nil && *req.Cursor <= 0 {
		return cerrors.E(op, cerrors.Invalid, "커서를 확인해주세요.")
	}

	return nil
}

type ListSubscriptionSchoolNewsResponse struct {
	SubscriptionSchoolNews []SubscriptionSchoolNewsDTO `json:"subscriptionSchoolNews"`
	Cursor                 *int                        `json:"cursor"`
}

func SubscriptionSchoolDTOFrom(subscriptionSchool SubscriptionSchool) SubscriptionSchoolDTO {
	return SubscriptionSchoolDTO{
		BaseDTO: BaseDTO{
			ID:         subscriptionSchool.ID,
			CreateDate: subscriptionSchool.CreateDate,
			UpdateDate: subscriptionSchool.UpdateDate,
		},
		SchoolID: subscriptionSchool.SchoolID,
		Name:     subscriptionSchool.Name,
		Region:   subscriptionSchool.Region,
	}
}

func SubscriptionSchoolNewsDTOFrom(news News) SubscriptionSchoolNewsDTO {
	return SubscriptionSchoolNewsDTO{
		BaseDTO: BaseDTO{
			ID:         news.ID,
			CreateDate: news.CreateDate,
			UpdateDate: news.UpdateDate,
		},
		SchoolID: news.SchoolID,
		Title:    news.Title,
	}
}
