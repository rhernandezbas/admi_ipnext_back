package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appauth "github.com/ipnext/admin-backend/internal/application/auth"
	"github.com/ipnext/admin-backend/internal/domain/usuario"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/middleware"
)

type AuthHandlerImpl struct {
	login *appauth.LoginUseCase
	me    *appauth.GetMeUseCase
}

func NewAuthHandler(login *appauth.LoginUseCase, me *appauth.GetMeUseCase) *AuthHandlerImpl {
	return &AuthHandlerImpl{login: login, me: me}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	resp, err := h.login.Execute(c.Request.Context(), appauth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, appauth.ErrCredencialesInvalidas) || errors.Is(err, appauth.ErrUsuarioInactivo) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "UNAUTHORIZED", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Error interno"},
		})
		return
	}

	c.SetCookie("token", resp.Token, int((8 * time.Hour).Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"usuario": toUsuarioResponse(resp.Usuario)},
	})
}

func (h *AuthHandlerImpl) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Sesión cerrada"}})
}

func (h *AuthHandlerImpl) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHORIZED", "message": "No autenticado"},
		})
		return
	}

	u, err := h.me.Execute(c.Request.Context(), claims.Sub)
	if err != nil {
		if errors.Is(err, appauth.ErrUsuarioNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Usuario no encontrado"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Error interno"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"usuario": toUsuarioResponse(u)}})
}

type usuarioResponse struct {
	ID        string           `json:"id"`
	Nombre    string           `json:"nombre"`
	Email     string           `json:"email"`
	Rol       string           `json:"rol"`
	Permisos  usuario.Permisos `json:"permisos"`
	Avatar    *string          `json:"avatar"`
	CreatedAt time.Time        `json:"creadoEn"`
}

func toUsuarioResponse(u *usuario.Usuario) *usuarioResponse {
	return &usuarioResponse{
		ID:        u.ID,
		Nombre:    u.Nombre,
		Email:     u.Email,
		Rol:       string(u.Rol),
		Permisos:  u.Permisos,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}
