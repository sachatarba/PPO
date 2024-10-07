package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sachatarba/course-db/internal/config"
	handler "github.com/sachatarba/course-db/internal/delivery/v1/rest"
)

type Server struct {
	Handler        *handler.Handler
	PaymentHandler *handler.PaymentHandler
	Conf           *config.ServerConfig
}

func (server *Server) Run() {
	router := gin.Default()
	router.Use(handler.CORSMiddleware)
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	server.PaymentHandler.InitPayment(router)
	server.Handler.InitApi(router)

	serv := &http.Server{
		Addr:    ":8080", 
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
