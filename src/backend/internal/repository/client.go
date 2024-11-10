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

type ClientRepo struct {
	db *gorm.DB
}

func NewClientRepo(db *gorm.DB) service.IClientRepository {
	return &ClientRepo{
		db: db,
	}
}

func (r *ClientRepo) RegisterNewClient(ctx context.Context, client entity.Client) error {
	conv := orm.NewClientConverter()
	tx := r.db.WithContext(ctx).Create(conv.ConvertFromEntity(client))

	return tx.Error
}

func (r *ClientRepo) ChangeClient(ctx context.Context, client entity.Client) error {
	conv := orm.NewClientConverter()
	tx := r.db.WithContext(ctx).Model(&orm.Client{ID: client.ID}).Updates(conv.ConvertFromEntity(client))

	return tx.Error
}

func (r *ClientRepo) DeleteClient(ctx context.Context, clientID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&orm.Client{
		ID: clientID,
	})

	return tx.Error
}

func (r *ClientRepo) GetClientByID(ctx context.Context, clientID uuid.UUID) (entity.Client, error) {
	client := orm.Client{
		ID: clientID,
	}
	
	tx := r.db.WithContext(ctx).First(&client)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(service.ErrGetByIDNotFound, tx.Error)
	}

	if tx.Error != nil {
		return entity.Client{}, tx.Error
	}

	conv := orm.NewClientConverter()

	return conv.ConvertToEntity(client), tx.Error
}

func (r *ClientRepo) GetClientByLogin(ctx context.Context, login string) (entity.Client, error) {
	client := &orm.Client{}
	tx := r.db.WithContext(ctx).Where("login = ?", login).First(client)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(service.ErrNoSuchClient, tx.Error)
	}

	if tx.Error != nil {
		return entity.Client{}, tx.Error
	}

	conv := orm.NewClientConverter()

	return conv.ConvertToEntity(*client), tx.Error
}

func (r *ClientRepo) ListClients(ctx context.Context) ([]entity.Client, error) {
	var clients []orm.Client

	tx := r.db.WithContext(ctx).Find(&clients)
	if tx.Error != nil {
		return []entity.Client{}, tx.Error
	}

	conv := orm.NewClientConverter()
	return conv.ConvertToEntitySlice(clients), tx.Error
}
