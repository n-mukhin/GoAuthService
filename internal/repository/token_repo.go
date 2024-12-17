package repository

import (
	"context"
	"time"

	"example.com/authservice/internal/models"
	"github.com/jackc/pgx/v5"
)

type TokenRepository interface {
	Create(ctx context.Context, userID, refreshHash, ipAddress string, expiresAt time.Time) error
	GetLatestForUser(ctx context.Context, userID string) (*models.RefreshTokenRecord, error)
	MarkUsed(ctx context.Context, id int) error
}

type tokenRepository struct {
	conn *pgx.Conn
}

func NewTokenRepository(conn *pgx.Conn) TokenRepository {
	return &tokenRepository{conn: conn}
}

func (r *tokenRepository) Create(ctx context.Context, userID, refreshHash, ipAddress string, expiresAt time.Time) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO refresh_tokens (user_id, refresh_hash, ip_address, expires_at) VALUES ($1, $2, $3, $4)",
		userID, refreshHash, ipAddress, expiresAt)
	return err
}

func (r *tokenRepository) GetLatestForUser(ctx context.Context, userID string) (*models.RefreshTokenRecord, error) {
	rt := &models.RefreshTokenRecord{}
	err := r.conn.QueryRow(ctx, "SELECT id, user_id, refresh_hash, ip_address, created_at, expires_at, used FROM refresh_tokens WHERE user_id=$1 AND used=false ORDER BY created_at DESC LIMIT 1",
		userID).Scan(&rt.ID, &rt.UserID, &rt.RefreshHash, &rt.IPAddress, &rt.CreatedAt, &rt.ExpiresAt, &rt.Used)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *tokenRepository) MarkUsed(ctx context.Context, id int) error {
	_, err := r.conn.Exec(ctx, "UPDATE refresh_tokens SET used=true WHERE id=$1", id)
	return err
}
