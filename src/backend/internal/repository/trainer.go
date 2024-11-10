package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	"github.com/sachatarba/course-db/internal/service"
	"gorm.io/gorm"
)

type TrainerRepo struct {
	db        *gorm.DB
	converter entity.IConverter[entity.Trainer, orm.Trainer]
}

func NewTrainerRepo(db *gorm.DB) service.ITrainerRepository {
	return &TrainerRepo{
		db:        db,
		converter: orm.NewTrainerConverter(),
	}
}

func (r *TrainerRepo) RegisterNewTrainer(ctx context.Context, trainer entity.Trainer) error {
	var err error
	trainerOrm := r.converter.ConvertFromEntity(trainer)
	trainerOrm.Gyms = []*orm.Gym{}

	// log.Println("berba")
	// log.Fatalf("bebra")
	// tx := r.db.WithContext(ctx).Create(&trainerOrm)
	// r.db.Logger.LogMode()

	for _, id := range trainer.GymsID {
		log.Println("here")
		err = r.db.Model(&orm.Gym{
			ID: id,
		}).Association("Trainers").Append(&trainerOrm)

		if err != nil {
			break
		}
	}

	return err
}

func (r *TrainerRepo) ChangeTrainer(ctx context.Context, trainer entity.Trainer) error {
	trainerOrm := r.converter.ConvertFromEntity(trainer)
	trainerOrm.Gyms = []*orm.Gym{}
	tx := r.db.WithContext(ctx).Save(&trainerOrm)

	return tx.Error
}

func (r *TrainerRepo) DeleteTrainer(ctx context.Context, trainerID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&orm.Trainer{ID: trainerID})

	return tx.Error
}

func (r *TrainerRepo) GetTrainerByID(ctx context.Context, trainerID uuid.UUID) (entity.Trainer, error) {
	trainerOrm := orm.Trainer{
		ID: trainerID,
	}

	tx := r.db.WithContext(ctx).First(&trainerOrm)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return entity.Trainer{}, service.ErrGetByIDNotFound
	}

	return r.converter.ConvertToEntity(trainerOrm), tx.Error
}

func (r *TrainerRepo) ListTrainers(ctx context.Context) ([]entity.Trainer, error) {
	var trainerOrms []orm.Trainer
	tx := r.db.WithContext(ctx).Find(&trainerOrms)

	return r.converter.ConvertToEntitySlice(trainerOrms), tx.Error
}

func (r *TrainerRepo) ListTrainersByGymID(ctx context.Context, gymID uuid.UUID) ([]entity.Trainer, error) {
	var trainersOrmPtr []*orm.Trainer
	err := r.db.Model(&orm.Gym{
		ID: gymID,
	}).Association("Trainers").Find(&trainersOrmPtr)

	if err != nil {
		return nil, err
	}
	trainersOrm := make([]orm.Trainer, len(trainersOrmPtr))
	for i, trainer := range trainersOrmPtr {
		trainersOrm[i] = *trainer
	}

	return r.converter.ConvertToEntitySlice(trainersOrm), err
}
