package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appusr "github.com/ipnext/admin-backend/internal/application/usuario"
)

type UsuarioHandlerImpl struct {
	list   *appusr.ListUseCase
	create *appusr.CreateUseCase
	update *appusr.UpdateUseCase
	delete *appusr.DeleteUseCase
}

func NewUsuarioHandler(
	list *appusr.ListUseCase,
	create *appusr.CreateUseCase,
	update *appusr.UpdateUseCase,
	del *appusr.DeleteUseCase,
) *UsuarioHandlerImpl {
	return &UsuarioHandlerImpl{list: list, create: create, update: update, delete: del}
}

func (h *UsuarioHandlerImpl) List(c *gin.Context) {
	users, err := h.list.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UsuarioHandlerImpl) Create(c *gin.Context) {
	var req appusr.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	u, err := h.create.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": u})
}

func (h *UsuarioHandlerImpl) Update(c *gin.Context) {
	var req appusr.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	u, err := h.update.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, appusr.ErrUsuarioNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Usuario no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

func (h *UsuarioHandlerImpl) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, appusr.ErrUsuarioNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Usuario no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Usuario eliminado"}})
}
