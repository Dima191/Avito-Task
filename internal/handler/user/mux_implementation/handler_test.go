package usermuximpl

import (
	userhandler "avito/internal/handler/user"
	userhandlermodel "avito/internal/handler/user/model"
	"avito/internal/middleware"
	sessionservice "avito/internal/service/session"
	userservice "avito/internal/service/user"
	"avito/internal/validator"
	stubwriter "avito/pkg/stub_writer"
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {
	ctrl, mockUserService, mockSessionService, _, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name        string
		statusCode  int
		prepareFunc func() userhandlermodel.User
	}{
		{
			name:       "OK",
			statusCode: http.StatusOK,
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "test@gmail.com",
					Password: "123456",
				}

				mockUserService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				mockSessionService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)
				return user
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user := c.prepareFunc()

			userBytes, err := json.Marshal(user)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.RegistrationUrl, bytes.NewReader(userBytes))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
		})
	}
}

func TestRegistrationErr(t *testing.T) {
	ctrl, mockUserService, mockSessionService, _, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() userhandlermodel.User
	}{
		{
			name:           "INVALID PASSWORD",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: validator.ErrInvalidPassword.Error(),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "test@gmail.com",
					Password: "123",
				}

				return user
			},
		},
		{
			name:           "INVALID EMAIL",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: validator.ErrInvalidEmail.Error(),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "invalid email",
					Password: "123456",
				}

				return user
			},
		},
		{
			name:           "INVALID ROLE",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: validator.ErrInvalidRole.Error(),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "invalid role",
					Email:    "test@gmail.com",
					Password: "123456",
				}

				return user
			},
		},
		{
			name:           "EMAIL ALREADY TAKEN",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: userhandler.ErrEmailAlreadyTaken.Error(),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "test@gmail.com",
					Password: "123456",
				}

				mockUserService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(userservice.ErrEmailAlreadyTaken)
				return user
			},
		},
		{
			name:           "USER SERVICE INTERNAL ERROR",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "test@gmail.com",
					Password: "123456",
				}

				mockUserService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(userservice.ErrInternal)
				return user
			},
		},
		{
			name:           "SESSION SERVICE INTERNAL ERROR",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() userhandlermodel.User {
				user := userhandlermodel.User{
					Role:     "client",
					Email:    "test@gmail.com",
					Password: "123456",
				}

				mockUserService.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				mockSessionService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", sessionservice.ErrInternal)
				return user
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user := c.prepareFunc()

			userBytes, err := json.Marshal(user)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.RegistrationUrl, bytes.NewReader(userBytes))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), c.expectedErrMsg)
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl, mockUserService, mockSessionService, _, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name        string
		statusCode  int
		prepareFunc func() *http.Request
	}{
		{
			name:       "OK",
			statusCode: http.StatusOK,
			prepareFunc: func() *http.Request {
				user := &userhandlermodel.User{
					Email:    "test@gmail.com",
					Password: "123456",
				}

				userBytes, err := json.Marshal(user)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.LoginUrl, bytes.NewReader(userBytes))

				mockUserService.EXPECT().LogIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), nil)
				mockSessionService.EXPECT().ResetSession(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)

				return req
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := c.prepareFunc()
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
		})
	}
}

func TestLoginErr(t *testing.T) {
	ctrl, mockUserService, mockSessionService, _, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() *http.Request
	}{
		{
			name:           "CREDENTIAL ERROR",
			statusCode:     http.StatusNotFound,
			expectedErrMsg: userservice.ErrCredentialsInvalid.Error(),
			prepareFunc: func() *http.Request {
				user := &userhandlermodel.User{
					Email:    "test@gmail.com",
					Password: "123456",
				}

				userBytes, err := json.Marshal(user)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.LoginUrl, bytes.NewReader(userBytes))

				mockUserService.EXPECT().LogIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), userservice.ErrCredentialsInvalid)

				return req
			},
		},
		{
			name:           "USER SERVICE INTERNAL ERROR",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				user := &userhandlermodel.User{
					Email:    "test@gmail.com",
					Password: "123456",
				}

				userBytes, err := json.Marshal(user)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.LoginUrl, bytes.NewReader(userBytes))

				mockUserService.EXPECT().LogIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), userservice.ErrInternal)
				//mockSessionService.EXPECT().ResetSession(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)

				return req
			},
		},
		{
			name:           "SESSION SERVICE INTERNAL ERROR",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				user := &userhandlermodel.User{
					Email:    "test@gmail.com",
					Password: "123456",
				}

				userBytes, err := json.Marshal(user)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, userhandler.APIUrl+userhandler.LoginUrl, bytes.NewReader(userBytes))

				mockUserService.EXPECT().LogIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), nil)
				mockSessionService.EXPECT().ResetSession(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", sessionservice.ErrInternal)

				return req
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := c.prepareFunc()
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), c.expectedErrMsg)
		})
	}
}

