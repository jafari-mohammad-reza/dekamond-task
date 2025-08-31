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
	GetUsers(page, limit int) ([]*models.User, error)
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
		return "", errors.New("invalid credentials")
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	token, err := u.tokenService.GenerateToken(map[string]any{
		"mobile": mobile,
	})
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

func (u *UserService) CreateUser(mobile string) error {
	ctx := context.Background()

	user, err := u.repo.GetUser(ctx, mobile)
	if err != nil && err.Error() != "user not found" {
		return errors.New("failed to check existing user")
	}

	if user != nil {
		return errors.New("invalid credentials")
	}

	if err := u.repo.CreateUser(ctx, mobile); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (u *UserService) GetUser(mobile string) (*models.User, error) {
	user, err := u.repo.GetUser(context.Background(), mobile)
	if err != nil {
		return nil, errors.New("failed to get user")
	}
	return user, nil
}

func (u *UserService) GetUsers(page, limit int) ([]*models.User, error) {
	users, err := u.repo.GetUsers(context.Background(), page, limit)
	if err != nil {
		return nil, errors.New("failed to get users")
	}
	return users, nil
}

func NewUserService(cfg *config.Config) (IUserService, error) {
	repo, err := db.NewDB(cfg)
	if err != nil {
		return nil, errors.New("failed to create db connection")
	}
	return &UserService{
		cfg:          cfg,
		repo:         repo,
		tokenService: NewTokenService(cfg),
	}, nil
}
