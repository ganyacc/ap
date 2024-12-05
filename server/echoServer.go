package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganyacc/ap.git/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func StartServer() error {
	echoApp := echo.New()

	echoApp.Logger.SetLevel(log.DEBUG)
	echoApp.Use(middleware.Recover())
	echoApp.Use(middleware.Logger())

	//
	InitializeRoutes(echoApp)

	return echoApp.Start(":8080")

}

func InitializeRoutes(e *echo.Echo) {

	userHandler := controller.NewUserHttpHandler()

	e.POST("/v1/user/signup", userHandler.SignUp)
	e.POST("/v1/user/signin", userHandler.SignIn)
	e.GET("/v1/user/refresh_token", userHandler.RefreshToken)
}

func ShutdownServer(e *echo.Echo) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown signal received")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	//Attempt to gracefully shutdown the server
	if err := echo.New().Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Print("Server exited gracefully")

}
