package service

import (
	"context"
	"testing"
	"time"

	"example.com/authservice/internal/models"
	"example.com/authservice/internal/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTokenRepo struct {
	mock.Mock
}

func (m *mockTokenRepo) Create(ctx context.Context, userID, refreshHash, ipAddress string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, refreshHash, ipAddress, expiresAt)
	return args.Error(0)
}

func (m *mockTokenRepo) GetLatestForUser(ctx context.Context, userID string) (*models.RefreshTokenRecord, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.RefreshTokenRecord), args.Error(1)
}

func (m *mockTokenRepo) MarkUsed(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestIssueTokens(t *testing.T) {
	tokenRepo := new(mockTokenRepo)
	userRepo := new(mockUserRepo)
	emailService := NewEmailService("no-reply@example.com")

	userRepo.On("GetByID", mock.Anything, "11111111-1111-1111-1111-111111111111").Return(&models.User{Email:"test@example.com"}, nil)
	tokenRepo.On("Create", mock.Anything, "11111111-1111-1111-1111-111111111111", mock.Anything, "127.0.0.1", mock.Anything).Return(nil)

	authSrv := NewAuthService(tokenRepo, userRepo, "supersecretkey", emailService)

	access, refresh, err := authSrv.IssueTokens(context.Background(), "11111111-1111-1111-1111-111111111111", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)
}
