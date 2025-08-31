package db

import (
	"context"
	"dekamond-task/internal/config"

	"os"
	"testing"
	"time"
)

func setupDB(t *testing.T) Repo {
	tmpFile := "test.db"
	os.Remove(tmpFile)

	cfg := &config.Config{
		Database: struct{ Url string }{Url: tmpFile},
	}
	database, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	if err := database.(*DB).InitTables(); err != nil {
		t.Fatalf("failed to init tables: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tmpFile)
	})

	return database
}

func TestDB_CreateGetUser(t *testing.T) {
	repo := setupDB(t)
	ctx := context.Background()
	mobile := "09123456789"

	if err := repo.CreateUser(ctx, mobile); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user, err := repo.GetUser(ctx, mobile)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Mobile != mobile {
		t.Fatalf("expected mobile=%s, got %s", mobile, user.Mobile)
	}

	_, err = repo.GetUser(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestDB_GetUsers(t *testing.T) {
	repo := setupDB(t)
	ctx := context.Background()
	mobiles := []string{"09111111111", "09222222222", "09333333333"}

	for _, m := range mobiles {
		if err := repo.CreateUser(ctx, m); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	users, err := repo.GetUsers(ctx, 1, 2)
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	allUsers, err := repo.GetUsers(ctx, 1, 10)
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}

	if len(allUsers) != len(mobiles) {
		t.Fatalf("expected %d users, got %d", len(mobiles), len(allUsers))
	}
}
