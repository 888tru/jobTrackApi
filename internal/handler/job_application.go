package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/888tru/jobtrack-api/internal/model"
	"github.com/888tru/jobtrack-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type JobApplicationHandler struct {
	repo repository.JobApplicationRepository
}

func NewJobApplicationHandler(repo repository.JobApplicationRepository) *JobApplicationHandler {
	return &JobApplicationHandler{repo: repo}
}

type createAppRequest struct {
	Company   string                  `json:"company" binding:"required"`
	Position  string                  `json:"position" binding:"required"`
	Status    model.ApplicationStatus `json:"status"`
	AppliedAt string                  `json:"applied_at"`
	Notes     string                  `json:"notes"`
	URL       string                  `json:"url"`
}

type updateAppRequest struct {
	Company  string                  `json:"company"`
	Position string                  `json:"position"`
	Status   model.ApplicationStatus `json:"status"`
	Notes    string                  `json:"notes"`
	URL      string                  `json:"url"`
}

func (h *JobApplicationHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var req createAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := req.Status
	if status == "" {
		status = model.StatusApplied
	}
	appliedAt := time.Now()
	if req.AppliedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.AppliedAt); err == nil {
			appliedAt = t
		}
	}
	app := &model.JobApplication{
		UserID:    userID,
		Company:   req.Company,
		Position:  req.Position,
		Status:    status,
		AppliedAt: appliedAt,
		Notes:     req.Notes,
		URL:       req.URL,
	}
	if err := h.repo.Create(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create application"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": app})
}

func (h *JobApplicationHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	apps, err := h.repo.FindAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": apps})
}

func (h *JobApplicationHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := h.repo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func (h *JobApplicationHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := h.repo.FindByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	var req updateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Company != "" {
		app.Company = req.Company
	}
	if req.Position != "" {
		app.Position = req.Position
	}
	if req.Status != "" {
		app.Status = req.Status
	}
	app.Notes = req.Notes
	app.URL = req.URL

	if err := h.repo.Update(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func (h *JobApplicationHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(param), 10, 64)
	return uint(v), err
}
