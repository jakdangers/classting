package domain

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (int, error)
	FindUserByUserName(ctx context.Context, userName string) (*User, error)
}

// UserType 고랭에는 별도의 enum 타입이 없어서 string으로 정의
type UserType string

const (
	UserUseTypeAdmin   UserType = "ADMIN"
	UserUseTypeStudent UserType = "STUDENT"
)

type User struct {
	Base
	UserName string
	Password string
	UseType  UserType
}
