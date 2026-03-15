package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apprep "github.com/ipnext/admin-backend/internal/application/reporte"
)

type ReporteHandlerImpl struct {
	financiero  *apprep.ReporteFinancieroUseCase
	nomina      *apprep.ReporteNominaUseCase
	proveedores *apprep.ReporteProveedoresUseCase
	inmuebles   *apprep.ReporteInmueblesUseCase
}

func NewReporteHandler(
	financiero *apprep.ReporteFinancieroUseCase,
	nomina *apprep.ReporteNominaUseCase,
	proveedores *apprep.ReporteProveedoresUseCase,
	inmuebles *apprep.ReporteInmueblesUseCase,
) *ReporteHandlerImpl {
	return &ReporteHandlerImpl{
		financiero: financiero, nomina: nomina,
		proveedores: proveedores, inmuebles: inmuebles,
	}
}

func (h *ReporteHandlerImpl) Financiero(c *gin.Context) {
	desde := time.Now().AddDate(0, -6, 0)
	hasta := time.Now()
	if d := c.Query("desde"); d != "" {
		if t, err := time.Parse("2006-01-02", d); err == nil {
			desde = t
		}
	}
	if hs := c.Query("hasta"); hs != "" {
		if t, err := time.Parse("2006-01-02", hs); err == nil {
			hasta = t
		}
	}
	items, err := h.financiero.Execute(c.Request.Context(), desde, hasta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ReporteHandlerImpl) Nomina(c *gin.Context) {
	periodo := c.Query("periodo")
	if periodo == "" {
		now := time.Now()
		periodo = now.Format("2006-01")
	}
	items, err := h.nomina.Execute(c.Request.Context(), periodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ReporteHandlerImpl) Proveedores(c *gin.Context) {
	items, err := h.proveedores.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ReporteHandlerImpl) Inmuebles(c *gin.Context) {
	items, err := h.inmuebles.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ReporteHandlerImpl) Exportar(c *gin.Context) {
	tipo := c.Query("tipo")
	desde := time.Now().AddDate(0, -6, 0)
	hasta := time.Now()

	switch tipo {
	case "financiero":
		items, err := h.financiero.Execute(c.Request.Context(), desde, hasta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		data, err := apprep.ToCSV(items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		c.Header("Content-Disposition", "attachment; filename=reporte_financiero.csv")
		c.Data(http.StatusOK, "text/csv", data)
	case "proveedores":
		items, err := h.proveedores.Execute(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		data, err := apprep.ToCSV(items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		c.Header("Content-Disposition", "attachment; filename=reporte_proveedores.csv")
		c.Data(http.StatusOK, "text/csv", data)
	case "inmuebles":
		items, err := h.inmuebles.Execute(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		data, err := apprep.ToCSV(items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errInterno())
			return
		}
		c.Header("Content-Disposition", "attachment; filename=reporte_inmuebles.csv")
		c.Data(http.StatusOK, "text/csv", data)
	default:
		c.JSON(http.StatusBadRequest, errValidacion("tipo debe ser: financiero, proveedores o inmuebles"))
	}
}
