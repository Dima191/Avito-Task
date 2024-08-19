package sessionserviceimpl

import (
	"avito/internal/model"
	sessionrepository "avito/internal/repository/session"
	sessionservice "avito/internal/service/session"
	"avito/pkg/hasher"
	"avito/pkg/logger"
	tokenmanager "avito/pkg/token_manager"
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type service struct {
	sessionRepository sessionrepository.Repository

	tokenManager tokenmanager.Manager

	logger *slog.Logger
}

func (s *service) generateTokens(userID uint32, role string, l *slog.Logger) (accessToken string, refreshToken string, refreshTokenExpiresAt time.Time, err error) {
	accessToken, err = s.tokenManager.GenerateAccessToken(userID, role)
	if err != nil {
		l.Error("Failed to generate access token", "error", err.Error())
		return "", "", time.Time{}, sessionservice.ErrInternal
	}

	refreshToken, refreshTokenExpiresAt, err = s.tokenManager.GenerateRefreshToken()
	if err != nil {
		l.Error("Failed to generate refresh token", "error", err.Error())
		return "", "", time.Time{}, sessionservice.ErrInternal
	}

	return accessToken, refreshToken, refreshTokenExpiresAt, nil
}

func (s *service) Create(ctx context.Context, userID uint32, role string) (accessToken, refreshToken string, err error) {
	l := logger.EndToEndLogging(ctx, s.logger)

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(userID, role, l)
	if err != nil {
		return "", "", err
	}

	hashRefreshToken, err := hasher.Hash(refreshToken)
	if err != nil {
		l.Error("Failed to hash refresh token", "error", err.Error())
		return "", "", sessionservice.ErrInternal
	}

	session := model.Session{
		SessionID:        uuid.New().ID(),
		UserID:           userID,
		HashRefreshToken: hashRefreshToken,
		ExpiresAt:        refreshTokenExpiresAt,
	}

	if err = s.sessionRepository.Create(ctx, session); err != nil {
		return "", "", sessionservice.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *service) Update(ctx context.Context, userID uint32, role string, expiredRefreshToken string) (accessToken, refreshToken string, err error) {
	l := logger.EndToEndLogging(ctx, s.logger)

	session, err := s.sessionRepository.SessionByUserId(ctx, userID)
	switch {
	case errors.Is(err, sessionrepository.ErrNoSession):
		return "", "", sessionservice.ErrNoSession
	case err != nil:
		return "", "", sessionservice.ErrInternal
	}

	if time.Now().After(session.ExpiresAt) {
		return "", "", sessionservice.ErrInvalidRefreshToken
	}

	err = hasher.Compare(expiredRefreshToken, session.HashRefreshToken)
	if err != nil {
		return "", "", sessionservice.ErrInvalidRefreshToken
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(userID, role, l)
	if err != nil {
		return "", "", err
	}

	hashRefreshToken, err := hasher.Hash(refreshToken)
	if err != nil {
		l.Error("Failed to hash refresh token", "error", err.Error())
		return "", "", sessionservice.ErrInternal
	}

	session = model.Session{
		UserID:           userID,
		HashRefreshToken: hashRefreshToken,
		ExpiresAt:        refreshTokenExpiresAt,
	}

	if err = s.sessionRepository.ResetSession(ctx, session); err != nil {
		return "", "", sessionservice.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *service) checkSessionByUserId(ctx context.Context, userID uint32) (exist bool, err error) {
	exist, err = s.sessionRepository.CheckSessionByUserId(ctx, userID)
	if err != nil {
		return false, sessionservice.ErrInternal
	}

	return exist, nil
}

func (s *service) ResetSession(ctx context.Context, userID uint32, role string) (accessToken, refreshToken string, err error) {
	l := logger.EndToEndLogging(ctx, s.logger)
	//SEARCHING USER SESSION
	exist, err := s.checkSessionByUserId(ctx, userID)
	if err != nil {
		return "", "", err
	}

	//CREATE A NEW ONE SESSION
	if !exist {
		return s.Create(ctx, userID, role)
	}

	//GENERATE TOKENS
	var refreshTokenExpiresAt time.Time
	accessToken, refreshToken, refreshTokenExpiresAt, err = s.generateTokens(userID, role, l)
	if err != nil {
		return "", "", err
	}

	hashRefreshToken, err := hasher.Hash(refreshToken)
	if err != nil {
		l.Error("Failed to hash refresh token", "error", err.Error())
		return "", "", sessionservice.ErrInternal
	}

	session := model.Session{
		UserID:           userID,
		HashRefreshToken: hashRefreshToken,
		ExpiresAt:        refreshTokenExpiresAt,
	}

	//RESET SESSION
	err = s.sessionRepository.ResetSession(ctx, session)
	if err != nil {
		return "", "", sessionservice.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func New(sessionRepository sessionrepository.Repository, tokenManager tokenmanager.Manager, logger *slog.Logger) sessionservice.Service {
	s := &service{
		sessionRepository: sessionRepository,
		tokenManager:      tokenManager,
		logger:            logger,
	}
	return s
}
