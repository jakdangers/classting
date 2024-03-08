package domain

import "classting/pkg/cerrors"

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
