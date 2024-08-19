package app

import (
	apartmentrepository "avito/internal/repository/apartment"
	apartmentrepositorypostgres "avito/internal/repository/apartment/postgres"
	houserepository "avito/internal/repository/house"
	houserepositorypostgres "avito/internal/repository/house/postgres"
	sessionrepository "avito/internal/repository/session"
	sessionrepositorypostgres "avito/internal/repository/session/postgres"
	userrepository "avito/internal/repository/user"
	userrepositorypostgres "avito/internal/repository/user/postgres"
	apartmentservice "avito/internal/service/apartment"
	apartmentserviceimpl "avito/internal/service/apartment/implementation"
	houseservice "avito/internal/service/house"
	houseserviceimpl "avito/internal/service/house/implementation"
	sessionservice "avito/internal/service/session"
	sessionserviceimpl "avito/internal/service/session/implementation"
	userservice "avito/internal/service/user"
	userserviceimpl "avito/internal/service/user/implementation"
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"log/slog"
	"time"
)

type serviceProvider struct {
	tokenManager tokenmanager.Manager

	dbURL string

	jwtSignedKey          string
	accessTokenExpiresIn  time.Duration
	refreshTokenExpiresIn time.Duration

	sessionRepository sessionrepository.Repository
	sessionService    sessionservice.Service

	userRepository userrepository.Repository
	userService    userservice.Service

	apartmentRepository apartmentrepository.Repository
	apartmentService    apartmentservice.Service

	houseRepository houserepository.Repository
	houseService    houseservice.Service

	logger *slog.Logger
}

func (sp *serviceProvider) SessionRepository() (sessionrepository.Repository, error) {
	if sp.sessionRepository == nil {
		rep, err := sessionrepositorypostgres.New(sp.logger, sp.dbURL)
		if err != nil {
			return nil, err
		}
		sp.sessionRepository = rep
	}
	return sp.sessionRepository, nil
}

func (sp *serviceProvider) SessionService() (sessionservice.Service, error) {
	if sp.sessionService == nil {
		rep, err := sp.SessionRepository()
		if err != nil {
			return nil, err
		}

		sp.sessionService = sessionserviceimpl.New(rep, sp.TokenManager(), sp.logger)
	}

	return sp.sessionService, nil
}

func (sp *serviceProvider) UserRepository() (userrepository.Repository, error) {
	if sp.userRepository == nil {
		rep, err := userrepositorypostgres.New(sp.dbURL, sp.logger)
		if err != nil {
			return nil, err
		}

		sp.userRepository = rep
	}

	return sp.userRepository, nil
}

func (sp *serviceProvider) UserService() (userservice.Service, error) {
	if sp.userService == nil {
		rep, err := sp.UserRepository()
		if err != nil {
			return nil, err
		}

		sp.userService = userserviceimpl.New(rep, sp.TokenManager(), sp.logger)
	}

	return sp.userService, nil
}

func (sp *serviceProvider) ApartmentRepository() (apartmentrepository.Repository, error) {
	if sp.apartmentRepository == nil {
		apartmentRepository, err := apartmentrepositorypostgres.New(sp.dbURL, sp.logger)
		if err != nil {
			return nil, err
		}

		sp.apartmentRepository = apartmentRepository
	}
	return sp.apartmentRepository, nil
}

func (sp *serviceProvider) ApartmentService() (apartmentservice.Service, error) {
	if sp.apartmentService == nil {
		apartmentRepository, err := sp.ApartmentRepository()
		if err != nil {
			return nil, err
		}

		sp.apartmentService = apartmentserviceimpl.New(apartmentRepository, sp.logger)
	}
	return sp.apartmentService, nil
}

func (sp *serviceProvider) HouseRepository() (houserepository.Repository, error) {
	if sp.houseRepository == nil {
		houseRepository, err := houserepositorypostgres.New(sp.dbURL, sp.logger)
		if err != nil {
			return nil, err
		}

		sp.houseRepository = houseRepository
	}

	return sp.houseRepository, nil
}
func (sp *serviceProvider) HouseService() (houseservice.Service, error) {
	if sp.houseService == nil {
		houseRepository, err := sp.HouseRepository()
		if err != nil {
			return nil, err
		}
		sp.houseService = houseserviceimpl.New(houseRepository, sp.logger)
	}

	return sp.houseService, nil
}

func (sp *serviceProvider) TokenManager() tokenmanager.Manager {
	if sp.tokenManager == nil {
		sp.tokenManager = tokenmanagerimpl.New(sp.jwtSignedKey, sp.accessTokenExpiresIn, sp.refreshTokenExpiresIn)
	}
	return sp.tokenManager
}

func newServiceProvider(dbURL string, jwtSignedKey string, accessTokenExpiresIn, refreshTokenExpiresIn time.Duration, logger *slog.Logger) *serviceProvider {
	sp := &serviceProvider{
		dbURL:                 dbURL,
		jwtSignedKey:          jwtSignedKey,
		accessTokenExpiresIn:  accessTokenExpiresIn,
		refreshTokenExpiresIn: refreshTokenExpiresIn,
		logger:                logger,
	}

	return sp
}
