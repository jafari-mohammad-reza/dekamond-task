package dto

import "dekamond-task/internal/models"

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
type UsersResponse struct {
	Users []*models.User `json:"users"`
}
