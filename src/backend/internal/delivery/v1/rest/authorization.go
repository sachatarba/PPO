package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) IsAuthorize(ctx *gin.Context) {
	log.Print("IsAuthorize request:", ctx.Request)

	sessionID, err := ctx.Cookie("session")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	session, err := h.authorizationService.IsAuthorize(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if session == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"session": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *Handler) Logout(ctx *gin.Context) {
	log.Print("Logout request:", ctx.Request)

	session, err := ctx.Cookie("session")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
	}

	ctx.SetCookie(
		"session",
		session,
		0,
		"/",
		"",
		true,
		true,
	)
	ctx.Status(http.StatusOK)
}

func (h *Handler) Login(ctx *gin.Context) {
	log.Print("Login request:", ctx.Request)

	var req request.LoginReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	session, err := h.authorizationService.Authorize(ctx.Request.Context(), req.Login, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	cookie := &http.Cookie{
		Name:    "session",
		Value:   session.SessionID.String(),
		Path:    "/",
		Domain:  "localhost", 
		Expires: session.TTL, 
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, 
	}

	http.SetCookie(ctx.Writer, cookie)
	ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *Handler) RegisterNewUser(ctx *gin.Context) {
	log.Print("RegisterNewUser request:", ctx.Request)

	var req request.ClientReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})

		return
	}

	session, err := h.authorizationService.Register(ctx.Request.Context(), entity.Client{
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

	ctx.SetCookie(
		"session",
		session.SessionID.String(),
		session.TTL.Second(),
		"/",
		"",
		true,
		true,
	)
	ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{"session": session})
}
