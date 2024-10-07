package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListTrainingsByTrainerID(ctx *gin.Context) {
	log.Print("ListTrainingsByTrainerID:", ctx.Request)

	id := ctx.Param("id")

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

func (h *Handler) CreateNewTraining(ctx *gin.Context) {
	log.Print("CreateNewTraining:", ctx.Request)

	var req request.TrainingReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.trainingService.CreateNewTraining(ctx.Request.Context(), entity.Training{
		ID:           req.ID,
		Title:        req.Title,
		Description:  req.Description,
		TrainingType: req.TrainingType,
		TrainerID:    req.TrainerID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeTraining(ctx *gin.Context) {
	log.Print("ChangeTraining: ", ctx.Request)

	var req request.TrainingReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.trainingService.ChangeTraining(ctx.Request.Context(), entity.Training{
		ID:           req.ID,
		Title:        req.Title,
		Description:  req.Description,
		TrainingType: req.TrainingType,
		TrainerID:    req.TrainerID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteTraining(ctx *gin.Context) {
	log.Print("DeleteTraining request: ", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainingService.DeleteTraining(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
