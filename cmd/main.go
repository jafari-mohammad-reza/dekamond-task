package main

import (
	"context"
	"dekamond-task/internal/config"
	"dekamond-task/internal/db"
	"dekamond-task/internal/server"
	"fmt"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // hierarchical context handling
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
		os.Exit(1)
	}
	db, err := db.NewDB(cfg)
	if err != nil {
		fmt.Printf("failed to create database: %s", err.Error())
		os.Exit(1)
	}
	if err := db.InitTables(); err != nil {
		fmt.Printf("failed to create tables: %s", err.Error())
		os.Exit(1)
	}
	server := server.NewServer(cfg)
	err = server.Start(ctx)
	if err != nil {
		fmt.Printf("failed to start server: %s", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	defer cancel()
	defer server.Stop(ctx)
}
