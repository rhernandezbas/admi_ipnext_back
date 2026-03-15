package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ipnext/admin-backend/internal/domain/usuario"
	"gorm.io/gorm"
)

type usuarioModel struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	Nombre    string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Rol       string    `gorm:"type:enum('admin','sub-usuario');not null"`
	Permisos  string    `gorm:"type:json;not null"`
	Avatar    *string
	Activo    bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (usuarioModel) TableName() string { return "usuarios" }

type MySQLUsuarioRepository struct {
	db *gorm.DB
}

func NewMySQLUsuarioRepository(db *gorm.DB) usuario.Repository {
	return &MySQLUsuarioRepository{db: db}
}

func (r *MySQLUsuarioRepository) FindByID(ctx context.Context, id string) (*usuario.Usuario, error) {
	var m usuarioModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toUsuarioDomain(&m)
}

func (r *MySQLUsuarioRepository) FindByEmail(ctx context.Context, email string) (*usuario.Usuario, error) {
	var m usuarioModel
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toUsuarioDomain(&m)
}

func (r *MySQLUsuarioRepository) FindAll(ctx context.Context) ([]*usuario.Usuario, error) {
	var models []usuarioModel
	if err := r.db.WithContext(ctx).Where("activo = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*usuario.Usuario, 0, len(models))
	for _, m := range models {
		u, err := toUsuarioDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *MySQLUsuarioRepository) Save(ctx context.Context, u *usuario.Usuario) error {
	m, err := toUsuarioModel(u)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MySQLUsuarioRepository) Update(ctx context.Context, u *usuario.Usuario) error {
	m, err := toUsuarioModel(u)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MySQLUsuarioRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&usuarioModel{}).Where("id = ?", id).Update("activo", false).Error
}

func toUsuarioDomain(m *usuarioModel) (*usuario.Usuario, error) {
	var permisos usuario.Permisos
	if err := json.Unmarshal([]byte(m.Permisos), &permisos); err != nil {
		return nil, err
	}
	return &usuario.Usuario{
		ID:        m.ID,
		Nombre:    m.Nombre,
		Email:     m.Email,
		Password:  m.Password,
		Rol:       usuario.Rol(m.Rol),
		Permisos:  permisos,
		Avatar:    m.Avatar,
		Activo:    m.Activo,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func toUsuarioModel(u *usuario.Usuario) (*usuarioModel, error) {
	permisosJSON, err := json.Marshal(u.Permisos)
	if err != nil {
		return nil, err
	}
	return &usuarioModel{
		ID:        u.ID,
		Nombre:    u.Nombre,
		Email:     u.Email,
		Password:  u.Password,
		Rol:       string(u.Rol),
		Permisos:  string(permisosJSON),
		Avatar:    u.Avatar,
		Activo:    u.Activo,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
