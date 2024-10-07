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

type ScheduleHandler struct {
	scheduleService service.IScheduleService
}

func NewScheduleHandler(scheduleService service.IScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) GetSchedules(ctx *gin.Context) {
	log.Print("GetSchedules:", ctx.Request)

	id := ctx.Param("clientId")

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

func (h *ScheduleHandler) PostSchedule(ctx *gin.Context) {
	log.Print("PostSchedule:", ctx.Request)

	var req dto.PostSchedule

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id := ctx.Param("clientId")

	clientId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.CreateNewSchedule(ctx.Request.Context(), entity.Schedule{
		ID:           req.Id,
		DayOfTheWeek: req.DayOfTheWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ClientID:     clientId,
		TrainingID:   req.TrainingID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *ScheduleHandler) PutSchedule(ctx *gin.Context) {
	log.Print("PutSchedule: ", ctx.Request)

	var req dto.PutSchedule
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id := ctx.Param("clientId")

	clientId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id = ctx.Param("scheduleId")

	scheduleId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.ChangeSchedule(ctx.Request.Context(), entity.Schedule{
		ID:           scheduleId,
		DayOfTheWeek: req.DayOfTheWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		ClientID:     clientId,
		TrainingID:   req.TrainingID,
	})
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *ScheduleHandler) DeleteSchedule(ctx *gin.Context) {
	log.Print("DeleteSchedule: ", ctx.Request)

	id := ctx.Param("scheduleId")

	scheduleId, err := uuid.Parse(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.scheduleService.DeleteSchedule(ctx.Request.Context(), scheduleId)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
