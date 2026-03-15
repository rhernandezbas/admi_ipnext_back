package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	apptes "github.com/ipnext/admin-backend/internal/application/tesoreria"
)

type TesoreriaHandlerImpl struct {
	flujoCaja          *apptes.GetFlujoCajaUseCase
	proyecciones       *apptes.GetProyeccionesUseCase
	listCuentas        *apptes.ListCuentasUseCase
	createCuenta       *apptes.CreateCuentaUseCase
	updateCuenta       *apptes.UpdateCuentaUseCase
	getConciliacion    *apptes.GetConciliacionUseCase
	createMovimiento   *apptes.CreateMovimientoUseCase
	conciliarMovimiento *apptes.ConciliarMovimientoUseCase
}

func NewTesoreriaHandler(
	flujoCaja *apptes.GetFlujoCajaUseCase,
	proyecciones *apptes.GetProyeccionesUseCase,
	listCuentas *apptes.ListCuentasUseCase,
	createCuenta *apptes.CreateCuentaUseCase,
	updateCuenta *apptes.UpdateCuentaUseCase,
	getConciliacion *apptes.GetConciliacionUseCase,
	createMovimiento *apptes.CreateMovimientoUseCase,
	conciliarMovimiento *apptes.ConciliarMovimientoUseCase,
) *TesoreriaHandlerImpl {
	return &TesoreriaHandlerImpl{
		flujoCaja: flujoCaja, proyecciones: proyecciones,
		listCuentas: listCuentas, createCuenta: createCuenta,
		updateCuenta: updateCuenta, getConciliacion: getConciliacion,
		createMovimiento: createMovimiento, conciliarMovimiento: conciliarMovimiento,
	}
}

func (h *TesoreriaHandlerImpl) GetFlujoCaja(c *gin.Context) {
	desde := time.Now().AddDate(0, -1, 0)
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
	items, err := h.flujoCaja.Execute(c.Request.Context(), desde, hasta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *TesoreriaHandlerImpl) GetProyecciones(c *gin.Context) {
	meses := 6
	if m := c.Query("meses"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			meses = v
		}
	}
	items, err := h.proyecciones.Execute(c.Request.Context(), meses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *TesoreriaHandlerImpl) ListCuentas(c *gin.Context) {
	cuentas, err := h.listCuentas.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cuentas})
}

func (h *TesoreriaHandlerImpl) CreateCuenta(c *gin.Context) {
	var req apptes.CreateCuentaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	cuenta, err := h.createCuenta.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cuenta})
}

func (h *TesoreriaHandlerImpl) UpdateCuenta(c *gin.Context) {
	var req apptes.UpdateCuentaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	cuenta, err := h.updateCuenta.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, apptes.ErrCuentaNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Cuenta no encontrada"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cuenta})
}

func (h *TesoreriaHandlerImpl) GetConciliacion(c *gin.Context) {
	var cuentaID *string
	if id := c.Query("cuenta_id"); id != "" {
		cuentaID = &id
	}
	movs, err := h.getConciliacion.Execute(c.Request.Context(), cuentaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": movs})
}

func (h *TesoreriaHandlerImpl) CreateMovimiento(c *gin.Context) {
	var req apptes.CreateMovimientoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	mov, err := h.createMovimiento.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": mov})
}

func (h *TesoreriaHandlerImpl) ConciliarMovimiento(c *gin.Context) {
	if err := h.conciliarMovimiento.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Movimiento conciliado"}})
}
