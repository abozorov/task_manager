package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abozorov/task_manager/internal/config"
	"github.com/abozorov/task_manager/internal/handlers"
	"github.com/abozorov/task_manager/internal/models"
	"github.com/abozorov/task_manager/internal/repo"
	"github.com/abozorov/task_manager/internal/service"
	"github.com/abozorov/task_manager/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// get config
	cfg, err := config.NewConfig("internal/config/config.env")
	if err != nil {
		log.Fatal(err)
	}

	//connect to db server
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// migration db
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// create logger
	logger, err := logger.NewLogger(true)
	if err != nil {
		log.Println("Func main", zap.Error(err))
		return
	}

	// create layers
	repo := repo.NewTaskRepo(db)
	service := service.NewTaskService(repo)
	tashHandler := handlers.NewTaskHandler(service, logger)

	// create router
	router := handlers.NewRouter(*tashHandler)
	server := &http.Server{
		Addr:    cfg.HttpHost,
		Handler: router,
	}

	// statrt server
	go func() {
		logger.Info(fmt.Sprintf("Server started localhost:%s started", server.Addr))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("main", zap.Error(err))
			return
		}
	}()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("Shutdown server started")
	stopCtx, stopCancle := context.WithTimeout(context.Background(), time.Second*5)
	defer stopCancle()

	server.Shutdown(stopCtx)

	logger.Info("Server shutdown completed")
}
