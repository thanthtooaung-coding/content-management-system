package types

import "github.com/google/uuid"

type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
	}
}
