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

type EquipmentHandler struct {
	equipmentService service.IEquipmentService
}

func NewEquipmentHandler(equipmentService service.IEquipmentService) *EquipmentHandler {
	return &EquipmentHandler{
		equipmentService: equipmentService,
	}
}

func (h *EquipmentHandler) GetEquipments(ctx *gin.Context) {
	log.Print("GetEquipment:", ctx.Request)

	id := ctx.Param("gymId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	equipments, err := h.equipmentService.ListEquipmentsByGymID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"equipments": equipments})
}

func (h *EquipmentHandler) PostEquipment(ctx *gin.Context) {
	log.Print("PostEquipment:", ctx.Request)

	var req dto.PostEquipment

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id := ctx.Param("gymId")

	gymId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.equipmentService.CreateNewEquipment(ctx.Request.Context(), entity.Equipment{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		GymID:       gymId,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *EquipmentHandler) PutEquipment(ctx *gin.Context) {
	log.Print("PutEquipment: ", ctx.Request)

	var req dto.PutEquipment
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id := ctx.Param("gymId")

	gymId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id = ctx.Param("equipmentId")

	equipmentId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.equipmentService.ChangeEquipment(ctx.Request.Context(), entity.Equipment{
		ID:          equipmentId,
		Name:        req.Name,
		Description: req.Description,
		GymID:       gymId,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *EquipmentHandler) DeleteEquipment(ctx *gin.Context) {
	log.Print("DeleteEquipment: ", ctx.Request)

	id := ctx.Param("equipmentId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.equipmentService.DeleteEquipment(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
