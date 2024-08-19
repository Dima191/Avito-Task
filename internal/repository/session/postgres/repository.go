package sessionrepositorypostgres

import (
	"avito/internal/model"
	sessionrepository "avito/internal/repository/session"
	sessionrepositoryconverter "avito/internal/repository/session/converter"
	sessionrepositorymodel "avito/internal/repository/session/model"
	"avito/pkg/logger"
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log/slog"
)

const (
	postgresDriverName = "postgres"
)

type repository struct {
	db *sql.DB

	logger *slog.Logger
}

func (r *repository) CheckSessionByUserId(ctx context.Context, userID uint32) (exists bool, err error) {
	l := logger.EndToEndLogging(ctx, r.logger)

	q := "SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = $1)"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to check for session existence", "error", err.Error())
		return false, sessionrepository.ErrInternal
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, userID).Scan(&exists); err != nil {
		l.Error("Failed to check for session existence", "error", err.Error())
		return false, sessionrepository.ErrInternal
	}

	return exists, nil
}

func (r *repository) SessionByUserId(ctx context.Context, userID uint32) (session model.Session, err error) {
	l := logger.EndToEndLogging(ctx, r.logger)

	q := "SELECT session_id, user_id, hash_refresh_token, expires_at FROM sessions WHERE user_id = $1"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to check for session existence", "error", err.Error())
		return model.Session{}, sessionrepository.ErrInternal
	}
	defer stmt.Close()

	sessionRepModel := sessionrepositorymodel.Session{}

	if err = stmt.QueryRowContext(ctx, userID).Scan(
		&sessionRepModel.SessionID,
		&sessionRepModel.UserID,
		&sessionRepModel.HashRefreshToken,
		&sessionRepModel.ExpiresAt); err != nil {
		l.Error("Failed to check for session existence", "error", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return model.Session{}, sessionrepository.ErrNoSession
		}

		return model.Session{}, sessionrepository.ErrInternal
	}

	return sessionrepositoryconverter.ToSessionDTO(sessionRepModel), nil
}

func (r *repository) Create(ctx context.Context, session model.Session) error {
	sessionRepModel := sessionrepositoryconverter.ToSessionRepModel(session)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := "INSERT INTO sessions (session_id, user_id, hash_refresh_token, expires_at) VALUES ($1, $2, $3, $4)"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for save session", "error", err.Error())
		return sessionrepository.ErrInternal
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sessionRepModel.SessionID,
		sessionRepModel.UserID,
		sessionRepModel.HashRefreshToken,
		sessionRepModel.ExpiresAt)
	if err != nil {
		l.Error("Failed to save session", "error", err.Error())
		return sessionrepository.ErrInternal
	}

	return nil
}

func (r *repository) ResetSession(ctx context.Context, session model.Session) error {
	sessionRepModel := sessionrepositoryconverter.ToSessionRepModel(session)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := "UPDATE sessions SET hash_refresh_token = $1, expires_at = $2 WHERE user_id = $3"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for session reset", "error", err.Error())
		return sessionrepository.ErrInternal
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		sessionRepModel.HashRefreshToken,
		sessionRepModel.ExpiresAt,
		sessionRepModel.UserID); err != nil {
		l.Error("Failed to reset session", "error", err.Error())
		return sessionrepository.ErrInternal
	}

	return nil
}

func (r *repository) CloseConnection() error {
	return r.db.Close()
}

func New(logger *slog.Logger, dataSourceName string) (sessionrepository.Repository, error) {
	r := &repository{
		logger: logger,
	}
	db, err := sql.Open(postgresDriverName, dataSourceName)
	if err != nil {
		logger.Error("failed to open postgres database connection", "error", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("failed to ping postgres database connection", "error", err.Error())
		return nil, err
	}

	r.db = db

	return r, nil
}
