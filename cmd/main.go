package main

import (
	"compare/config"
	"compare/internal/db"
	"compare/internal/handler"
	"compare/internal/jobs"
	"compare/internal/repository"
	"compare/internal/service"
	"compare/pkg/logging"
	"github.com/gin-gonic/gin"
	"net"
)

var logger = logging.GetLogger()

func main() {
	router := gin.Default()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Errorf("failed while getting config: %v", err)
		return
	}

	dbConn, err := db.GetDbConn()

	newRepository := repository.NewRepository(dbConn)

	newService := service.NewService(newRepository)

	newHandler := handler.NewHandler(router, newService)
	newHandler.InitRoutes()

	go jobs.StartJobs(newService)

	addr := net.JoinHostPort(cfg.Host, cfg.Port)

	logger.Fatalln(router.Run(addr))
}
