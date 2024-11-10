package repository

// import (
// 	"context"
// 	"log"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/sachatarba/course-db/internal/repository/builder"
// 	"github.com/sachatarba/course-db/internal/repository/mocks"
// )

// func Bebra() {
// 	db, mock := mocks.NewMockDB()

// 	// trainerID := uuid.New()
// 	// gymID := uuid.New()
// 	// trainer := builder.NewTrainerBuilder().
// 	// 	SetID(trainerID).
// 	// 	SetFullname("John Doe").
// 	// 	SetEmail("john@example.com").
// 	// 	SetPhone("+1-800-555-1234").
// 	// 	SetQualification("Certified Trainer").
// 	// 	SetUnitPrice(100).
// 	// 	SetGymsID([]uuid.UUID{gymID}).
// 	// 	Build()

// 	// mock.ExpectBegin()
// 	// mock.ExpectExec(`^(.+)$`).
// 	// 	// WithArgs(trainer.ID, trainer.Fullname, trainer.Email, trainer.Phone, trainer.Qualification, trainer.Qualification).
// 	// 	WillReturnResult(sqlmock.NewResult(1, 1))
// 	// trainerID := uuid.New()
// 	// gymID := uuid.New()

// 	// trainer := builder.NewTrainerBuilder().
// 	// 	SetID(trainerID).
// 	// 	SetFullname("John Doe").
// 	// 	SetEmail("john@example.com").
// 	// 	SetPhone("+1-800-555-1234").
// 	// 	SetQualification("Certified Trainer").
// 	// 	SetUnitPrice(100).
// 	// 	SetGymsID([]uuid.UUID{gymID}).
// 	// 	Build()

// 	mock.ExpectExec(`^INSERT INTO "trainers" (.+)$`).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	mock.ExpectExec(`^INSERT INTO "gym_trainers" (.+)$`).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// mock.ExpectExec("h")

// 	repo := NewTrainerRepo(db)
// 	// log.Println("bebra")

// 	err := repo.RegisterNewTrainer(context.TODO(), builder.NewTrainerBuilder().Build())

// 	log.Println(err)
// }
