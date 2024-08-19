package apartmentrepositorypostgres

import (
	"avito/internal/model"
	apartmentrepository "avito/internal/repository/apartment"
	apartmentrepositoryconverter "avito/internal/repository/apartment/converter"
	apartmentrepositorymodel "avito/internal/repository/apartment/model"
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

func (r *repository) Create(ctx context.Context, apartment model.Apartment) error {
	apartmentRepModel := apartmentrepositoryconverter.ToApartmentRepModel(apartment)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := "INSERT INTO apartments(apartment_id, apartment_number, house_id, price, number_of_rooms, moderation_status) VALUES ($1, $2, $3, $4, $5, $6)"
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for create apartment", "error", err.Error())
		return apartmentrepository.ErrInternal
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		apartmentRepModel.ID,
		apartmentRepModel.ApartmentNumber,
		apartmentRepModel.HouseID,
		apartmentRepModel.Price,
		apartmentRepModel.NumberOfRooms,
		apartmentRepModel.ModerationStatus); err != nil {
		l.Error("Failed to create apartment", "error", err.Error())

		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			switch {
			case pgerr.Code == pgerrcode.ForeignKeyViolation:
				return apartmentrepository.ErrInvalidHouseID
			}
		}

		return apartmentrepository.ErrInternal
	}

	return nil
}

func (r *repository) Update(ctx context.Context, apartment model.Apartment) error {
	apartmentRepModel := apartmentrepositoryconverter.ToApartmentRepModel(apartment)

	l := logger.EndToEndLogging(ctx, r.logger)

	q := `UPDATE apartments SET 
					  apartment_number = $1,
					  house_id = $2,
					  price = $3,
					  number_of_rooms = $4,
					  moderation_status = $5 WHERE apartment_id = $6`
	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for update apartment", "error", err.Error())
		return apartmentrepository.ErrInternal
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		apartmentRepModel.ApartmentNumber,
		apartmentRepModel.HouseID,
		apartmentRepModel.Price,
		apartmentRepModel.NumberOfRooms,
		apartmentRepModel.ModerationStatus,
		apartmentRepModel.ID); err != nil {
		l.Info("Failed to update apartment", "error", err.Error())

		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			switch {
			case pgerr.Code == pgerrcode.ForeignKeyViolation:
				return apartmentrepository.ErrInvalidHouseID
			}
		}

		return apartmentrepository.ErrInternal
	}

	return nil
}

func (r *repository) Apartments(ctx context.Context, houseID uint32, offset int, limit int, moderationStatusConstraint bool) ([]model.Apartment, error) {
	apartments := make([]model.Apartment, 0, limit)

	l := logger.EndToEndLogging(ctx, r.logger)

	constraint := " WHERE house_id = $1"

	if moderationStatusConstraint {
		constraint += " AND moderation_status = 'approved'"
	}

	q := "SELECT apartment_id, apartment_number, house_id, price, number_of_rooms, moderation_status FROM apartments" + constraint + " OFFSET $2 LIMIT $3"

	stmt, err := r.db.Prepare(q)
	if err != nil {
		l.Error("Failed to prepare statement for get apartments", "error", err.Error())
		return nil, apartmentrepository.ErrInternal
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, houseID, offset, limit)
	if err != nil {
		l.Error("Failed to get apartments", "error", err.Error())
		return nil, apartmentrepository.ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		apartment := apartmentrepositorymodel.Apartment{}

		if err = rows.Scan(
			&apartment.ID,
			&apartment.ApartmentNumber,
			&apartment.HouseID,
			&apartment.Price,
			&apartment.NumberOfRooms,
			&apartment.ModerationStatus); err != nil {
			l.Error("Failed to get apartments", "error", err.Error())
			return nil, apartmentrepository.ErrInternal
		}

		apartments = append(apartments, apartmentrepositoryconverter.ToApartmentDTO(apartment))
	}

	return apartments, nil
}

func (r *repository) CloseConnection() error {
	return r.db.Close()
}

func New(dataSourceName string, logger *slog.Logger) (apartmentrepository.Repository, error) {
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

	return r, err
}
