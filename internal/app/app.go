package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/vadim-rm/bmstu-web-backend/internal/config"
	"github.com/vadim-rm/bmstu-web-backend/internal/repository"
	"github.com/vadim-rm/bmstu-web-backend/internal/service"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/external_routes"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/handler"
	"github.com/vadim-rm/bmstu-web-backend/internal/transport/http/router"
	"github.com/vadim-rm/bmstu-web-backend/pkg/logger"
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

	contractRepository := repository.NewContractImpl()
	contractService := service.NewContractImpl(contractRepository)

	orderRepository := repository.NewAccountImpl()
	orderService := service.NewAccountImpl(orderRepository)

	orderHandler := handler.NewAccountImpl(orderService)
	contractHandler := handler.NewContractImpl(contractService, orderService)

	engine := router.New(router.Config{
		DebugCors:     cfg.App.Debug,
		TemplatesPath: cfg.App.TemplatesPath,
	})

	external_routes.Initialize(engine, contractHandler, orderHandler)

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
