package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	apptrans "github.com/ipnext/admin-backend/internal/application/transferencia"
	"github.com/ipnext/admin-backend/internal/domain/transferencia"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/middleware"
)

type TransferenciaHandlerImpl struct {
	list         *apptrans.ListUseCase
	get          *apptrans.GetUseCase
	create       *apptrans.CreateUseCase
	update       *apptrans.UpdateUseCase
	delete       *apptrans.DeleteUseCase
	cambiarEstado *apptrans.CambiarEstadoUseCase
	calendario   *apptrans.CalendarioUseCase
	recurrentes  *apptrans.RecurrentesUseCase
}

func NewTransferenciaHandler(
	list *apptrans.ListUseCase,
	get *apptrans.GetUseCase,
	create *apptrans.CreateUseCase,
	update *apptrans.UpdateUseCase,
	del *apptrans.DeleteUseCase,
	cambiarEstado *apptrans.CambiarEstadoUseCase,
	calendario *apptrans.CalendarioUseCase,
	recurrentes *apptrans.RecurrentesUseCase,
) *TransferenciaHandlerImpl {
	return &TransferenciaHandlerImpl{
		list: list, get: get, create: create, update: update,
		delete: del, cambiarEstado: cambiarEstado,
		calendario: calendario, recurrentes: recurrentes,
	}
}

func (h *TransferenciaHandlerImpl) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := h.list.Execute(c.Request.Context(), transferencia.Filtros{
		Estado:    c.Query("estado"),
		Categoria: c.Query("categoria"),
		Query:     c.Query("q"),
		Orden:     c.Query("orden"),
		Dir:       c.Query("dir"),
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result.Items,
		"meta": gin.H{"total": result.Total, "page": result.Page, "per_page": result.PerPage, "total_pages": result.TotalPages},
	})
}

func (h *TransferenciaHandlerImpl) Get(c *gin.Context) {
	t, err := h.get.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, apptrans.ErrNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Transferencia no encontrada"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

type createTransferenciaRequest struct {
	Beneficiario string  `json:"beneficiario" binding:"required"`
	CBU          *string `json:"cbu"`
	Alias        *string `json:"alias"`
	Categoria    string  `json:"categoria" binding:"required"`
	Monto        float64 `json:"monto" binding:"required,gt=0"`
	Moneda       string  `json:"moneda" binding:"required,oneof=ARS USD"`
	FechaPago    string  `json:"fechaPago" binding:"required"`
	Frecuencia   string  `json:"frecuencia" binding:"required"`
	MetodoPago   string  `json:"metodoPago" binding:"required"`
	Notas        *string `json:"notas"`
	ProveedorID  *string `json:"proveedorId"`
}

func (h *TransferenciaHandlerImpl) Create(c *gin.Context) {
	var req createTransferenciaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	fecha, err := time.Parse("2006-01-02", req.FechaPago)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("fechaPago debe ser YYYY-MM-DD"))
		return
	}
	claims := middleware.GetClaims(c)
	t, err := h.create.Execute(c.Request.Context(), apptrans.CreateRequest{
		Beneficiario: req.Beneficiario,
		CBU:          req.CBU,
		Alias:        req.Alias,
		Categoria:    req.Categoria,
		Monto:        req.Monto,
		Moneda:       req.Moneda,
		FechaPago:    fecha,
		Frecuencia:   transferencia.Frecuencia(req.Frecuencia),
		MetodoPago:   transferencia.MetodoPago(req.MetodoPago),
		Notas:        req.Notas,
		ProveedorID:  req.ProveedorID,
		CreadoPor:    claims.Sub,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": t})
}

func (h *TransferenciaHandlerImpl) Update(c *gin.Context) {
	var body struct {
		Beneficiario     *string  `json:"beneficiario"`
		CBU              *string  `json:"cbu"`
		Alias            *string  `json:"alias"`
		Categoria        *string  `json:"categoria"`
		Monto            *float64 `json:"monto"`
		Moneda           *string  `json:"moneda"`
		FechaPago        *string  `json:"fechaPago"`
		FechaVencimiento *string  `json:"fechaVencimiento"`
		Frecuencia       *string  `json:"frecuencia"`
		MetodoPago       *string  `json:"metodoPago"`
		Notas            *string  `json:"notas"`
		ProveedorID      *string  `json:"proveedorId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	req := apptrans.UpdateRequest{
		Beneficiario: body.Beneficiario, CBU: body.CBU, Alias: body.Alias,
		Categoria: body.Categoria, Monto: body.Monto, Moneda: body.Moneda,
		Frecuencia: body.Frecuencia, MetodoPago: body.MetodoPago,
		Notas: body.Notas, ProveedorID: body.ProveedorID,
	}
	if body.FechaPago != nil {
		t, err := time.Parse("2006-01-02", *body.FechaPago)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, errValidacion("fechaPago debe ser YYYY-MM-DD"))
			return
		}
		req.FechaPago = &t
	}
	if body.FechaVencimiento != nil {
		t, err := time.Parse("2006-01-02", *body.FechaVencimiento)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, errValidacion("fechaVencimiento debe ser YYYY-MM-DD"))
			return
		}
		req.FechaVencimiento = &t
	}
	t, err := h.update.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, apptrans.ErrNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Transferencia no encontrada"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func (h *TransferenciaHandlerImpl) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, apptrans.ErrNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Transferencia no encontrada"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Eliminada correctamente"}})
}

func (h *TransferenciaHandlerImpl) CambiarEstado(c *gin.Context) {
	var body struct {
		Estado string `json:"estado" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	t, err := h.cambiarEstado.Execute(c.Request.Context(), c.Param("id"), transferencia.Estado(body.Estado))
	if err != nil {
		if errors.Is(err, apptrans.ErrNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Transferencia no encontrada"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func (h *TransferenciaHandlerImpl) Calendario(c *gin.Context) {
	desde := time.Now()
	hasta := desde.AddDate(0, 1, 0)
	items, err := h.calendario.Execute(c.Request.Context(), desde, hasta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *TransferenciaHandlerImpl) Recurrentes(c *gin.Context) {
	items, err := h.recurrentes.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// helpers de error compartidos
func errInterno() gin.H {
	return gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Error interno del servidor"}}
}
func errNotFound(msg string) gin.H {
	return gin.H{"error": gin.H{"code": "NOT_FOUND", "message": msg}}
}
func errValidacion(msg string) gin.H {
	return gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": msg}}
}