func TestUpdateTokens(t *testing.T) {
	ctrl, _, mockSessionService, mockTokenManager, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name        string
		statusCode  int
		prepareFunc func() *http.Request
	}{
		{
			name:       "OK",
			statusCode: http.StatusOK,
			prepareFunc: func() *http.Request {
				values := make(url.Values)
				values.Set("refresh_token", "refresh_token")

				u := fmt.Sprintf("%s%s?%s", userhandler.APIUrl, userhandler.UpdateTokensUrl, values.Encode())

				req := httptest.NewRequest(http.MethodGet, u, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = time.Now().Add(5 * time.Minute)

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockSessionService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)

				return req
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := c.prepareFunc()
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
		})
	}
}

func TestUpdateTokensErr(t *testing.T) {
	ctrl, _, mockSessionService, mockTokenManager, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() *http.Request
	}{
		{
			name:           "ERR NO SESSION",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: userhandler.ErrNoSession.Error(),
			prepareFunc: func() *http.Request {
				values := make(url.Values)
				values.Set("refresh_token", "refresh_token")

				u := fmt.Sprintf("%s%s?%s", userhandler.APIUrl, userhandler.UpdateTokensUrl, values.Encode())

				req := httptest.NewRequest(http.MethodGet, u, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = time.Now().Add(5 * time.Minute)

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockSessionService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", sessionservice.ErrNoSession)

				return req
			},
		},
		{
			name:           "INVALID REFRESH TOKEN",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: userhandler.ErrInvalidRefreshToken.Error(),
			prepareFunc: func() *http.Request {
				values := make(url.Values)
				values.Set("refresh_token", "refresh_token")

				u := fmt.Sprintf("%s%s?%s", userhandler.APIUrl, userhandler.UpdateTokensUrl, values.Encode())

				req := httptest.NewRequest(http.MethodGet, u, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = time.Now().Add(5 * time.Minute)

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockSessionService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", sessionservice.ErrInvalidRefreshToken)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				values := make(url.Values)
				values.Set("refresh_token", "refresh_token")

				u := fmt.Sprintf("%s%s?%s", userhandler.APIUrl, userhandler.UpdateTokensUrl, values.Encode())

				req := httptest.NewRequest(http.MethodGet, u, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = time.Now().Add(5 * time.Minute)

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockSessionService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", sessionservice.ErrInternal)

				return req
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := c.prepareFunc()
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, c.statusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), c.expectedErrMsg)
		})
	}
}

func testHandler(t *testing.T) (ctrl *gomock.Controller, mockUserService *userservice.MockService, mockSessionService *sessionservice.MockService, mockTokenManager *tokenmanager.MockManager, router *mux.Router) {
	ctrl = gomock.NewController(t)

	mockUserService = userservice.NewMockService(ctrl)
	mockSessionService = sessionservice.NewMockService(ctrl)
	mockTokenManager = tokenmanager.NewMockManager(ctrl)

	router = mux.NewRouter()

	logger := slog.New(slog.NewTextHandler(&stubwriter.Writer{}, nil))

	err := Register(router, mockUserService, mockSessionService, mockTokenManager, logger)
	assert.NoError(t, err)

	return ctrl, mockUserService, mockSessionService, mockTokenManager, router
}
