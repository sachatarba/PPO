package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v2/dto"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service"
)

type TrainingHandler struct {
	trainingService service.ITrainingService
}

func NewTrainingHandler(trainingService service.ITrainingService) *TrainingHandler {
	return &TrainingHandler{
		trainingService: trainingService,
	}
}

func (h *TrainingHandler) GetTrainingsByTrainerID(ctx *gin.Context) {
	log.Print("GetTrainingsByTrainerID:", ctx.Request)

	id := ctx.Param("trainerId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	trainings, err := h.trainingService.ListTrainingsByTrainerID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trainings": trainings})
}

func (h *TrainingHandler) PostTraining(ctx *gin.Context) {
	log.Print("PostTraining:", ctx.Request)

	var req dto.PostTraining

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id := ctx.Param("trainerId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainingService.CreateNewTraining(ctx.Request.Context(), entity.Training{
		ID:           req.Id,
		Title:        req.Title,
		Description:  req.Description,
		TrainingType: req.TrainingType,
		TrainerID:    uuID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *TrainingHandler) PutTraining(ctx *gin.Context) {
	log.Print("ChangeTraining: ", ctx.Request)

	var req dto.PutTraining
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id := ctx.Param("trainerId")

	trainerId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id = ctx.Param("trainingId")

	trainingId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainingService.ChangeTraining(ctx.Request.Context(), entity.Training{
		ID:           trainingId,
		Title:        req.Title,
		Description:  req.Description,
		TrainingType: req.TrainingType,
		TrainerID:    trainerId,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *TrainingHandler) DeleteTraining(ctx *gin.Context) {
	log.Print("DeleteTraining: ", ctx.Request)

	id := ctx.Param("trainingId")

	trainingId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainingService.DeleteTraining(ctx.Request.Context(), trainingId)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
