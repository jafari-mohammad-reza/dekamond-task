package server

import (
	"context"
	_ "dekamond-task/docs"
	"dekamond-task/internal/config"
	"dekamond-task/internal/dto"
	"dekamond-task/internal/service"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const maxFailedAttempts = 3

type Server struct {
	cfg            *config.Config
	e              *echo.Echo
	failedAttempts map[string]int
	registerOtps   map[string]OtpItem
	userService    service.IUserService
	mu             sync.Mutex
}

func NewServer(cfg *config.Config) (*Server, error) {
	userService, err := service.NewUserService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %w", err)
	}
	return &Server{
		cfg:            cfg,
		failedAttempts: make(map[string]int),
		registerOtps:   make(map[string]OtpItem),
		userService:    userService,
		mu:             sync.Mutex{},
	}, nil
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
	e.GET("/users", s.getUsers)
	e.GET("/users/:mobile", s.getByMobile)
	e.POST("/auth", s.auth)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) getUsers(ctx echo.Context) error {
	page := getQueryNum(ctx, "page", 1)
	limit := getQueryNum(ctx, "limit", 10)
	users, err := s.userService.GetUsers(page, limit)
	if err != nil {
		return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
	}
	return ctx.JSON(200, dto.UsersResponse{Users: users})
}

func (s *Server) getByMobile(ctx echo.Context) error {
	mobile := ctx.Param("mobile")
	user, err := s.userService.GetUser(mobile)
	if err != nil {
		return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
	}
	return ctx.JSON(200, dto.UserResponse{User: user})
}

func (s *Server) auth(ctx echo.Context) error {
	var request dto.AuthRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
	}
	if err := dto.ValidateModel(&request); err != nil {
		return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
	}

	userExists := s.userService.UserExists(request.MobileNumber)

	if request.OTP == nil {
		if !userExists {
			if err := s.userService.CreateUser(request.MobileNumber); err != nil {
				return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
			}
		}
		otp := rand.Intn(999999)
		fmt.Printf("mobile number %s otp: %d\n", request.MobileNumber, otp)
		s.mu.Lock()
		s.registerOtps[request.MobileNumber] = OtpItem{
			Otp:       otp,
			CreatedAt: time.Now(),
		}
		s.mu.Unlock()
		return ctx.JSON(200, dto.AuthResponse{Message: "OTP sent successfully"})
	}

	if failed, ok := s.failedAttempts[request.MobileNumber]; ok {
		if failed >= maxFailedAttempts {
			return ctx.JSON(403, dto.MessageResponse{Message: "Too many failed attempts"})
		}
	}
	otp, ok := s.registerOtps[request.MobileNumber]
	if !ok {
		return ctx.JSON(403, dto.MessageResponse{Message: "Invalid input"})
	}
	if otp.Otp != *request.OTP {
		s.mu.Lock()
		s.failedAttempts[request.MobileNumber]++
		s.mu.Unlock()
		return ctx.JSON(403, dto.MessageResponse{Message: "Invalid input"})
	}
	delete(s.registerOtps, request.MobileNumber)
	token, err := s.userService.Login(request.MobileNumber)
	if err != nil {
		return ctx.JSON(400, dto.MessageResponse{Message: err.Error()})
	}
	return ctx.JSON(200, dto.AuthResponse{
		Message: "User logged in successfully",
		Token:   token,
	})
}

func (s *Server) checkOtps(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("main context cancelled")
			return
		case <-ticker.C:
			s.mu.Lock()
			for key, val := range s.registerOtps {
				if time.Since(val.CreatedAt) > 3*time.Minute {
					delete(s.registerOtps, key)
				}
			}
			s.mu.Unlock()
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
