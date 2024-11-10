package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	"github.com/sachatarba/course-db/internal/service"
	"gorm.io/gorm"
)

type EquipmentRepo struct {
	db *gorm.DB
}

func NewEquipmentRepo(db *gorm.DB) service.IEquipmentRepository {
	return &EquipmentRepo{
		db: db,
	}
}

func (r *EquipmentRepo) CreateNewEquipment(ctx context.Context, equipment entity.Equipment) error {
	equipmentOrm := orm.NewEquipmentConverter().ConvertFromEntity(equipment)
	err := r.db.WithContext(ctx).Model(
		&orm.Gym{
			ID: equipment.GymID,
		},
	).Association("Equipments").Append(&equipmentOrm)

	return err
}

func (r *EquipmentRepo) ChangeEquipment(ctx context.Context, equipment entity.Equipment) error {
	equipmentOrm := orm.NewEquipmentConverter().ConvertFromEntity(equipment)
	tx := r.db.WithContext(ctx).Model(&orm.Equipment{
		ID: equipment.ID,
	}).UpdateColumns(&orm.Equipment{
		Name: equipmentOrm.Name,
		Description: equipmentOrm.Description,
		GymID: equipmentOrm.GymID,
	})

	return tx.Error
}
func (r *EquipmentRepo) DeleteEquipment(ctx context.Context, equipmentID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&orm.Equipment{
		ID: equipmentID,
	})

	return tx.Error
}
func (r *EquipmentRepo) GetEquipmentByID(ctx context.Context, equipmentID uuid.UUID) (entity.Equipment, error) {
	equipment := &orm.Equipment{
		ID: equipmentID,
	}

	tx := r.db.WithContext(ctx).First(equipment)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(service.ErrGetByIDNotFound, tx.Error)
	}

	if tx.Error != nil {
		return entity.Equipment{}, tx.Error
	}

	conv := orm.NewEquipmentConverter()

	return conv.ConvertToEntity(*equipment), tx.Error
}

func (r *EquipmentRepo) ListEquipmentsByGymID(ctx context.Context, gymID uuid.UUID) ([]entity.Equipment, error) {
	var equipmentsOrm []orm.Equipment

	err := r.db.WithContext(ctx).Model(
		&orm.Gym{
			ID: gymID,
		},
	).Association("Equipments").Find(&equipmentsOrm)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.Join(err, service.ErrListByIDNotFound)
	}

	if err != nil {
		return []entity.Equipment{}, err
	}

	conv := orm.NewEquipmentConverter()

	return conv.ConvertToEntitySlice(equipmentsOrm), err
}
