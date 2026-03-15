package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ipnext/admin-backend/internal/domain/usuario"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/middleware"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrUsuarioInactivo       = errors.New("usuario inactivo")
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Usuario *usuario.Usuario
	Token   string
}

type LoginUseCase struct {
	repo      usuario.Repository
	jwtSecret string
	jwtExpH   int
}

func NewLoginUseCase(repo usuario.Repository, jwtSecret string, jwtExpH int) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwtSecret: jwtSecret, jwtExpH: jwtExpH}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	u, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrCredencialesInvalidas
	}
	if !u.Activo {
		return nil, ErrUsuarioInactivo
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrCredencialesInvalidas
	}

	token, err := uc.signJWT(u)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Usuario: u, Token: token}, nil
}

func (uc *LoginUseCase) signJWT(u *usuario.Usuario) (string, error) {
	claims := middleware.Claims{
		Sub:      u.ID,
		Rol:      u.Rol,
		Permisos: u.Permisos,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(uc.jwtExpH) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
