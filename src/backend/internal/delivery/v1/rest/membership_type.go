package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListMembershipTypeByGymID(ctx *gin.Context) {
	log.Print("ListMembershipTypeByGymID:", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	membershipTypes, err := h.membershipTypeService.ListMembershipTypesByGymID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"membershipTypes": membershipTypes})
}

func (h *Handler) CreateNewMembershipType(ctx *gin.Context) {
	log.Print("CreateNewMembershipType:", ctx.Request)

	var req request.MembershipTypeReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.membershipTypeService.RegisterNewMembershipType(ctx.Request.Context(), entity.MembershipType{
		ID:           req.ID,
		Type:         req.Type,
		Description:  req.Description,
		Price:        req.Price,
		DaysDuration: req.DaysDuration,
		GymID:        req.GymID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeMembershipType(ctx *gin.Context) {
	log.Print("ChangeMembershipType request: ", ctx.Request)

	var req request.MembershipTypeReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.membershipTypeService.ChangeMembershipType(ctx.Request.Context(), entity.MembershipType{
		ID:           req.ID,
		Type:         req.Type,
		Price:        req.Price,
		DaysDuration: req.DaysDuration,
		GymID:        req.GymID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteMembershipType(ctx *gin.Context) {
	log.Print("DeleteMembershipType request: ", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.membershipTypeService.DeleteMembershipType(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
