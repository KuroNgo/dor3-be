package user_domain

import "context"

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

//go:generate mockery --name IProfileUseCase
type IProfileUseCase interface {
	GetProfileByID(c context.Context, userID string) (*Profile, error)
}
