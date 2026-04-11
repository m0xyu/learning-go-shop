package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m0xyu/learning-go-shop/internal/config"
	"github.com/m0xyu/learning-go-shop/internal/database"
	"github.com/m0xyu/learning-go-shop/internal/logger"
	"github.com/m0xyu/learning-go-shop/internal/providers"
	"github.com/m0xyu/learning-go-shop/internal/server"
	"github.com/m0xyu/learning-go-shop/internal/services"
)

func main() {
	log := logger.New()
	ctg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&ctg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}

	defer mainDB.Close()
	gin.SetMode(ctg.Server.GinMode)

	authService := services.NewAuthService(db, ctg)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	uploadService := services.NewUploadService(providers.NewLocalUploadProvider(ctg.Upload.Path))

	srv := server.New(ctg, db, &log, authService, productService, userService, uploadService)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", ctg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// HTTPサーバーを別のゴルーチンで起動
	go func() {
		log.Info().Str("port", ctg.Server.Port).Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	// チャンネルの作成とシグナルの通知設定
	quit := make(chan os.Signal, 1)
	// OSの割り込みシグナルを受け取るように設定
	signal.Notify(quit, os.Interrupt)
	// シグナルを受け取るまでブロック
	<-quit
	log.Info().Msg("shutting down server")

	// タイムアウト付きのコンテキストを作成して、HTTPサーバーのシャットダウンを実行
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// HTTPサーバーのシャットダウンを実行
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
		return
	}
	log.Info().Msg("shutting down database")
}
