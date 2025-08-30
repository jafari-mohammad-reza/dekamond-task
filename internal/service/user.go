package service

import (
	"dekamond-task/internal/config"
	"dekamond-task/internal/db"
	"dekamond-task/internal/models"
)

type IUserService interface {
	CreateUser(mobile string) error
	Login(mobile string, otp string) (string, error)
	GetUser(mobile string) (*models.User, error)
	GetUsers() ([]*models.User, error)
}

type UserService struct {
	cfg  *config.Config
	repo *db.DB
}

func (u *UserService) Login(mobile string, otp string) (string, error) {
	panic("unimplemented")
}

func (u *UserService) CreateUser(mobile string) error {
	panic("unimplemented")
}

func (u *UserService) GetUser(mobile string) (*models.User, error) {
	panic("unimplemented")
}

func (u *UserService) GetUsers() ([]*models.User, error) {
	panic("unimplemented")
}

func NewUserService(cfg *config.Config) (IUserService, error) {
	repo, err := db.NewDB(cfg)
	if err != nil {
		return nil, err
	}
	return &UserService{
		cfg:  cfg,
		repo: repo,
	}, nil
}
