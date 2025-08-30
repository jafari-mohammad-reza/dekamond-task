package service

import (
	"context"
	"dekamond-task/internal/config"
	"dekamond-task/internal/db"
	"dekamond-task/internal/models"
	"errors"
	"time"
)

type IUserService interface {
	CreateUser(mobile string) error
	Login(mobile string) (string, error)
	GetUser(mobile string) (*models.User, error)
	GetUsers(page , limit int) ([]*models.User, error)
}

type UserService struct {
	cfg          *config.Config
	repo         *db.DB
	tokenService *TokenService
}

func (u *UserService) Login(mobile string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := u.repo.GetUser(ctx, mobile)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	token, err := u.tokenService.GenerateToken(map[string]any{
		"mobile": mobile,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserService) CreateUser(mobile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	user, err := u.repo.GetUser(ctx, mobile)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("invalid credentials")
	}
	if err := u.CreateUser(mobile); err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUser(mobile string) (*models.User, error) {
	user, err := u.repo.GetUser(context.Background(), mobile)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUsers(page, limit int) ([]*models.User, error) {
	users, err := u.repo.GetUsers(context.Background(), page, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func NewUserService(cfg *config.Config) (IUserService, error) {
	repo, err := db.NewDB(cfg)
	if err != nil {
		return nil, err
	}
	return &UserService{
		cfg:          cfg,
		repo:         repo,
		tokenService: NewTokenService(cfg),
	}, nil
}
