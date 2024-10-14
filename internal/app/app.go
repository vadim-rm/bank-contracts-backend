package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vadim-rm/bmstu-web-backend/internal/config"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/external_routes"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/middleware"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/router"
	"github.com/vadim-rm/bmstu-web-backend/pkg/logger"
	"github.com/vadim-rm/bmstu-web-backend/pkg/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("error loading config: %s", err.Error())
	}

	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
				cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.Port)),
		&gorm.Config{})
	if err != nil {
		logger.Fatalf("error connecting to postgres: %s", err.Error())
	}

	redisClient, err := redis.New(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		User:     cfg.Redis.User,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		logger.Fatalf("error connecting to redis: %s", err.Error())
	}

	err = db.SetupJoinTable(&entity.Account{}, "Contracts", &entity.AccountContracts{})
	if err != nil {
		logger.Fatalf("error setting up join table: %s", err.Error())
		return
	}

	err = db.AutoMigrate(&entity.Contract{}, &entity.User{}, &entity.Account{})
	if err != nil {
		logger.Fatalf("error migrating entities: %s", err.Error())
		return
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.Id, cfg.Minio.Secret, ""),
		Secure: false, // Use true if using HTTPS
	})
	if err != nil {
		logger.Fatalf("error initializing minio: %s", err.Error())
		return
	}

	contractRepository := repository.NewContractImpl(db)
	accountRepository := repository.NewAccountImpl(db)
	accountContractsRepository := repository.NewAccountContractsImpl(db)
	usersRepository := repository.NewUserImpl(db)
	tokenRepository := repository.NewTokenImpl(
		repository.TokenConfig{
			ExpiresIn: cfg.Jwt.ExpiresIn,
			Token:     cfg.Jwt.Token,
			Issuer:    cfg.Jwt.Issuer,
		},
		redisClient,
	)

	imageRepository := repository.NewImageImpl(minioClient, repository.ImageConfig{
		BucketName: cfg.Minio.BucketName,
		BaseUrl:    cfg.Minio.BaseUrl,
	})

	contractService := service.NewContractImpl(contractRepository, accountRepository, imageRepository)
	accountService := service.NewAccountImpl(accountRepository)
	accountContractsService := service.NewAccountContractsImpl(accountContractsRepository, accountRepository)
	usersService := service.NewUserImpl(usersRepository, tokenRepository)

	accountHandler := handler.NewAccountImpl(accountService)
	contractHandler := handler.NewContractImpl(contractService, accountService)
	accountContractsHandler := handler.NewAccountContractsImpl(accountContractsService)
	usersHandler := handler.NewUserImpl(usersService)

	engine := router.New(router.Config{
		DebugCors: cfg.App.Debug,
	})

	authMiddleware := middleware.NewAuthMiddleware(tokenRepository)

	external_routes.Initialize(engine,
		authMiddleware,
		contractHandler, accountHandler, accountContractsHandler, usersHandler,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		Handler: engine.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
