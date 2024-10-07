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

type MembershipTypeHandler struct {
	membershipTypeService service.IMembershipTypeService
}

func NewMembershipTypeHandler(membershipTypeService service.IMembershipTypeService) *MembershipTypeHandler {
	return &MembershipTypeHandler{
		membershipTypeService: membershipTypeService,
	}
}

func (h *MembershipTypeHandler) GetMembershipTypeByGymID(ctx *gin.Context) {
	log.Print("GetMembershipTypeByGymID:", ctx.Request)

	id := ctx.Param("gymId")

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

func (h *MembershipTypeHandler) PostMembershipType(ctx *gin.Context) {
	log.Print("PostMembershipType:", ctx.Request)

	var req dto.PostMembershipType

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

	err = h.membershipTypeService.RegisterNewMembershipType(ctx.Request.Context(), entity.MembershipType{
		ID:           req.Id,
		Type:         req.Type,
		Description:  req.Description,
		Price:        req.Price,
		DaysDuration: req.DaysDuration,
		GymID:        gymId,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *MembershipTypeHandler) PutMembershipType(ctx *gin.Context) {
	log.Print("ChangeMembershipType request: ", ctx.Request)

	var req dto.PutMembershipType
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

	id = ctx.Param("membershipTypeId")

	membershipTypeId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.membershipTypeService.ChangeMembershipType(ctx.Request.Context(), entity.MembershipType{
		ID:           membershipTypeId,
		Type:         req.Type,
		Price:        req.Price,
		DaysDuration: req.DaysDuration,
		GymID:        gymId,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *MembershipTypeHandler) DeleteMembershipType(ctx *gin.Context) {
	log.Print("DeleteMembershipType: ", ctx.Request)

	id := ctx.Param("membershipTypeId")

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
