package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipnext/admin-backend/internal/application/dashboard"
)

type DashboardHandlerImpl struct {
	getKPIs               *dashboard.GetKPIsUseCase
	getPagosUrgentes      *dashboard.GetPagosUrgentesUseCase
	getDistribucionEgresos *dashboard.GetDistribucionEgresosUseCase
}

func NewDashboardHandler(
	getKPIs *dashboard.GetKPIsUseCase,
	getPagosUrgentes *dashboard.GetPagosUrgentesUseCase,
	getDistribucionEgresos *dashboard.GetDistribucionEgresosUseCase,
) *DashboardHandlerImpl {
	return &DashboardHandlerImpl{
		getKPIs:               getKPIs,
		getPagosUrgentes:      getPagosUrgentes,
		getDistribucionEgresos: getDistribucionEgresos,
	}
}

func (h *DashboardHandlerImpl) GetKPIs(c *gin.Context) {
	kpis, err := h.getKPIs.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

func (h *DashboardHandlerImpl) GetPagosUrgentes(c *gin.Context) {
	items, err := h.getPagosUrgentes.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *DashboardHandlerImpl) GetDistribucionEgresos(c *gin.Context) {
	items, err := h.getDistribucionEgresos.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *DashboardHandlerImpl) GetActividadReciente(c *gin.Context) {
	// Actividad reciente: placeholder para v1 — retorna lista vacía
	c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
}
