package apartmentmuximpl

import (
	apartmenthandler "avito/internal/handler/apartment"
	apartmenthandlerconverter "avito/internal/handler/apartment/converter"
	apartmenthandlermodel "avito/internal/handler/apartment/model"
	userhandler "avito/internal/handler/user"
	"avito/internal/middleware"
	apartmentservice "avito/internal/service/apartment"
	"avito/internal/validator"
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

const (
	defaultModerationStatus = "created"
	moderator               = "moderator"

	defaultLimit = 20
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeKey  = "Content-Type"
)

var _ apartmenthandler.Handler = &handler{}

type handler struct {
	router *mux.Router

	apartmentService apartmentservice.Service

	tm tokenmanager.Manager

	validator *validator.Validate

	logger *slog.Logger
}

func (h *handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		apartment := apartmenthandlermodel.Apartment{}
		if err := json.NewDecoder(r.Body).Decode(&apartment); err != nil {
			l.Error("Failed to decode request body", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		apartment.ModerationStatus = defaultModerationStatus

		//VALIDATION
		if err := h.validator.Validate(apartment); err != nil {
			l.Error("Invalid data", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		apartmentDTO := apartmenthandlerconverter.ToApartmentDTO(apartment)
		if err := h.apartmentService.Create(r.Context(), apartmentDTO); err != nil {
			switch {
			case errors.Is(err, apartmentservice.ErrInvalidHouseID):
				http.Error(w, apartmenthandler.ErrInvalidHouseID.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set(ContentTypeKey, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apartmenthandlerconverter.ToHandlerModelApartment(apartmentDTO))
	}
}

func (h *handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		apartmentIDStr := mux.Vars(r)[apartmenthandler.ApartmentID]
		apartmentID, err := strconv.Atoi(apartmentIDStr)
		if err != nil {
			l.Error("Invalid apartmentID", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		apartment := apartmenthandlermodel.Apartment{}
		if err = json.NewDecoder(r.Body).Decode(&apartment); err != nil {
			l.Error("Failed to decode request body", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		apartment.ID = uint32(apartmentID)

		if err = h.validator.Validate(apartment); err != nil {
			l.Error("Invalid data", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		apartmentDTO := apartmenthandlerconverter.ToApartmentDTO(apartment)
		if err = h.apartmentService.Update(r.Context(), apartmentDTO); err != nil {
			switch {
			case errors.Is(err, apartmentservice.ErrInvalidHouseID):
				http.Error(w, apartmenthandler.ErrInvalidHouseID.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set(ContentTypeKey, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apartmenthandlerconverter.ToHandlerModelApartment(apartmentDTO))
	}
}

func (h *handler) Apartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		houseIDStr := mux.Vars(r)[apartmenthandler.HouseID]
		houseID, err := strconv.Atoi(houseIDStr)
		if err != nil {
			l.Error("Invalid houseID", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

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

		limitStr := values.Get(apartmenthandler.LimitQueryParams)
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = defaultLimit
		}

		offsetStr := values.Get(apartmenthandler.OffsetQueryParams)
		offset, _ := strconv.Atoi(offsetStr)

		role, ok := r.Context().Value(middleware.RoleCtxKey).(string)
		if !ok {
			l.Error("Failed to get role from context")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		apartments, err := h.apartmentService.Apartments(r.Context(), uint32(houseID), offset, limit, role)
		if err != nil {
			l.Error("Failed to get apartments", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		apartmentsHandlerModel := make([]apartmenthandlermodel.Apartment, 0, len(apartments))
		for _, apartment := range apartments {
			apartmentsHandlerModel = append(apartmentsHandlerModel, apartmenthandlerconverter.ToHandlerModelApartment(apartment))
		}

		w.Header().Set(ContentTypeKey, ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apartmentsHandlerModel)
	}
}

func Register(router *mux.Router, apartmentService apartmentservice.Service, tm tokenmanager.Manager, logger *slog.Logger) error {
	h := &handler{
		router:           router,
		apartmentService: apartmentService,
		tm:               tm,
		validator:        validator.New(),
		logger:           logger,
	}

	if err := h.validator.RegisterTag(validator.ModerationStatusTag, apartmenthandlermodel.ModerationStatusValidation); err != nil {
		logger.Error("Failed to register moderation status validation", "error", err.Error())
		return err
	}

	apiRouter := router.PathPrefix(apartmenthandler.APIUrl).Subrouter()
	apiRouter.Use(middleware.Log(logger), middleware.AuthOnly(tm))

	apiRouter.Path(apartmenthandler.CreateApartmentUrl).Handler(h.Create()).Methods(http.MethodPost)
	apiRouter.Path(apartmenthandler.ApartmentsByHouseIDUrl).Handler(h.Apartments()).Methods(http.MethodGet)

	moderationRouter := apiRouter.NewRoute().Subrouter()
	moderationRouter.Use(middleware.CheckRole(tm, moderator))
	moderationRouter.Path(apartmenthandler.UpdateApartmentUrl).Handler(h.Update()).Methods(http.MethodPut)

	return nil
}
