package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/sachatarba/course-db/internal/config"
	handler "github.com/sachatarba/course-db/internal/delivery/v2/rest"
	"github.com/swaggo/files"
)

type ServiceHandlers interface {
	InitHandlers(gin.IRouter)
}

type Server struct {
	Handlers ServiceHandlers
	Conf     *config.ServerConfig
}

func (server *Server) Run() {
	log.Println("Server starting..")
	router := gin.Default()
	router.Use(handler.CORSMiddleware)
	router.Static("/docs", "./internal/server/v2/docs")
	url := ginSwagger.URL("/docs/openapi.yaml")
	log.Println(os.Getwd())
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// rouer.GET()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	server.Handlers.InitHandlers(router)

	serv := &http.Server{
		Addr:    ":"+server.Conf.Port,
		Handler: router,
	}

	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
