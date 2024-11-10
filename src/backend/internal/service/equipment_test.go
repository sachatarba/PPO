package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/stretchr/testify/mock"
)

type EquipmentServiceSuite struct {
	suite.Suite
}

func (s *EquipmentServiceSuite) TestCreateNewEquipment(t provider.T) {
	t.Title("[CreateNewEquipment] Successfully creates equipment")
	t.Tags("equipment", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created equipment", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		equipment := builder.NewEquipmentBuilder().Build()

		equipmentRepoMock.On("CreateNewEquipment", mock.Anything, equipment).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "equipment", equipment)

		err := equipmentService.CreateNewEquipment(ctx, equipment)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		invalidEquipment := builder.NewEquipmentBuilder().Invalid().Build()

		err := equipmentService.CreateNewEquipment(ctx, invalidEquipment)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *EquipmentServiceSuite) TestChangeEquipment(t provider.T) {
	t.Title("[ChangeEquipment] Successfully changes equipment")
	t.Tags("equipment", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed equipment", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		equipment := builder.NewEquipmentBuilder().Build()

		equipmentRepoMock.On("ChangeEquipment", mock.Anything, equipment).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "equipment", equipment)

		err := equipmentService.ChangeEquipment(ctx, equipment)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		invalidEquipment := builder.NewEquipmentBuilder().Invalid().Build()

		err := equipmentService.ChangeEquipment(ctx, invalidEquipment)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *EquipmentServiceSuite) TestDeleteEquipment(t provider.T) {
	t.Title("[DeleteEquipment] Successfully deletes equipment")
	t.Tags("equipment", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted equipment", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		equipmentID := uuid.New()

		equipmentRepoMock.On("DeleteEquipment", mock.Anything, equipmentID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "equipmentID", equipmentID)

		err := equipmentService.DeleteEquipment(ctx, equipmentID)

		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentServiceSuite) TestGetEquipmentByID(t provider.T) {
	t.Title("[GetEquipmentByID] Successfully retrieves equipment by ID")
	t.Tags("equipment", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved equipment", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		equipmentID := uuid.New()
		expectedEquipment := builder.NewEquipmentBuilder().Build()

		equipmentRepoMock.On("GetEquipmentByID", mock.Anything, equipmentID).Return(expectedEquipment, nil)

		sCtx.WithNewParameters("ctx", ctx, "equipmentID", equipmentID)

		equipment, err := equipmentService.GetEquipmentByID(ctx, equipmentID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedEquipment, equipment)
	})
}

func (s *EquipmentServiceSuite) TestListEquipmentsByGymID(t provider.T) {
	t.Title("[ListEquipmentsByGymID] Successfully lists equipments by gym ID")
	t.Tags("equipment", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed equipments", func(sCtx provider.StepCtx) {
		equipmentRepoMock := &mocks.IEquipmentRepository{}
		equipmentService := &EquipmentService{equipmentRepoMock}

		ctx := context.TODO()
		gymID := uuid.New()
		expectedEquipments := []entity.Equipment{
			builder.NewEquipmentBuilder().Build(),
			builder.NewEquipmentBuilder().Build(),
		}

		equipmentRepoMock.On("ListEquipmentsByGymID", mock.Anything, gymID).Return(expectedEquipments, nil)

		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		equipments, err := equipmentService.ListEquipmentsByGymID(ctx, gymID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedEquipments, equipments)
	})
}

func TestEquipmentServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(EquipmentServiceSuite))
}
