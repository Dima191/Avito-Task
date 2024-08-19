package userrepositorypostgres

import (
	"avito/internal/model"
	userrepository "avito/internal/repository/user"
	userrepositoryconverter "avito/internal/repository/user/converter"
	userrepositorymodel "avito/internal/repository/user/model"
	"avito/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"log/slog"
)

const (
	postgresDriverName = "postgres"
)

type repository struct {
	db *sql.DB

	logger *slog.Logger
}

func (r *repository) Save(ctx context.Context, user model.User) error {
	userRepModel := userrepositoryconverter.ToUserRepModel(user)
	l := logger.EndToEndLogging(ctx, r.logger)

	q := "INSERT INTO users(user_id, role, email, hash_password) VALUES ($1, $2, $3, $4)"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for save user", "error", err.Error())
		return userrepository.ErrInternal
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		userRepModel.ID,
		userRepModel.Role,
		userRepModel.Email,
		userRepModel.HashPassword); err != nil {
		l.Error("Failed to save user", "error", err.Error())

		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			switch {
			case pgerr.Code == pgerrcode.UniqueViolation:
				return userrepository.ErrEmailAlreadyTaken
			}
		}
		return userrepository.ErrInternal
	}

	return nil
}

func (r *repository) UserByEmail(ctx context.Context, email string) (model.User, error) {
	l := logger.EndToEndLogging(ctx, r.logger)

	q := "SELECT user_id, role, email, hash_password FROM users WHERE email = $1"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for save user", "error", err.Error())
		return model.User{}, userrepository.ErrInternal
	}
	defer stmt.Close()

	user := userrepositorymodel.User{}

	if err = stmt.QueryRowContext(ctx, email).Scan(
		&user.ID,
		&user.Role,
		&user.Email,
		&user.HashPassword); err != nil {
		l.Error("Failed to save user", "error", err.Error())

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.User{}, userrepository.ErrUserNotFound
		default:
			return model.User{}, userrepository.ErrInternal
		}
	}

	return userrepositoryconverter.ToUserDTO(user), nil
}

func (r *repository) CloseConnection() error {
	return r.db.Close()
}

func New(dataSourceName string, logger *slog.Logger) (userrepository.Repository, error) {
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
