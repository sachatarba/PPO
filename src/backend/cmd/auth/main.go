package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	pb "github.com/sachatarba/course-db/internal/api/grpc"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/delivery/v2/rest"
	redis_adapter "github.com/sachatarba/course-db/internal/redis"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"

	"google.golang.org/grpc"
)

func main() {
	conf := config.NewConfFromEnv()

	redisConnector := redis_adapter.RedisConnector{
		Conf: conf.RedisConf,
	}

	rdb := redisConnector.Connect()
	if rdb == nil {
		log.Fatal("Cant connect redis")
		return
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s",
			conf.GrpcClientsServerConfig.Host,
			conf.GrpcClientsServerConfig.Port),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewClientServiceClient(conn)

	repo := repository.NewSessionRepo(rdb)
	repoCode := repository.NewCodeRepo(rdb)
	smtpService := service.NewSmtpService(conf.SmtpConfig)

	service := service.NewAuthorizationNewService(repo, client, smtpService, repoCode)

	h := rest.NewAuthHandler(service)

	router := gin.Default()

	v2 := router.Group("/api/v2")
	{
		v2.POST("/register", h.RegisterNewUser)
		v2.POST("/login", h.Login)
		v2.GET("/logout", h.Logout)

		v2.GET("/isauthorize", h.IsAuthorize)
		v2.POST("/confirm", h.Confirm2FA)
		v2.POST("/change_password", h.ChangePassword)
	}

	router.Run(fmt.Sprintf(":%s", conf.AuthServerConfig.Port))
}
