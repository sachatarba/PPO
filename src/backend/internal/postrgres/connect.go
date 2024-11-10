package postrgres_adapter

import (
	"fmt"

	"github.com/sachatarba/course-db/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IPostgresConnector interface {
	Connect() (*gorm.DB, error)
}

type PostgresConnector struct {
	Conf *config.PostgresConfig
}

func (connector *PostgresConnector) Connect() (*gorm.DB, error) {
	conf := connector.Conf

	dsn := fmt.Sprintf(
		"host=172.17.0.1 port=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Port,
		conf.User,
		conf.DBName,
		conf.Password,
		conf.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
