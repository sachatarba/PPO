package rest


import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v2/dto"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/entity"
)

type ClientMembershipHandler struct {
	clientMembershipService	service.IClientMembershipsService
}

func NewClientMembershipHandler(
	clientMembershipService service.IClientMembershipsService) (*ClientMembershipHandler) {
		return &ClientMembershipHandler{
			clientMembershipService: clientMembershipService,
		}
}

func (h *ClientMembershipHandler) GetClientMemberships(ctx *gin.Context) {
	log.Print("GetClientMemberships:", ctx.Request)

	id := ctx.Param("clientId")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	memberships, err := h.clientMembershipService.ListClientMembershipsByClientID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"clientMemberships": memberships})
}

func (h *ClientMembershipHandler) PostClientMembership(ctx *gin.Context) {
	log.Print("PostClientMembership:", ctx.Request)

	clientId := ctx.Param("clientId")

	clientUuid, err := uuid.Parse(clientId)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	var req dto.PostClientMembership

	err = ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.CreateNewClientMembership(ctx.Request.Context(), entity.ClientMembership{
		ID:        req.Id,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		MembershipType: entity.MembershipType{
			ID: req.MembershipTypeID,
		},
		ClientID: clientUuid,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *ClientMembershipHandler) PutClientMembership(ctx *gin.Context) {
	log.Print("PutClientMembership: ", ctx.Request)

	var req dto.PutClientMembership
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	clientId := ctx.Param("clientId")

	clientUuid, err := uuid.Parse(clientId)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	id := ctx.Param("clientMembershipId")

	uuId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.ChangeClientMembership(ctx.Request.Context(), entity.ClientMembership{
		ID:        uuId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		MembershipType: entity.MembershipType{
			ID: req.MembershipTypeID,
		},
		ClientID: clientUuid,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *ClientMembershipHandler) DeleteClientMembership(ctx *gin.Context) {
	log.Print("DeleteClientMembershipt: ", ctx.Request)

	id := ctx.Param("clientMembershipId")

	uuId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientMembershipService.DeleteClientMembership(ctx.Request.Context(), uuId)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
