package v2

import (
	"log"
	"time"

	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/orm"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	"github.com/sachatarba/course-db/internal/server/v2"
	"gorm.io/gorm"
)

type ApiServer struct {
	conf config.Config
}

func (api *ApiServer) Run() {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}
	// redisConnector := redis_adapter.RedisConnector{
	// 	Conf: conf.RedisConf,
	// }
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = postgresConnector.Connect()
		if err != nil {
			log.Print("Cant connect postgres", err)
			time.Sleep(time.Second)
			// return
		}
	}
	if err != nil {
		log.Print("Cant connect postgres", err)
		return
	}

	// rdb := redisConnector.Connect()
	// if rdb == nil {
	// 	log.Print("Cant connect redis", err)
	// 	return
	// }

	postgresMigrator := postrgres_adapter.PostgresMigrator{
		DB:     db,
		Tables: orm.TablesORM,
	}

	err = postgresMigrator.Migrate()
	if err != nil {
		log.Print("Cant migrate", err)
		return
	}

	// paymentConfig := conf.PaymentConfig
	// paymentHandler := handler.NewPaymentHandler(paymentConfig)

	// services, err := director.NewServices()
	// if err != nil {
	// 	log.Print("Cant create services", err)
	// 	return
	// }

	log.Print("bebra")

	server := server.Server{
		// PaymentHandler: paymentHandler,
		Handlers: &ApiHandlers{
			Postgres: db,
		},
		Conf: &config.ServerConfig{
			Port: "8081",
		},
	}

	server.Run()
}
