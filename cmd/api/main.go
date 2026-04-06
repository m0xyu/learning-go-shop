package main

import (
	"github.com/gin-gonic/gin"
	"github.com/m0xyu/learning-go-shop/internal/config"
	"github.com/m0xyu/learning-go-shop/internal/database"
	"github.com/m0xyu/learning-go-shop/internal/logger"
)

func main() {
	log := logger.New()
	ctg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(ctg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to datbase")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}

	defer mainDB.Close()
	gin.SetMode(ctg.Server.GinMode)

	log.Info().Msg("starting server")
}
