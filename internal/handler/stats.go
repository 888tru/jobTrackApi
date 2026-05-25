package handler

import (
	"net/http"

	"github.com/888tru/jobtrack-api/internal/service"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	svc service.StatsService
}

func NewStatsHandler(svc service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	stats, err := h.svc.GetStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
