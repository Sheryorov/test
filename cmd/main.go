package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sheryorov/test/config"
	"github.com/sheryorov/test/internal/entity"
	"github.com/sheryorov/test/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// swagger embed files
)

var semaphore chan struct{}

func LimitMiddleware(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			if len(semaphore) == limit {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"message": "Timeout for method get"})
				return
			}
			semaphore <- struct{}{}
		}
		c.Next()
	}
}

// @title           Swagger USER API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /
func main() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	cfg := config.InitConfig()
	defer logger.Sync() // flushes buffer, if any
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		sugar.Fatalf("failed to connect database: %v", err)
	}
	semaphore = make(chan struct{}, 5)
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			<-ticker.C
			for i := 0; i < cfg.GETMethodPerMinute; i++ {
				<-semaphore
			}
		}
	}()
	db.AutoMigrate(&entity.User{})
	webAPIUser := usecase.NewUserIteractor(db, logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := gin.Default()
	r.Use(LimitMiddleware(cfg.GETMethodPerMinute))
	r.GET("/users", webAPIUser.GetAllUsers)
	r.GET("/users/:id", webAPIUser.GetUserByID)
	r.POST("/user", webAPIUser.CreateUser)
	r.PUT("/users/:id", webAPIUser.UpdateUser)
	r.DELETE("/users/:id", webAPIUser.DeleteUser)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	sugar.Infoln("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatal("Server forced to shutdown: ", err)
	}
}
