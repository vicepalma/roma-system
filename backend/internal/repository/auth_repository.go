package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, name, rawPassword string) (*UserBasic, error)
	GetUserBasic(ctx context.Context, userID string) (*UserBasic, error)
	DeriveRole(ctx context.Context, userID string) (string, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

type UserBasic struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateUser: inserta con argon2id (mismo formato que viste en tu BD).
func (r *authRepository) CreateUser(ctx context.Context, email, name, rawPassword string) (*UserBasic, error) {
	hash, err := argon2id.CreateHash(rawPassword, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// Insert con RETURNING
	const q = `
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
ON CONFLICT (email) DO NOTHING
RETURNING id, email, name;
`
	var out UserBasic
	err = r.db.WithContext(ctx).Raw(q, email, hash, name).Row().Scan(&out.ID, &out.Email, &out.Name)
	if err != nil {
		// Si no retornó fila por conflicto, lo detectamos:
		if errors.Is(err, sql.ErrNoRows) {
			// trae el existente
			const q2 = `SELECT id, email, name FROM users WHERE email = $1;`
			if err2 := r.db.WithContext(ctx).Raw(q2, email).Row().Scan(&out.ID, &out.Email, &out.Name); err2 != nil {
				return nil, err2
			}
			return &out, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *authRepository) GetUserBasic(ctx context.Context, userID string) (*UserBasic, error) {
	const q = `SELECT id, email, name FROM users WHERE id = $1 LIMIT 1;`
	var out UserBasic
	if err := r.db.WithContext(ctx).Raw(q, userID).Row().Scan(&out.ID, &out.Email, &out.Name); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *authRepository) DeriveRole(ctx context.Context, userID string) (string, error) {
	// coach si:
	// - tiene links como coach aceptados, o
	// - tiene programas creados
	const q = `
SELECT
  CASE
    WHEN EXISTS (SELECT 1 FROM coach_links cl WHERE cl.coach_id = $1 AND cl.status = 'accepted') THEN 'coach'
    WHEN EXISTS (SELECT 1 FROM programs p WHERE p.owner_id = $1) THEN 'coach'
    ELSE 'disciple'
  END AS role;
`
	var role string
	if err := r.db.WithContext(ctx).Raw(q, userID).Row().Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}
