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
