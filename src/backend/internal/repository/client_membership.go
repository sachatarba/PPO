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

type ClientMembershipRepo struct {
	db *gorm.DB
}

func NewClientMembershipRepo(db *gorm.DB) service.IClientMembershipsRepository {
	return &ClientMembershipRepo{
		db: db,
	}
}

func (r *ClientMembershipRepo) CreateNewClientMembership(ctx context.Context, clientMembership entity.ClientMembership) error {
	conv := orm.NewClientMembershipConverter()
	membership := conv.ConvertFromEntity(clientMembership)
	membership.MembershipTypeID = membership.MembershipType.ID
	membership.MembershipType = orm.MembershipType{}
	tx := r.db.Create(&membership)

	return tx.Error
}

func (r *ClientMembershipRepo) ChangeClientMembership(ctx context.Context, clientMembership entity.ClientMembership) error {
	conv := orm.NewClientMembershipConverter()
	membership := conv.ConvertFromEntity(clientMembership)
	tx := r.db.WithContext(ctx).Model(orm.ClientMembership{ID: membership.ID}).Updates(
		map[string]interface{}{
			"start_date": membership.StartDate,
			"end_date": membership.EndDate,
		},
	)

	return tx.Error
}

func (r *ClientMembershipRepo) DeleteClientMembership(ctx context.Context, clientMembershipID uuid.UUID) error {
	// r.db.Begin()
	tx := r.db.WithContext(ctx).Delete(&orm.ClientMembership{
		ID: clientMembershipID,
	})
	// r.db.Commit()
	return tx.Error
}

func (r *ClientMembershipRepo) GetClientMembershipByID(ctx context.Context, clientMembershipID uuid.UUID) (entity.ClientMembership, error) {
	membership := orm.ClientMembership{
		ID: clientMembershipID,
	}
	tx := r.db.WithContext(ctx).
	Preload("MembershipType").
	First(&membership)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(service.ErrGetByIDNotFound, tx.Error)
	}

	if tx.Error != nil {
		return entity.ClientMembership{}, tx.Error
	}

	conv := orm.NewClientMembershipConverter()

	return conv.ConvertToEntity(membership), tx.Error
}

func (r *ClientMembershipRepo) ListClientMembershipsByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.ClientMembership, error) {
	var memberships []orm.ClientMembership
	
	tx := r.db.WithContext(ctx).Preload("MembershipType").Where(&orm.ClientMembership{ClientID: clientID}).Find(&memberships)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(tx.Error, service.ErrListByIDNotFound)
	}

	if tx.Error != nil {
		return []entity.ClientMembership{}, tx.Error
	}

	conv := orm.NewClientMembershipConverter()
	clientMemberships := conv.ConvertToEntitySlice(memberships)

	return clientMemberships, tx.Error
}
