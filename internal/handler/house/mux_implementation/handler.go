package housemuximpl

import (
	househandler "avito/internal/handler/house"
	househandlerconverter "avito/internal/handler/house/converter"
	househandlermodel "avito/internal/handler/house/model"
	userhandler "avito/internal/handler/user"
	"avito/internal/middleware"
	houseservice "avito/internal/service/house"
	"avito/pkg/logger"
	tokenmanager "avito/pkg/token_manager"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

var _ househandler.Handler = &handler{}

const (
	moderator    = "moderator"
	defaultLimit = 20
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeKey  = "Content-Type"
)

type handler struct {
	router       *mux.Router
	houseService houseservice.Service

	tm tokenmanager.Manager

	logger *slog.Logger
}

func (h *handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		var house househandlermodel.House
		if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
			l.Error("Failed to decode request body", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		houseDTO := househandlerconverter.ToHouseDTO(house)
		if err := h.houseService.Create(r.Context(), houseDTO); err != nil {
			switch {
			case errors.Is(err, houseservice.ErrHouseAlreadyExists):
				http.Error(w, househandler.ErrHouseAlreadyExists.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set(ContentTypeKey, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(househandlerconverter.ToHouseHandlerModel(houseDTO))
	}
}

func (h *handler) Houses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		//PARSE URL PARAMS
		u, err := url.Parse(r.RequestURI)
		if err != nil {
			l.Error("Failed to parse request URI", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		values, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			l.Error("Failed to parse query parameters", slog.String("error", err.Error()))
			http.Error(w, userhandler.ErrInvalidURLParams.Error(), http.StatusBadRequest)
			return
		}

		limitStr := values.Get(househandler.LimitQueryParams)
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = defaultLimit
		}

		offsetStr := values.Get(househandler.OffsetQueryParams)
		offset, _ := strconv.Atoi(offsetStr)

		houses, err := h.houseService.Houses(r.Context(), offset, limit)
		if err != nil {
			switch {
			case errors.Is(err, houseservice.ErrHouseNotFound):
				http.Error(w, househandler.ErrHouseNotFound.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		housesHandlerModel := make([]househandlermodel.House, 0, len(houses))
		for _, house := range houses {
			housesHandlerModel = append(housesHandlerModel, househandlerconverter.ToHouseHandlerModel(house))
		}

		w.Header().Set(ContentTypeKey, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(housesHandlerModel)
	}
}

func Register(router *mux.Router, houseService houseservice.Service, tm tokenmanager.Manager, logger *slog.Logger) error {
	h := &handler{
		router:       router,
		houseService: houseService,
		tm:           tm,
		logger:       logger,
	}

	apiRouter := router.PathPrefix(househandler.APIUrl).Subrouter()
	apiRouter.Use(middleware.Log(logger), middleware.AuthOnly(tm))

	apiRouter.Path(househandler.HouseUrl).Handler(h.Houses()).Methods(http.MethodGet)

	moderationRouter := apiRouter.NewRoute().Subrouter()
	moderationRouter.Use(middleware.CheckRole(tm, moderator))
	moderationRouter.Path(househandler.CreateHouseUrl).Handler(h.Create()).Methods(http.MethodPost)

	return nil
}
