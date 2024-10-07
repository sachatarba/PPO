package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListClients(ctx *gin.Context) {
	log.Print("ListClients request:", ctx.Request)

	clients, err := h.clientService.ListClients(ctx.Request.Context())
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"clients": clients})
}

func (h *Handler) GetClientByLogin(ctx *gin.Context) {
	log.Print("GetClientByLogin request:", ctx.Request)

	login := ctx.Param("login")

	client, err := h.clientService.GetClientByLogin(ctx.Request.Context(), login)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"client": client})
}

func (h *Handler) GetClientByID(ctx *gin.Context) {
	log.Print("GetClientByID request:", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	client, err := h.clientService.GetClientByID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"client": client})
}

func (h *Handler) ChangeClient(ctx *gin.Context) {
	log.Print("ChangeClient request: ", ctx.Request)

	var req request.ClientReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientService.ChangeClient(ctx.Request.Context(), entity.Client{
		ID:        req.ID,
		Login:     req.Login,
		Password:  req.Password,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Phone:     req.Phone,
		Birthdate: req.Birthdate,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteClient(ctx *gin.Context) {
	log.Print("DeleteClient request: ", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})

		return
	}

	err = h.clientService.DeleteClient(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	ctx.Status(http.StatusOK)
}
