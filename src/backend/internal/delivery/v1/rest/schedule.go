package rest

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"
)

func (h *Handler) ListSchedulesByClientID(ctx *gin.Context) {
	log.Print("ListSchedulesByClientID:", ctx.Request)

	id := ctx.Param("id")

	uuID, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	schedules, err := h.scheduleService.ListSchedulesByClientID(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func (h *Handler) CreateNewSchedule(ctx *gin.Context) {
	log.Print("CreateNewSchedule:", ctx.Request)

	var req request.ScheduleReq

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.CreateNewSchedule(ctx.Request.Context(), entity.Schedule{
		ID:           req.ID,
		DayOfTheWeek: req.DayOfTheWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ClientID:     req.ClientID,
		TrainingID:   req.TrainingID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeSchedule(ctx *gin.Context) {
	log.Print("ChangeSchedule: ", ctx.Request)

	var req request.ScheduleReq
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.ChangeSchedule(ctx.Request.Context(), entity.Schedule{
		ID:           req.ID,
		DayOfTheWeek: req.DayOfTheWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ClientID:     req.ClientID,
		TrainingID:   req.TrainingID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteSchedule(ctx *gin.Context) {
	log.Print("DeleteSchedule request: ", ctx.Request)

	id, ok := ctx.Keys["id"]
	if !ok {
		log.Print()
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "No key in request"})
		return
	}

	uuID, err := uuid.Parse(id.(string))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.DeleteSchedule(ctx.Request.Context(), uuID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
