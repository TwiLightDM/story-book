package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"story-book/internal/config"
	"story-book/internal/middlewares"
	"story-book/internal/services/userservice"
	"story-book/package/databases/postgres"
	"story-book/package/services/encryptservice"
	"story-book/package/services/jwtservice"
	"story-book/package/services/validateservice"
	"syscall"
	"time"

	_ "story-book/internal/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func Run(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())
	db, err := postgres.InitDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	if err != nil {
		return err
	}

	jwtService := jwtservice.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessDuration, cfg.JWT.RefreshDuration)
	encryptService := encryptservice.NewEncryptionService(cfg.SaltLength)
	validateService := validateservice.NewValidationService(cfg.MinPasswordSize)

	authMiddleware := middlewares.AuthMiddleware(jwtService)

	userRepository := userservice.NewUserRepository(db)
	userService := userservice.NewUserService(userRepository, jwtService, encryptService, validateService)
	userHandler := userservice.NewUserHandler(userService)

	registerRoutes(e, authMiddleware, userHandler)

	server := &http.Server{
		Addr:    ":" + cfg.BackendPort,
		Handler: e,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func(db *gorm.DB) {
		log.Printf("Backend started on :%s", cfg.BackendPort)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}

		if errors.Is(err, http.ErrServerClosed) {
			err = postgres.CloseDB(db)
			if err != nil {
				log.Fatalf("failed to close database connection: %v", err)
			}
		}
	}(db)

	<-ctx.Done()
	log.Println("Shutting down backend server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Backend server stopped gracefully")
	return nil
}

func registerRoutes(e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	userHandler *userservice.UserHandler,
) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	public := e.Group("/auth")
	public.POST("/login", userHandler.Login)
	public.POST("/signup", userHandler.SignUp)
	public.POST("/refresh", userHandler.Refresh, authMiddleware)
	public.POST("/reset-password", userHandler.ResetPassword, authMiddleware)

	users := e.Group("/users", authMiddleware)
	users.GET("/me", userHandler.ReadSelf)
	users.GET("/:id", userHandler.ReadUser)
	users.PUT("/me", userHandler.UpdateUser)
	users.PATCH("/me/password", userHandler.ChangePassword)
	users.DELETE("/me", userHandler.DeleteUser)
}
