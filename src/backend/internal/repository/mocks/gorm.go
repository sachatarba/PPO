package mocks

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	}),
		&gorm.Config{
			// Logger: logger.New(
			// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
			// 	logger.Config{
			// 		SlowThreshold: time.Second,
			// 		LogLevel:      logger.Info,
			// 		Colorful:      true,
			// 	},
			// ),
		})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}
