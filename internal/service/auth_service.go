package service

import (
	"context"
	"errors"
	"time"

	"example.com/authservice/internal/repository"
	"example.com/authservice/internal/tokens"
	"example.com/authservice/internal/utils"
	"github.com/google/uuid"
)

const (
	accessTokenTTL  = time.Minute * 15
	refreshTokenTTL = time.Hour * 24 * 7 // неделя
)

type AuthService struct {
	tokenRepo    repository.TokenRepository
	userRepo     repository.UserRepository
	jwtSecret    string
	emailService *EmailService
}

func NewAuthService(tokenRepo repository.TokenRepository, userRepo repository.UserRepository, jwtSecret string, emailService *EmailService) *AuthService {
	return &AuthService{
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		emailService: emailService,
	}
}

func (a *AuthService) IssueTokens(ctx context.Context, userID, ip string) (accessToken string, refreshToken string, err error) {
	
	_, err = a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	accessToken, err = tokens.GenerateAccessToken(a.jwtSecret, userID, ip, accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	rawRefresh, err := tokens.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	refreshHash, err := utils.HashPassword(rawRefresh)
	if err != nil {
		return "", "", err
	}

	err = a.tokenRepo.Create(ctx, userID, refreshHash, ip, time.Now().Add(refreshTokenTTL))
	if err != nil {
		return "", "", err
	}

	return accessToken, rawRefresh, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, oldAccess, oldRefresh, currentIP string) (accessToken, refreshToken string, err error) {
	claims, err := tokens.ValidateAccessToken(a.jwtSecret, oldAccess)
	if err != nil {
		return "", "", errors.New("invalid access token")
	}

	userID := claims.UserID
	oldIP := claims.IPAddress


	rt, err := a.tokenRepo.GetLatestForUser(ctx, userID)
	if err != nil {
		return "", "", err
	}

	if rt.Used {
		return "", "", errors.New("refresh token already used")
	}


	if time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("refresh token expired")
	}


	if err := utils.CheckPasswordHash(oldRefresh, rt.RefreshHash); err != nil {
		return "", "", errors.New("refresh token does not match")
	}

	user, err := a.userRepo.GetByID(ctx, userID)
	if err == nil && oldIP != currentIP {
		a.emailService.SendWarningEmail(user.Email, oldIP, currentIP)
	}

	if err := a.tokenRepo.MarkUsed(ctx, rt.ID); err != nil {
		return "", "", err
	}

	return a.IssueTokens(ctx, userID, currentIP)
}
