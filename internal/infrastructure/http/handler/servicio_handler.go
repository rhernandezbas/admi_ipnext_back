package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appsvc "github.com/ipnext/admin-backend/internal/application/servicio"
	domservicio "github.com/ipnext/admin-backend/internal/domain/servicio"
)

type ServicioHandlerImpl struct {
	list    *appsvc.ListUseCase
	getKPIs *appsvc.GetKPIsUseCase
	get     *appsvc.GetUseCase
	create  *appsvc.CreateUseCase
	update  *appsvc.UpdateUseCase
	delete  *appsvc.DeleteUseCase
}

func NewServicioHandler(
	list *appsvc.ListUseCase, getKPIs *appsvc.GetKPIsUseCase,
	get *appsvc.GetUseCase, create *appsvc.CreateUseCase,
	update *appsvc.UpdateUseCase, del *appsvc.DeleteUseCase,
) *ServicioHandlerImpl {
	return &ServicioHandlerImpl{list: list, getKPIs: getKPIs, get: get, create: create, update: update, delete: del}
}

func (h *ServicioHandlerImpl) List(c *gin.Context) {
	items, err := h.list.Execute(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ServicioHandlerImpl) GetKPIs(c *gin.Context) {
	kpis, err := h.getKPIs.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

func (h *ServicioHandlerImpl) ListByTipo(c *gin.Context) {
	tipo := domservicio.Tipo(c.Param("tipo"))
	items, err := h.list.Execute(c.Request.Context(), &tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *ServicioHandlerImpl) Get(c *gin.Context) {
	s, err := h.get.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, appsvc.ErrServicioNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Servicio no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func (h *ServicioHandlerImpl) Create(c *gin.Context) {
	var req appsvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	s, err := h.create.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": s})
}

func (h *ServicioHandlerImpl) Update(c *gin.Context) {
	var req appsvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	s, err := h.update.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, appsvc.ErrServicioNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Servicio no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func (h *ServicioHandlerImpl) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, appsvc.ErrServicioNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Servicio no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Servicio eliminado"}})
}
