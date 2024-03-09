package domain

import "classting/pkg/cerrors"

type SchoolDTO struct {
	ID     int    `json:"id"`
	UserID int    `json:"userID"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type CreateSchoolRequest struct {
	UserID int    `swaggerignore:"true"`
	Name   string `json:"name" validate:"required" example:"클래스팅"`
	Region string `json:"region" validate:"required" example:"서울"`
}

func (req CreateSchoolRequest) Validate() error {
	var op cerrors.Op = "domain/CreateSchoolRequest.Validate"

	if req.Name == "" {
		return cerrors.E(op, cerrors.Invalid, "학교명을 확인해주세요.")
	}

	if req.Region == "" {
		return cerrors.E(op, cerrors.Invalid, "지역을 확인해주세요.")
	}

	return nil
}

type ListSchoolsRequest struct {
	UserID *int `form:"userID" example:"1"`
	Cursor *int `form:"cursor"`
}

func (req ListSchoolsRequest) Validate() error {
	var op cerrors.Op = "domain/ListSchoolsRequest.Validate"

	if req.UserID != nil && *req.UserID <= 0 {
		return cerrors.E(op, cerrors.Invalid, "사용자 ID를 확인해주세요.")
	}

	if req.Cursor != nil && *req.Cursor <= 0 {
		return cerrors.E(op, cerrors.Invalid, "커서를 확인해주세요.")
	}

	return nil
}

type ListSchoolsResponse struct {
	Schools []SchoolDTO `json:"schools"`
	Cursor  *int        `json:"cursor"`
}

func SchoolDTOFrom(school School) SchoolDTO {
	return SchoolDTO{
		ID:     school.ID,
		UserID: school.UserID,
		Name:   school.Name,
		Region: school.Region,
	}
}
