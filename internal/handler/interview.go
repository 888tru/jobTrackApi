package handler

import (
	"net/http"
	"time"

	"github.com/888tru/jobtrack-api/internal/model"
	"github.com/888tru/jobtrack-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type InterviewHandler struct {
	repo    repository.InterviewRepository
	appRepo repository.JobApplicationRepository
}

func NewInterviewHandler(repo repository.InterviewRepository, appRepo repository.JobApplicationRepository) *InterviewHandler {
	return &InterviewHandler{repo: repo, appRepo: appRepo}
}

type createInterviewRequest struct {
	ScheduledAt string              `json:"scheduled_at" binding:"required"`
	Type        model.InterviewType `json:"type"`
	Notes       string              `json:"notes"`
}

type updateInterviewRequest struct {
	ScheduledAt string              `json:"scheduled_at"`
	Type        model.InterviewType `json:"type"`
	Notes       string              `json:"notes"`
	Result      string              `json:"result"`
}

func (h *InterviewHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	appID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	if _, err := h.appRepo.FindByID(appID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	interviews, err := h.repo.FindByApplicationID(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch interviews"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": interviews})
}

func (h *InterviewHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	appID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	if _, err := h.appRepo.FindByID(appID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	var req createInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_at format, use RFC3339"})
		return
	}
	interviewType := req.Type
	if interviewType == "" {
		interviewType = model.InterviewOther
	}
	interview := &model.Interview{
		JobApplicationID: appID,
		ScheduledAt:      scheduledAt,
		Type:             interviewType,
		Notes:            req.Notes,
	}
	if err := h.repo.Create(interview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create interview"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": interview})
}

func (h *InterviewHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	appID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	if _, err := h.appRepo.FindByID(appID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	interviewID, err := parseUintParam(c, "interview_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}
	interview, err := h.repo.FindByID(interviewID, appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "interview not found"})
		return
	}
	var req updateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ScheduledAt != "" {
		if t, err := time.Parse(time.RFC3339, req.ScheduledAt); err == nil {
			interview.ScheduledAt = t
		}
	}
	if req.Type != "" {
		interview.Type = req.Type
	}
	interview.Notes = req.Notes
	interview.Result = req.Result

	if err := h.repo.Update(interview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update interview"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": interview})
}

func (h *InterviewHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	appID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	if _, err := h.appRepo.FindByID(appID, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	interviewID, err := parseUintParam(c, "interview_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interview id"})
		return
	}
	if err := h.repo.Delete(interviewID, appID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete interview"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
