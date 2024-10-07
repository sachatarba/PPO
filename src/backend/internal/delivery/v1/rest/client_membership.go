package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListClientMembershipsByClientID(ctx *gin.Context) {
	log.Print("ListClientMembershipsByClientID:", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	memberships, err := h.clientMembershipService.ListClientMembershipsByClientID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"clientMemberships": memberships})
}

func (h *Handler) CreateNewClientMembership(ctx *gin.Context) {
	log.Print("CreateNewClientMembership:", ctx.Request)

	var req request.ClientMembershipReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.CreateNewClientMembership(ctx.Request.Context(), entity.ClientMembership{
		ID:        req.ID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		MembershipType: entity.MembershipType{
			ID: req.MembershipTypeID,
		},
		ClientID: req.ClientID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeClientMembership(ctx *gin.Context) {
	log.Print("ChangeClientMembership: ", ctx.Request)

	var req request.ClientMembershipReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.ChangeClientMembership(ctx.Request.Context(), entity.ClientMembership{
		ID:        req.ID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		MembershipType: entity.MembershipType{
			ID: req.MembershipTypeID,
		},
		ClientID: req.ClientID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteClientMembership(ctx *gin.Context) {
	log.Print("DeleteClientMembership request: ", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.DeleteClientMembership(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
