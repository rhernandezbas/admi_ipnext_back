package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appprov "github.com/ipnext/admin-backend/internal/application/proveedor"
	"github.com/ipnext/admin-backend/internal/domain/proveedor"
)

type ProveedorHandlerImpl struct {
	list           *appprov.ListUseCase
	get            *appprov.GetUseCase
	create         *appprov.CreateUseCase
	update         *appprov.UpdateUseCase
	delete         *appprov.DeleteUseCase
	listContratos  *appprov.ListContratosUseCase
	createContrato *appprov.CreateContratoUseCase
	updateContrato *appprov.UpdateContratoUseCase
	getRanking     *appprov.GetRankingUseCase
}

func NewProveedorHandler(
	list *appprov.ListUseCase,
	get *appprov.GetUseCase,
	create *appprov.CreateUseCase,
	update *appprov.UpdateUseCase,
	del *appprov.DeleteUseCase,
	listContratos *appprov.ListContratosUseCase,
	createContrato *appprov.CreateContratoUseCase,
	updateContrato *appprov.UpdateContratoUseCase,
	getRanking *appprov.GetRankingUseCase,
) *ProveedorHandlerImpl {
	return &ProveedorHandlerImpl{
		list: list, get: get, create: create, update: update, delete: del,
		listContratos: listContratos, createContrato: createContrato,
		updateContrato: updateContrato, getRanking: getRanking,
	}
}

func (h *ProveedorHandlerImpl) List(c *gin.Context) {
	items, err := h.list.Execute(c.Request.Context(), c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ProveedorHandlerImpl) Get(c *gin.Context) {
	p, err := h.get.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, appprov.ErrProveedorNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Proveedor no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProveedorHandlerImpl) Create(c *gin.Context) {
	var req appprov.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	p, err := h.create.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *ProveedorHandlerImpl) Update(c *gin.Context) {
	var req appprov.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	p, err := h.update.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, appprov.ErrProveedorNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Proveedor no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProveedorHandlerImpl) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, appprov.ErrProveedorNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Proveedor no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Proveedor eliminado"}})
}

func (h *ProveedorHandlerImpl) ListContratos(c *gin.Context) {
	var provID *string
	if id := c.Query("proveedor_id"); id != "" {
		provID = &id
	}
	items, err := h.listContratos.Execute(c.Request.Context(), provID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ProveedorHandlerImpl) CreateContrato(c *gin.Context) {
	var body struct {
		ProveedorID   string  `json:"proveedorId" binding:"required"`
		Descripcion   string  `json:"descripcion"`
		VigenciaDesde string  `json:"vigenciaDesde" binding:"required"`
		VigenciaHasta string  `json:"vigenciaHasta" binding:"required"`
		MontoAnual    float64 `json:"montoAnual" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	desde, err1 := time.Parse("2006-01-02", body.VigenciaDesde)
	hasta, err2 := time.Parse("2006-01-02", body.VigenciaHasta)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("vigenciaDesde y vigenciaHasta deben ser YYYY-MM-DD"))
		return
	}
	req := appprov.CreateContratoRequest{
		ProveedorID: body.ProveedorID, Descripcion: body.Descripcion,
		VigenciaDesde: desde, VigenciaHasta: hasta, MontoAnual: body.MontoAnual,
	}
	cont, err := h.createContrato.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, appprov.ErrProveedorNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Proveedor no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cont})
}

func (h *ProveedorHandlerImpl) UpdateContrato(c *gin.Context) {
	var body struct {
		Estado string `json:"estado" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	cont, err := h.updateContrato.Execute(c.Request.Context(), c.Param("id"), proveedor.EstadoContrato(body.Estado))
	if err != nil {
		if errors.Is(err, appprov.ErrContratoNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Contrato no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cont})
}

func (h *ProveedorHandlerImpl) GetRanking(c *gin.Context) {
	items, err := h.getRanking.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
