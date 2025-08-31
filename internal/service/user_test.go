package service

import (
	"context"
	"dekamond-task/internal/config"
	"dekamond-task/internal/models"
	"errors"
	"testing"
	"time"
)

type mockRepo struct {
	users map[string]*models.User
}

func (m *mockRepo) CreateUser(ctx context.Context, mobile string) error {
	if _, ok := m.users[mobile]; ok {
		return errors.New("already exists")
	}
	m.users[mobile] = &models.User{ID: 1, Mobile: mobile, CreatedAt: time.Now()}
	return nil
}

func (m *mockRepo) GetUser(ctx context.Context, mobile string) (*models.User, error) {
	if u, ok := m.users[mobile]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockRepo) GetUsers(ctx context.Context, page, limit int) ([]*models.User, error) {
	var res []*models.User
	for _, u := range m.users {
		res = append(res, u)
	}
	return res, nil
}

func (m *mockRepo) InitTables() error {
	return nil
}
func (m *mockRepo) Close() error {
	return nil
}

func TestUserService(t *testing.T) {
	cfg := &config.Config{Token: struct{ Secret string }{Secret: "supersecretkey"}}
	mock := &mockRepo{users: make(map[string]*models.User)}
	us := &UserService{cfg: cfg, repo: mock, tokenService: NewTokenService(cfg)}

	mobile := "09123456789"

	if err := us.CreateUser(mobile); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := us.CreateUser(mobile); err == nil {
		t.Fatal("expected error for duplicate user")
	}

	token, err := us.Login(mobile)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}

	user, err := us.GetUser(mobile)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if user.Mobile != mobile {
		t.Fatalf("expected mobile=%s, got %s", mobile, user.Mobile)
	}

	users, err := us.GetUsers(1, 10)
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}
