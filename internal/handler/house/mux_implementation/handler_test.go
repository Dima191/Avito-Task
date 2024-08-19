package housemuximpl

import (
	househandler "avito/internal/handler/house"
	househandlermodel "avito/internal/handler/house/model"
	"avito/internal/middleware"
	houseservice "avito/internal/service/house"
	stubwriter "avito/pkg/stub_writer"
	tokenmanager "avito/pkg/token_manager"
	tokenmanagerimpl "avito/pkg/token_manager/implementation"
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctrl, mockHouseService, mockTokenManager, router := testHandler(t)
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
				houseBytes, err := json.Marshal(househandlermodel.House{
					Address:   "Address",
					Year:      2024,
					Developer: "Developer",
				})
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, househandler.APIUrl+househandler.CreateHouseUrl, bytes.NewReader(houseBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

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

func TestCreateErr(t *testing.T) {
	ctrl, mockHouseService, mockTokenManager, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() *http.Request
	}{
		{
			name:       "ERR STATUS FORBIDDEN",
			statusCode: http.StatusForbidden,
			prepareFunc: func() *http.Request {
				houseBytes, err := json.Marshal(househandlermodel.House{
					Address:   "Address",
					Year:      2024,
					Developer: "Developer",
				})
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, househandler.APIUrl+househandler.CreateHouseUrl, bytes.NewReader(houseBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)

				return req
			},
		},
		{
			name:           "ERR HOUSE ALREADY EXISTS",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: househandler.ErrHouseAlreadyExists.Error(),
			prepareFunc: func() *http.Request {
				houseBytes, err := json.Marshal(househandlermodel.House{
					Address:   "Address",
					Year:      2024,
					Developer: "Developer",
				})
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, househandler.APIUrl+househandler.CreateHouseUrl, bytes.NewReader(houseBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(houseservice.ErrHouseAlreadyExists)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				houseBytes, err := json.Marshal(househandlermodel.House{
					Address:   "Address",
					Year:      2024,
					Developer: "Developer",
				})
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, househandler.APIUrl+househandler.CreateHouseUrl, bytes.NewReader(houseBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(houseservice.ErrInternal)

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

func TestHouses(t *testing.T) {
	ctrl, mockHouseService, mockTokenManager, router := testHandler(t)
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
				req := httptest.NewRequest(http.MethodGet, househandler.APIUrl+househandler.HouseUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Houses(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

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

func TestHousesErr(t *testing.T) {
	ctrl, mockHouseService, mockTokenManager, router := testHandler(t)
	_ = mockTokenManager
	_ = mockHouseService
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() *http.Request
	}{
		{
			name:       "UNAUTHORIZED",
			statusCode: http.StatusUnauthorized,
			prepareFunc: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, househandler.APIUrl+househandler.HouseUrl, http.NoBody)

				return req
			},
		},
		{
			name:           "ERR NO HOUSES",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: househandler.ErrHouseNotFound.Error(),
			prepareFunc: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, househandler.APIUrl+househandler.HouseUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Houses(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, houseservice.ErrHouseNotFound)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, househandler.APIUrl+househandler.HouseUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockHouseService.EXPECT().Houses(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, houseservice.ErrInternal)

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

//func (h *handler) Houses() http.HandlerFunc {

func testHandler(t *testing.T) (ctrl *gomock.Controller, mockHouseService *houseservice.MockService, mockTokenManager *tokenmanager.MockManager, router *mux.Router) {
	ctrl = gomock.NewController(t)

	mockHouseService = houseservice.NewMockService(ctrl)
	mockTokenManager = tokenmanager.NewMockManager(ctrl)

	router = mux.NewRouter()
	logger := slog.New(slog.NewTextHandler(&stubwriter.Writer{}, nil))

	err := Register(router, mockHouseService, mockTokenManager, logger)
	assert.NoError(t, err)

	return ctrl, mockHouseService, mockTokenManager, router
}
