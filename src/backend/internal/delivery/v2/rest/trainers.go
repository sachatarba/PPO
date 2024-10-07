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

type TrainerHandler struct {
	trainerService service.ITrainerService
}

func NewTrainerHandler(trainerService service.ITrainerService) *TrainerHandler {
	return &TrainerHandler{
		trainerService: trainerService,
	}
}

func (h *TrainerHandler) GetTrainers(ctx *gin.Context) {
	log.Print("GetTrainer:", ctx.Request)

	trainers, err := h.trainerService.ListTrainers(ctx.Request.Context())
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trainers": trainers})
}

func (h *TrainerHandler) GetTrainersByGymID(ctx *gin.Context) {
	log.Print("GetTrainersByGymID:", ctx.Request)

	id := ctx.Param("gymId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	trainers, err := h.trainerService.ListTrainersByGymID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trainers": trainers})
}

func (h *TrainerHandler) PostTrainer(ctx *gin.Context) {
	log.Print("PostTrainer:", ctx.Request)

	var req dto.PostTrainer

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainerService.RegisterNewTrainer(ctx.Request.Context(), entity.Trainer{
		ID:            req.Id,
		Fullname:      req.Fullname,
		Email:         req.Email,
		Phone:         req.Phone,
		Qualification: req.Qualification,
		UnitPrice:     req.UnitPrice,
		GymsID:        req.GymsID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *TrainerHandler) PutTrainer(ctx *gin.Context) {
	log.Print("PutTrainer: ", ctx.Request)

	var req dto.PutTrainer
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

	err = h.trainerService.ChangeTrainer(ctx.Request.Context(), entity.Trainer{
		ID:            uuID,
		Fullname:      req.Fullname,
		Email:         req.Email,
		Phone:         req.Phone,
		Qualification: req.Qualification,
		UnitPrice:     req.UnitPrice,
		GymsID:        req.GymsID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *TrainerHandler) DeleteTrainer(ctx *gin.Context) {
	log.Print("DeleteTrainer: ", ctx.Request)

	id := ctx.Param("trainerId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.trainerService.DeleteTrainer(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
