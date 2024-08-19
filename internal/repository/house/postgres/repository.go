package houserepositorypostgres

import (
	"avito/internal/model"
	houserepository "avito/internal/repository/house"
	houserepositoryconverter "avito/internal/repository/house/converter"
	houserepositorymodel "avito/internal/repository/house/model"
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
	db     *sql.DB
	logger *slog.Logger
}

func (r *repository) Create(ctx context.Context, house model.House) error {
	houseRepoModel := houserepositoryconverter.ToHouseRepModel(house)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := `INSERT INTO houses (
                    house_id,
                    address,
                    year,
                    developer) VALUES ($1, $2, $3, $4)`
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for create house", "error", err.Error())
		return houserepository.ErrInternal
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		houseRepoModel.HouseId,
		houseRepoModel.Address,
		houseRepoModel.Year,
		houseRepoModel.Developer); err != nil {
		l.Error("Failed to create house", "error", err.Error())

		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			switch {
			case pgerr.Code == pgerrcode.UniqueViolation:
				return houserepository.ErrHouseAlreadyExists
			}
		}

		return houserepository.ErrInternal
	}

	return nil
}

func (r *repository) Houses(ctx context.Context, offset int, limit int) ([]model.House, error) {
	houses := make([]model.House, 0, limit)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := "SELECT house_id, address, year, developer, created_at, last_apartment_added_at FROM houses OFFSET $1 LIMIT $2"

	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for get houses list", "error", err.Error())
		return nil, houserepository.ErrInternal
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, offset, limit)
	if err != nil {
		l.Error("Failed to get houses list", "error", err.Error())
		return nil, houserepository.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		house := houserepositorymodel.House{}
		if err = rows.Scan(&house.HouseId,
			&house.Address,
			&house.Year,
			&house.Developer,
			&house.CreatedAt,
			&house.LastApartmentAddedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, houserepository.ErrHouseNotFound
			}

			l.Error("Failed to get houses list", "error", err.Error())
			return nil, houserepository.ErrInternal
		}

		houses = append(houses, houserepositoryconverter.ToHouseDto(house))
	}

	return houses, nil
}

func (r *repository) CloseConnection() error {
	return r.db.Close()
}

func New(dataSourceName string, logger *slog.Logger) (houserepository.Repository, error) {
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
