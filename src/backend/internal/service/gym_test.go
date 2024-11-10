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

type GymServiceSuite struct {
	suite.Suite
}

func (s *GymServiceSuite) TestRegisterNewGym(t provider.T) {
	t.Title("[RegisterNewGym] Successfully registers gym")
	t.Tags("gym", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered gym", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()

		gymRepoMock.On("RegisterNewGym", mock.Anything, gym).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		err := gymService.RegisterNewGym(ctx, gym)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		invalidGym := builder.NewGymBuilder().Invalid().Build()

		err := gymService.RegisterNewGym(ctx, invalidGym)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *GymServiceSuite) TestChangeGym(t provider.T) {
	t.Title("[ChangeGym] Successfully changes gym")
	t.Tags("gym", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed gym", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()

		gymRepoMock.On("ChangeGym", mock.Anything, gym).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		err := gymService.ChangeGym(ctx, gym)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		invalidGym := builder.NewGymBuilder().Invalid().Build()

		err := gymService.ChangeGym(ctx, invalidGym)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *GymServiceSuite) TestDeleteGym(t provider.T) {
	t.Title("[DeleteGym] Successfully deletes gym")
	t.Tags("gym", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted gym", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		gymID := uuid.New()

		gymRepoMock.On("DeleteGym", mock.Anything, gymID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		err := gymService.DeleteGym(ctx, gymID)

		sCtx.Assert().NoError(err)
	})
}

func (s *GymServiceSuite) TestGetGymByID(t provider.T) {
	t.Title("[GetGymByID] Successfully retrieves gym by ID")
	t.Tags("gym", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved gym", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		gymID := uuid.New()
		expectedGym := builder.NewGymBuilder().Build()

		gymRepoMock.On("GetGymByID", mock.Anything, gymID).Return(expectedGym, nil)

		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		gym, err := gymService.GetGymByID(ctx, gymID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedGym, gym)
	})
}

func (s *GymServiceSuite) TestListGyms(t provider.T) {
	t.Title("[ListGyms] Successfully lists gyms")
	t.Tags("gym", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed gyms", func(sCtx provider.StepCtx) {
		gymRepoMock := &mocks.IGymRepository{}
		gymService := &GymService{gymRepoMock}

		ctx := context.TODO()
		expectedGyms := []entity.Gym{
			builder.NewGymBuilder().Build(),
			builder.NewGymBuilder().Build(),
		}

		gymRepoMock.On("ListGyms", mock.Anything).Return(expectedGyms, nil)

		sCtx.WithNewParameters("ctx", ctx)

		gyms, err := gymService.ListGyms(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedGyms, gyms)
	})
}

func TestGymServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(GymServiceSuite))
}
