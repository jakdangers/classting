package domain

import (
	"classting/pkg/cerrors"
)

type NewsDTO struct {
	BaseDTO
	SchoolID int    `json:"schoolID"`
	Title    string `json:"title"`
}

type CreateNewsRequest struct {
	UserID   int    `swaggerignore:"true"`
	SchoolID int    `json:"schoolID"`
	Title    string `json:"title"`
}

func (req CreateNewsRequest) Validate() error {
	const op cerrors.Op = "domain/CreateNewsRequest.Validate"

	if req.SchoolID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "학교 ID를 확인해주세요.")
	}

	if req.Title == "" {
		return cerrors.E(op, cerrors.Invalid, "제목을 확인해주세요.")
	}

	return nil
}

type ListNewsRequest struct {
	UserID   int  `swaggerignore:"true"`
	SchoolID int  `form:"schoolID" validate:"required" example:"1"`
	Cursor   *int `form:"cursor" validate:"optional" example:"1"`
}

func (req ListNewsRequest) Validate() error {
	const op cerrors.Op = "domain/ListNewsRequest.Validate"

	if req.SchoolID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "학교 ID를 확인해주세요.")
	}

	if req.Cursor != nil && *req.Cursor <= 0 {
		return cerrors.E(op, cerrors.Invalid, "커서를 확인해주세요.")
	}

	return nil
}

type ListNewsResponse struct {
	News   []NewsDTO `json:"news"`
	Cursor *int      `json:"cursor"`
}

type UpdateNewsRequest struct {
	UserID int    `swaggerignore:"true"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

func (req UpdateNewsRequest) Validate() error {
	const op cerrors.Op = "domain/UpdateNewsRequest.Validate"

	if req.ID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "ID를 확인해주세요.")
	}

	if req.Title == "" {
		return cerrors.E(op, cerrors.Invalid, "제목을 확인해주세요.")
	}

	return nil
}

type DeleteNewsRequest struct {
	UserID int `swaggerignore:"true"`
	ID     int `uri:"newsID"`
}

func (req DeleteNewsRequest) Validate() error {
	const op cerrors.Op = "domain/DeleteNewsRequest.Validate"

	if req.ID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "ID를 확인해주세요.")
	}

	return nil
}

func NewsDTOFrom(news News) NewsDTO {
	return NewsDTO{
		BaseDTO: BaseDTO{
			ID:         news.ID,
			CreateDate: news.CreateDate,
			UpdateDate: news.UpdateDate,
		},
		SchoolID: news.SchoolID,
		Title:    news.Title,
	}
}
