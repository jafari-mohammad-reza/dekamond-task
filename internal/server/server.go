package server

import (
	"context"
	"dekamond-task/internal/config"
	"dekamond-task/internal/dto"
	"dekamond-task/internal/service"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

const maxFailedAttempts = 3

type Server struct {
	cfg            *config.Config
	e              *echo.Echo
	failedAttempts map[string]int
	registerOtps   map[string]OtpItem
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:            cfg,
		failedAttempts: make(map[string]int),
		registerOtps:   make(map[string]OtpItem),
	}
}

type Resp map[string]any
type OtpItem struct {
	Otp       int
	CreatedAt time.Time
}

func (s *Server) Start(ctx context.Context) error {

	go s.checkOtps(ctx)
	e := echo.New()
	s.e = e
	userService, err := service.NewUserService(s.cfg)
	if err != nil {
		return fmt.Errorf("failed to create user service: %w", err)
	}
	e.GET("/users", func(ctx echo.Context) error {
		page := getQueryNum(ctx, "page", 1)
		limit := getQueryNum(ctx, "limit", 10)
		users, err := userService.GetUsers(page, limit)
		if err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		return ctx.JSON(200, Resp{"users": users})
	})
	e.GET("/users/:mobile", func(ctx echo.Context) error {
		mobile := ctx.Param("mobile")
		user, err := userService.GetUser(mobile)
		if err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		return ctx.JSON(200, Resp{"user": user})
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
		s.registerOtps[request.MobileNumber] = OtpItem{
			Otp:       otp,
			CreatedAt: time.Now(),
		}
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
		if failed, ok := s.failedAttempts[request.MobileNumber]; ok {
			if failed >= maxFailedAttempts {
				return ctx.JSON(403, Resp{"message": "Too many failed attempts"})
			}
		}
		otp, ok := s.registerOtps[request.MobileNumber]
		if !ok {
			return ctx.JSON(403, Resp{"message": "Invalid input"})
		}
		if otp.Otp != request.OTP {
			s.failedAttempts[request.MobileNumber]++
			return ctx.JSON(403, Resp{"message": "Invalid input"})
		}
		token, err := userService.Login(request.MobileNumber)
		if err != nil {
			return ctx.JSON(400, Resp{"message": err.Error()})
		}
		return ctx.JSON(200, Resp{"message": "User logged in successfully", "token": token})
	})
	return e.Start(fmt.Sprintf(":%d", s.cfg.Port))
}
func (s *Server) checkOtps(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("main context cancelled")
			break
		case <-ticker.C:
			for key, val := range s.registerOtps {
				if time.Since(val.CreatedAt) > 3*time.Minute {
					delete(s.registerOtps, key)
				}
			}
		default:
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
func getQueryNum(ctx echo.Context, item string, def int) int {
	itemVal := ctx.QueryParams().Get(item)
	itemNum, err := strconv.Atoi(itemVal)
	if err != nil {
		return def
	}
	return itemNum
}
