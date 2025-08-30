package server

import (
	"context"
	"dekamond-task/internal/config"
	"dekamond-task/internal/dto"
	"dekamond-task/internal/service"
	"fmt"
	"math/rand"

	"github.com/labstack/echo"
)

const maxFailedAttempts = 3

type Server struct {
	cfg *config.Config
	e   *echo.Echo
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

type Resp map[string]any

func (s *Server) Start(ctx context.Context) error {
	failedAttempts := map[string]int{}
	registerOtps := map[string]int{}
	e := echo.New()
	s.e = e
	userService, err := service.NewUserService(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to create user service: %w", err)
	}
	e.GET("/users", func(ctx echo.Context) error {
		// query := ctx.QueryParams()
		// users := userService.GetUsers(filters)
		return ctx.JSON(200, Resp{"message": "User registered successfully"})
	})
	e.POST("/register", func(ctx echo.Context) error {
		var request dto.RegisterRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		if err := dto.ValidateModel(&request); err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}

		if err := userService.CreateUser(request.MobileNumber); err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		otp := rand.Intn(999999)
		fmt.Printf("mobile numner %s otp: %d\n", request.MobileNumber, otp)
		registerOtps[request.MobileNumber] = otp
		return ctx.JSON(200, Resp{"message": "User registered successfully"})
	})
	e.POST("/login", func(ctx echo.Context) error {
		var request dto.LoginRequest
		if err := ctx.Bind(&request); err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		if err := dto.ValidateModel(&request); err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		if failed, ok := failedAttempts[request.MobileNumber]; ok {
			if failed >= maxFailedAttempts {
				return ctx.JSON(403, Resp{"message": "Too many failed attempts"})
			}
		}
		token, err := userService.Login(request.MobileNumber, request.OTP)
		if err != nil {
			failedAttempts[request.MobileNumber]++
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		return ctx.JSON(200, Resp{"message": "User logged in successfully", "token": token})
	})
	return e.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
