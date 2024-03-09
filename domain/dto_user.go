package domain

import (
	cerrors "classting/pkg/cerrors"
)

type CreateUserRequest struct {
	UserName string   `json:"userName" validate:"required" example:"classting_admin"`
	Password string   `json:"password" validate:"required" example:"1234"`
	UserType UserType `json:"userType" validate:"required" enum:"ADMIN,STUDENT" example:"ADMIN"`
}

func (ur CreateUserRequest) Validate() error {
	const op cerrors.Op = "user/controller/valid"

	if ur.UserName == "" {
		return cerrors.E(op, cerrors.Invalid, "잘못된 아이디입니다.")
	}

	if ur.Password == "" {
		return cerrors.E(op, cerrors.Invalid, "잘못된 비밀번호입니다.")
	}

	return nil
}

type LoginUserRequest struct {
	UserName string `json:"userName" validate:"required" example:"classting_admin"`
	Password string `json:"password" validate:"required" example:"1234"`
}

func (ur LoginUserRequest) Validate() error {
	const op cerrors.Op = "user/controller/valid"

	if ur.UserName == "" {
		return cerrors.E(op, cerrors.Invalid, "아이디 또는 비밀번호를 확인해주세요.")
	}

	if ur.Password == "" {
		return cerrors.E(op, cerrors.Invalid, "아이디 또는 비밀번호를 확인해주세요.")
	}

	return nil
}

type LoginUserResponse struct {
	AccessToken string `json:"accessToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDg4ODgxOTIsInVzZXJJRCI6MX0.WVQGpeNbCpWSKuvYO7rFv6HoXaEA4_VQZSl7oMhmROk"`
	ExpiresIn   int64  `json:"expiresIn" validate:"required" example:"1708888192"`
}
