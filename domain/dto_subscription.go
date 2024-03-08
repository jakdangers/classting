package domain

type SubscriptionDTO struct {
	BaseDTO
	UserID   int    `json:"userID" validate:"required" example:"1"`
	SchoolID int    `json:"schoolID" validate:"required" example:"1"`
	Name     string `json:"name" validate:"required" example:"클래스팅"`
	Region   string `json:"region" validate:"required" example:"서울"`
}
