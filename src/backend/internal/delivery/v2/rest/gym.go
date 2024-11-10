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

type GymHandler struct {
	gymService service.IGymService
}

func NewGymHandler(gymService service.IGymService) *GymHandler {
	return &GymHandler{
		gymService: gymService,
	}
}

func (h *GymHandler) GetGymByID(ctx *gin.Context) {
	log.Print("GetGym:", ctx.Request)

	id := ctx.Param("gymId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	gym, err := h.gymService.GetGymByID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"gym": gym})
}

func (h *GymHandler) GetGyms(ctx *gin.Context) {
	log.Print("GetGyms request:", ctx.Request)

	gyms, err := h.gymService.ListGyms(ctx.Request.Context())
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"gyms": gyms})
}

func (h *GymHandler) PostGym(ctx *gin.Context) {
	log.Print("PostGym:", ctx.Request)

	var req dto.PostGym

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.gymService.RegisterNewGym(ctx.Request.Context(), entity.Gym{
		ID:      req.Id,
		Name:    req.Name,
		Phone:   req.Phone,
		City:    req.City,
		Addres:  req.Addres,
		IsChain: req.IsChain,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *GymHandler) ChangeGym(ctx *gin.Context) {
	log.Print("ChangeGym request: ", ctx.Request)

	var req dto.PutGym
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id := ctx.Param("gymId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.gymService.ChangeGym(ctx.Request.Context(), entity.Gym{
		ID:      uuID,
		Name:    req.Name,
		Phone:   req.Phone,
		City:    req.City,
		Addres:  req.Addres,
		IsChain: req.IsChain,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *GymHandler) DeleteGym(ctx *gin.Context) {
	log.Print("DeleteGym: ", ctx.Request)

	id := ctx.Param("gymId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.gymService.DeleteGym(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
