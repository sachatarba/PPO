package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListEquipmentsByGymID(ctx *gin.Context) {
	log.Print("GetListEquipmentsByGymID:", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	equipments, err := h.equipmentService.ListEquipmentsByGymID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"equipments": equipments})
}

func (h *Handler) CreateNewEquipment(ctx *gin.Context) {
	log.Print("CreateNewEquipment:", ctx.Request)

	var req request.EquipmentReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.equipmentService.CreateNewEquipment(ctx.Request.Context(), entity.Equipment{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		GymID:       req.GymID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeEquipment(ctx *gin.Context) {
	log.Print("ChangeEquipment request: ", ctx.Request)

	var req request.EquipmentReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.equipmentService.ChangeEquipment(ctx.Request.Context(), entity.Equipment{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		GymID:       req.GymID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteEquipment(ctx *gin.Context) {
	log.Print("DeleteEquipment request: ", ctx.Request)

	id := ctx.Param("id")

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
