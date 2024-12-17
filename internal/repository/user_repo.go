package repository

import (
	"context"
	"errors"

	"example.com/authservice/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	u := &models.User{}
	err = r.conn.QueryRow(ctx, "SELECT id, email FROM users WHERE id=$1", uid).Scan(&u.ID, &u.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return u, nil
}
