package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type MembershipTypeServiceSuite struct {
	suite.Suite
}

func (s *MembershipTypeServiceSuite) TestRegisterNewMembershipType(t provider.T) {
	t.Title("[RegisterNewMembershipType] Successfully registers a new membership type")
	t.Tags("membership_type", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new membership type", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		membershipType := builder.NewMembershipTypeBuilder().Build()

		repoMock.On("RegisterNewMembershipType", ctx, membershipType).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "membershipType", membershipType)

		err := service.RegisterNewMembershipType(ctx, membershipType)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		invalidMembershipType := builder.NewMembershipTypeBuilder().Invalid().Build()

		err := service.RegisterNewMembershipType(ctx, invalidMembershipType)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *MembershipTypeServiceSuite) TestChangeMembershipType(t provider.T) {
	t.Title("[ChangeMembershipType] Successfully changes membership type")
	t.Tags("membership_type", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed membership type", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		membershipType := builder.NewMembershipTypeBuilder().Build()

		repoMock.On("ChangeMembershipType", ctx, membershipType).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "membershipType", membershipType)

		err := service.ChangeMembershipType(ctx, membershipType)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		invalidMembershipType := builder.NewMembershipTypeBuilder().Invalid().Build()

		err := service.ChangeMembershipType(ctx, invalidMembershipType)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *MembershipTypeServiceSuite) TestDeleteMembershipType(t provider.T) {
	t.Title("[DeleteMembershipType] Successfully deletes membership type")
	t.Tags("membership_type", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted membership type", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		membershipTypeID := uuid.New()

		repoMock.On("DeleteMembershipType", ctx, membershipTypeID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "membershipTypeID", membershipTypeID)

		err := service.DeleteMembershipType(ctx, membershipTypeID)

		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeServiceSuite) TestGetMembershipTypeByID(t provider.T) {
	t.Title("[GetMembershipTypeByID] Successfully retrieves a membership type by ID")
	t.Tags("membership_type", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved membership type", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		membershipTypeID := uuid.New()
		expectedMembershipType := builder.NewMembershipTypeBuilder().Build()

		repoMock.On("GetMembershipTypeByID", ctx, membershipTypeID).Return(expectedMembershipType, nil)

		sCtx.WithNewParameters("ctx", ctx, "membershipTypeID", membershipTypeID)

		membershipType, err := service.GetMembershipTypeByID(ctx, membershipTypeID)

		sCtx.Assert().Equal(expectedMembershipType, membershipType)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: error while retrieving membership type", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		membershipTypeID := uuid.New()

		repoMock.On("GetMembershipTypeByID", ctx, membershipTypeID).Return(entity.MembershipType{}, errors.New("not found"))

		sCtx.WithNewParameters("ctx", ctx, "membershipTypeID", membershipTypeID)

		membershipType, err := service.GetMembershipTypeByID(ctx, membershipTypeID)

		sCtx.Assert().Equal(entity.MembershipType{}, membershipType)
		sCtx.Assert().Error(err)
	})
}

func (s *MembershipTypeServiceSuite) TestListMembershipTypesByGymID(t provider.T) {
	t.Title("[ListMembershipTypesByGymID] Successfully lists membership types by gym ID")
	t.Tags("membership_type", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed membership types", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		gymID := uuid.New()
		expectedMembershipTypes := []entity.MembershipType{
			builder.NewMembershipTypeBuilder().Build(),
			builder.NewMembershipTypeBuilder().Build(),
		}

		repoMock.On("ListMembershipTypesByGymID", ctx, gymID).Return(expectedMembershipTypes, nil)

		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		membershipTypes, err := service.ListMembershipTypesByGymID(ctx, gymID)

		sCtx.Assert().Equal(expectedMembershipTypes, membershipTypes)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: error while listing membership types", func(sCtx provider.StepCtx) {
		repoMock := &mocks.IMembershipTypeRepository{}
		service := &MembershipTypeService{repoMock}

		ctx := context.TODO()
		gymID := uuid.New()

		repoMock.On("ListMembershipTypesByGymID", ctx, gymID).Return([]entity.MembershipType{}, errors.New("no types found"))

		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		membershipTypes, err := service.ListMembershipTypesByGymID(ctx, gymID)

		sCtx.Assert().Equal([]entity.MembershipType{}, membershipTypes)
		sCtx.Assert().Error(err)
	})
}

func TestMembershipTypeServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(MembershipTypeServiceSuite))
}
