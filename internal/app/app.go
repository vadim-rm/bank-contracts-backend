package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/config"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository/entity"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/external_routes"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/router"
	"github.com/vadim-rm/bmstu-web-backend/pkg/logger"
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
		logger.Fatalf("error loading config: %w", err)
	}

	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
				cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.Port)),
		&gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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

	contractRepository := repository.NewContractImpl(db)
	accountRepository := repository.NewAccountImpl(db)
	accountContractsRepository := repository.NewAccountContractsImpl(db)
	usersRepository := repository.NewUserImpl(db)

	contractService := service.NewContractImpl(contractRepository, accountRepository)
	accountService := service.NewAccountImpl(accountRepository)
	accountContractsService := service.NewAccountContractsImpl(accountContractsRepository, accountRepository)
	usersService := service.NewUserImpl(usersRepository)

	accountHandler := handler.NewAccountImpl(accountService)
	contractHandler := handler.NewContractImpl(contractService, accountService)
	accountContractsHandler := handler.NewAccountContractsImpl(accountContractsService)
	usersHandler := handler.NewUserImpl(usersService)

	engine := router.New(router.Config{
		DebugCors: cfg.App.Debug,
	})

	external_routes.Initialize(engine, contractHandler, accountHandler, accountContractsHandler, usersHandler)

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
