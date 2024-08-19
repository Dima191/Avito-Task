package app

import (
	"avito/internal/config"
	apartmentmuximpl "avito/internal/handler/apartment/mux_implementation"
	housemuximpl "avito/internal/handler/house/mux_implementation"
	usermuximpl "avito/internal/handler/user/mux_implementation"
	"avito/pkg/logger"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
)

type App struct {
	configPath string
	cfg        *config.Config

	server *http.Server

	router *mux.Router

	sp *serviceProvider

	logger  *slog.Logger
	isDebug bool
}

func (a *App) initLogger(_ context.Context) error {
	a.logger = logger.New(a.isDebug)
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.New(a.configPath, a.logger)
	if err != nil {
		return err
	}

	a.cfg = cfg
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.sp = newServiceProvider(a.cfg.DBUrl, a.cfg.JWTSignedKey, a.cfg.AccessTokenExpiresIn, a.cfg.RefreshTokenExpiresIn, a.logger)
	return nil
}

func (a *App) initMuxHandler(_ context.Context) error {
	a.router = mux.NewRouter()
	return nil
}

func (a *App) initUserHandler(_ context.Context) error {
	userService, err := a.sp.UserService()
	if err != nil {
		return err
	}

	sessionService, err := a.sp.SessionService()
	if err != nil {
		return err
	}

	if err = usermuximpl.Register(a.router, userService, sessionService, a.sp.TokenManager(), a.logger); err != nil {
		return err
	}

	return nil
}

func (a *App) initApartmentHandler(_ context.Context) error {
	apartmentService, err := a.sp.ApartmentService()
	if err != nil {
		return err
	}

	if err = apartmentmuximpl.Register(a.router, apartmentService, a.sp.TokenManager(), a.logger); err != nil {
		return err
	}
	return nil
}

func (a *App) initHouseHandler(_ context.Context) error {
	houseService, err := a.sp.HouseService()
	if err != nil {
		return err
	}

	if err = housemuximpl.Register(a.router, houseService, a.sp.TokenManager(), a.logger); err != nil {
		return err
	}
	return nil
}

func (a *App) initDependencies(ctx context.Context) error {
	deps := []func(ctx context.Context) error{
		a.initLogger,
		a.initConfig,
		a.initServiceProvider,
		a.initMuxHandler,
		a.initUserHandler,
		a.initApartmentHandler,
		a.initHouseHandler,
	}

	for _, f := range deps {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) runHTTPServer() error {
	a.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port),
		Handler: a.router,
	}

	if err := a.server.ListenAndServe(); err != nil {
		a.logger.Error("Failed to listen and serve HTTP server", "error", err.Error())
		return err
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return a.runHTTPServer()
	})

	if err := g.Wait(); err != nil {
		a.logger.Error("")
		a.logger.Error("Failed to start server", "error", err.Error())
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown server", "error", err.Error())
		return err
	}

	if a.sp.houseRepository != nil {
		if err := a.sp.houseRepository.CloseConnection(); err != nil {
			a.logger.Error("Failed to close house repository", "error", err.Error())
			return err
		}
	}

	if a.sp.sessionRepository != nil {
		if err := a.sp.sessionRepository.CloseConnection(); err != nil {
			a.logger.Error("Failed to close session repository", "error", err.Error())
			return err
		}
	}

	if a.sp.userRepository != nil {
		if err := a.sp.userRepository.CloseConnection(); err != nil {
			a.logger.Error("Failed to close user repository", "error", err.Error())
			return err
		}
	}

	return nil
}

func New(ctx context.Context, configPath string, isDebug bool) (*App, error) {
	a := &App{
		configPath: configPath,
		isDebug:    isDebug,
	}

	if err := a.initDependencies(ctx); err != nil {
		return nil, err
	}

	return a, nil
}
