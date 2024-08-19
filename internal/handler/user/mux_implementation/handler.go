package usermuximpl

import (
	userhandler "avito/internal/handler/user"
	userhandlerconverter "avito/internal/handler/user/converter"
	userhandlermodel "avito/internal/handler/user/model"
	"avito/internal/middleware"
	sessionservice "avito/internal/service/session"
	userservice "avito/internal/service/user"
	"avito/internal/validator"
	"avito/pkg/logger"
	tokenmanager "avito/pkg/token_manager"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"net/url"
)

type handler struct {
	router *mux.Router

	userService    userservice.Service
	sessionService sessionservice.Service

	validator *validator.Validate

	logger *slog.Logger
}

func (h *handler) Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		user := userhandlermodel.User{}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			l.Error("Failed to decode request body", "error", err.Error())
			http.Error(w, userhandler.ErrDecodeBody.Error(), http.StatusBadRequest)
			return
		}

		//VALIDATION
		if err := h.validator.Validate(user); err != nil {
			l.Error("Invalid data", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//CONVERT TO DTO OBJECT

		userDto, err := userhandlerconverter.ToUserDto(user)
		if err != nil {
			l.Error("Failed to convert user to dto", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		//SAVE USER
		if err = h.userService.Save(r.Context(), userDto); err != nil {
			switch {
			case errors.Is(err, userservice.ErrEmailAlreadyTaken):
				http.Error(w, userhandler.ErrEmailAlreadyTaken.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		//CREATE SESSION
		accessToken, refreshToken, err := h.sessionService.Create(r.Context(), userDto.ID, userDto.Role)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tokens := userhandlermodel.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		tokensBytes, err := json.Marshal(tokens)
		if err != nil {
			l.Error("Failed to marshal tokens to json", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(tokensBytes)
	}
}

func (h *handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.EndToEndLogging(r.Context(), h.logger)

		user := &userhandlermodel.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			l.Error("Failed to decode request body", "error", err.Error())
			http.Error(w, userhandler.ErrDecodeBody.Error(), http.StatusBadRequest)
			return
		}

		userID, err := h.userService.LogIn(r.Context(), user.Email, user.Password)
		if err != nil {
			switch {
			case errors.Is(err, userservice.ErrCredentialsInvalid):
				http.Error(w, userhandler.ErrCredentialsInvalid.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		accessToken, refreshToken, err := h.sessionService.ResetSession(r.Context(), userID, user.Role)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tokens := userhandlermodel.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		tokensBytes, err := json.Marshal(tokens)
		if err != nil {
			l.Error("Failed to marshal tokens to json", "error", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(tokensBytes)
	}
}

func (h *handler) UpdateTokens() http.HandlerFunc {
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

		refreshToken := values.Get(userhandler.RefreshTokenQueryParam)
		userID := r.Context().Value(middleware.UserIDCtxKey).(uint32)
		role := r.Context().Value(middleware.RoleCtxKey).(string)

		accessToken, refreshToken, err := h.sessionService.Update(r.Context(), userID, role, refreshToken)
		if err != nil {
			switch {
			case errors.Is(err, sessionservice.ErrNoSession):
				http.Error(w, userhandler.ErrNoSession.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, sessionservice.ErrInvalidRefreshToken):
				http.Error(w, userhandler.ErrInvalidRefreshToken.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		tokens := userhandlermodel.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	}
}

func Register(router *mux.Router, userService userservice.Service, sessionService sessionservice.Service, tm tokenmanager.Manager, logger *slog.Logger) error {
	h := &handler{
		router:         router,
		userService:    userService,
		sessionService: sessionService,
		validator:      validator.New(),
		logger:         logger,
	}

	//ADD PASSWORD VALIDATION
	if err := h.validator.RegisterTag(validator.PasswordTag, userhandlermodel.PasswordValidation); err != nil {
		logger.Error("Failed to register password validation", "error", err.Error())
		return err
	}

	//ADD ROLE VALIDATION
	if err := h.validator.RegisterTag(validator.RoleTag, userhandlermodel.RoleValidation); err != nil {
		logger.Error("Failed to register role validation", "error", err.Error())
		return err
	}

	apiRouter := router.PathPrefix(userhandler.APIUrl).Subrouter()

	apiRouter.Use(middleware.Log(logger))

	apiRouter.Path(userhandler.RegistrationUrl).Handler(h.Registration()).Methods(http.MethodPost)
	apiRouter.Path(userhandler.LoginUrl).Handler(h.Login()).Methods(http.MethodPost)

	moderationRouter := apiRouter.NewRoute().Subrouter()
	moderationRouter.Use(middleware.ParseAuthToken(tm))
	moderationRouter.Path(userhandler.UpdateTokensUrl).Handler(h.UpdateTokens()).Methods(http.MethodGet)

	return nil
}
