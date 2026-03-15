package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	appalq "github.com/ipnext/admin-backend/internal/application/alquiler"
)

type AlquilerHandlerImpl struct {
	listInmuebles   *appalq.ListInmueblesUseCase
	getInmueble     *appalq.GetInmuebleUseCase
	createInmueble  *appalq.CreateInmuebleUseCase
	updateInmueble  *appalq.UpdateInmuebleUseCase
	deleteInmueble  *appalq.DeleteInmuebleUseCase
	listContratos   *appalq.ListContratosUseCase
	createContrato  *appalq.CreateContratoUseCase
	vencimientos    *appalq.VencimientosUseCase
	listPagos       *appalq.ListPagosUseCase
	createPago      *appalq.CreatePagoUseCase
}

func NewAlquilerHandler(
	listInmuebles *appalq.ListInmueblesUseCase,
	getInmueble *appalq.GetInmuebleUseCase,
	createInmueble *appalq.CreateInmuebleUseCase,
	updateInmueble *appalq.UpdateInmuebleUseCase,
	deleteInmueble *appalq.DeleteInmuebleUseCase,
	listContratos *appalq.ListContratosUseCase,
	createContrato *appalq.CreateContratoUseCase,
	vencimientos *appalq.VencimientosUseCase,
	listPagos *appalq.ListPagosUseCase,
	createPago *appalq.CreatePagoUseCase,
) *AlquilerHandlerImpl {
	return &AlquilerHandlerImpl{
		listInmuebles: listInmuebles, getInmueble: getInmueble,
		createInmueble: createInmueble, updateInmueble: updateInmueble,
		deleteInmueble: deleteInmueble, listContratos: listContratos,
		createContrato: createContrato, vencimientos: vencimientos,
		listPagos: listPagos, createPago: createPago,
	}
}

func (h *AlquilerHandlerImpl) List(c *gin.Context) {
	items, err := h.listInmuebles.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *AlquilerHandlerImpl) Get(c *gin.Context) {
	i, err := h.getInmueble.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, appalq.ErrInmuebleNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Inmueble no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": i})
}

func (h *AlquilerHandlerImpl) Create(c *gin.Context) {
	var req appalq.CreateInmuebleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	i, err := h.createInmueble.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": i})
}

func (h *AlquilerHandlerImpl) Update(c *gin.Context) {
	var req appalq.UpdateInmuebleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	i, err := h.updateInmueble.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, appalq.ErrInmuebleNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Inmueble no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": i})
}

func (h *AlquilerHandlerImpl) Delete(c *gin.Context) {
	if err := h.deleteInmueble.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, appalq.ErrInmuebleNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Inmueble no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Inmueble eliminado"}})
}

func (h *AlquilerHandlerImpl) ListContratos(c *gin.Context) {
	var inmuebleID *string
	if id := c.Query("inmueble_id"); id != "" {
		inmuebleID = &id
	}
	contratos, err := h.listContratos.Execute(c.Request.Context(), inmuebleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": contratos})
}

func (h *AlquilerHandlerImpl) CreateContrato(c *gin.Context) {
	var body struct {
		InmuebleID       string  `json:"inmuebleId" binding:"required"`
		VigenciaDesde    string  `json:"vigenciaDesde" binding:"required"`
		VigenciaHasta    string  `json:"vigenciaHasta" binding:"required"`
		AjusteFrecuencia string  `json:"ajusteFrecuencia" binding:"required"`
		MontoMensual     float64 `json:"montoMensual" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	desde, err := time.Parse("2006-01-02", body.VigenciaDesde)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("vigenciaDesde debe ser YYYY-MM-DD"))
		return
	}
	hasta, err := time.Parse("2006-01-02", body.VigenciaHasta)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("vigenciaHasta debe ser YYYY-MM-DD"))
		return
	}
	contrato, err := h.createContrato.Execute(c.Request.Context(), appalq.CreateContratoRequest{
		InmuebleID: body.InmuebleID, VigenciaDesde: desde, VigenciaHasta: hasta,
		AjusteFrecuencia: body.AjusteFrecuencia, MontoMensual: body.MontoMensual,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": contrato})
}

func (h *AlquilerHandlerImpl) ListPagos(c *gin.Context) {
	var inmuebleID *string
	if id := c.Query("inmueble_id"); id != "" {
		inmuebleID = &id
	}
	pagos, err := h.listPagos.Execute(c.Request.Context(), inmuebleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pagos})
}

func (h *AlquilerHandlerImpl) CreatePago(c *gin.Context) {
	var req appalq.CreatePagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	pago, err := h.createPago.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": pago})
}

func (h *AlquilerHandlerImpl) GetVencimientos(c *gin.Context) {
	dias := 30
	if d := c.Query("dias"); d != "" {
		if v, err := strconv.Atoi(d); err == nil {
			dias = v
		}
	}
	contratos, err := h.vencimientos.Execute(c.Request.Context(), dias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": contratos})
}
