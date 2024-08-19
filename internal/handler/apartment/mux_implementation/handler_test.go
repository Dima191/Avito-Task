package apartmentmuximpl

import (
	apartmenthandler "avito/internal/handler/apartment"
	apartmenthandlermodel "avito/internal/handler/apartment/model"
	"avito/internal/middleware"
	apartmentservice "avito/internal/service/apartment"
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
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
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
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber: 1,
					HouseID:         1,
					Price:           1,
					NumberOfRooms:   1,
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, apartmenthandler.APIUrl+apartmenthandler.CreateApartmentUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

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
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
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
				apartment := apartmenthandlermodel.Apartment{}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, apartmenthandler.APIUrl+apartmenthandler.CreateApartmentUrl, bytes.NewReader(apartmentBytes))

				return req
			},
		},
		{
			name:       "ERR INVALID DATA",
			statusCode: http.StatusBadRequest,
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, apartmenthandler.APIUrl+apartmenthandler.CreateApartmentUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)

				return req
			},
		},
		{
			name:           "ERR INVALID HOUSE ID",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: apartmenthandler.ErrInvalidHouseID.Error(),
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber: 1,
					HouseID:         1,
					Price:           1,
					NumberOfRooms:   1,
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, apartmenthandler.APIUrl+apartmenthandler.CreateApartmentUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(apartmentservice.ErrInvalidHouseID)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber: 1,
					HouseID:         1,
					Price:           1,
					NumberOfRooms:   1,
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, apartmenthandler.APIUrl+apartmenthandler.CreateApartmentUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(apartmentservice.ErrInternal)

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

func TestUpdate(t *testing.T) {
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
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
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber:  1,
					HouseID:          1,
					Price:            1,
					NumberOfRooms:    1,
					ModerationStatus: "approved",
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

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

func TestUpdateErr(t *testing.T) {
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
	_ = mockApartmentService
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
				apartment := apartmenthandlermodel.Apartment{}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))

				return req
			},
		},
		{
			name:       "INVALID DATA",
			statusCode: http.StatusBadRequest,
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)

				return req
			},
		},
		{
			name:       "ERR FORBIDDEN",
			statusCode: http.StatusForbidden,
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))
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
			name:           "ERR INVALID HOUSE ID",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: apartmenthandler.ErrInvalidHouseID.Error(),
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber:  1,
					HouseID:          1,
					Price:            1,
					NumberOfRooms:    1,
					ModerationStatus: "approved",
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(apartmentservice.ErrInvalidHouseID)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				apartment := apartmenthandlermodel.Apartment{
					ApartmentNumber:  1,
					HouseID:          1,
					Price:            1,
					NumberOfRooms:    1,
					ModerationStatus: "approved",
				}

				apartmentBytes, err := json.Marshal(apartment)
				assert.NoError(t, err)

				updateUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.UpdateApartmentUrl, fmt.Sprintf("{%s}", apartmenthandler.ApartmentID), "1")

				req := httptest.NewRequest(http.MethodPut, updateUrl, bytes.NewReader(apartmentBytes))
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "moderator"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(apartmentservice.ErrInternal)

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

func TestApartments(t *testing.T) {
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
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
				apartmentsUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.ApartmentsByHouseIDUrl, fmt.Sprintf("{%s}", apartmenthandler.HouseID), "1")
				req := httptest.NewRequest(http.MethodGet, apartmentsUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Apartments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

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

func TestApartmentsErr(t *testing.T) {
	ctrl, mockApartmentService, mockTokenManager, router := testHandler(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		prepareFunc    func() *http.Request
	}{
		{
			name:           "ERR INVALID HOUSE ID",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: http.StatusText(http.StatusBadRequest),
			prepareFunc: func() *http.Request {
				apartmentsUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.ApartmentsByHouseIDUrl, fmt.Sprintf("{%s}", apartmenthandler.HouseID), "invalid_house_id")
				req := httptest.NewRequest(http.MethodGet, apartmentsUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)

				return req
			},
		},
		{
			name:       "UNAUTHORIZED",
			statusCode: http.StatusUnauthorized,
			prepareFunc: func() *http.Request {
				apartmentsUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.ApartmentsByHouseIDUrl, fmt.Sprintf("{%s}", apartmenthandler.HouseID), "1")
				req := httptest.NewRequest(http.MethodGet, apartmentsUrl, http.NoBody)

				return req
			},
		},
		{
			name:           "ERR INTERNAL",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: http.StatusText(http.StatusInternalServerError),
			prepareFunc: func() *http.Request {
				apartmentsUrl := strings.ReplaceAll(apartmenthandler.APIUrl+apartmenthandler.ApartmentsByHouseIDUrl, fmt.Sprintf("{%s}", apartmenthandler.HouseID), "1")
				req := httptest.NewRequest(http.MethodGet, apartmentsUrl, http.NoBody)
				req.Header.Set(middleware.AuthorizationHeader, "Bearer access-token")

				m := make(jwt.MapClaims)
				m[tokenmanagerimpl.UserIDClaimsTag] = float64(uuid.New().ID())
				m[tokenmanagerimpl.RoleClaimsTag] = "client"
				m[tokenmanagerimpl.ExpClaimsTag] = float64(time.Now().Add(5 * time.Minute).Unix())

				mockTokenManager.EXPECT().Parse(gomock.Any()).Return(m, nil)
				mockApartmentService.EXPECT().Apartments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, apartmentservice.ErrInternal)

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

func testHandler(t *testing.T) (ctrl *gomock.Controller, mockApartmentService *apartmentservice.MockService, mockTokenManager *tokenmanager.MockManager, router *mux.Router) {
	ctrl = gomock.NewController(t)

	mockApartmentService = apartmentservice.NewMockService(ctrl)
	mockTokenManager = tokenmanager.NewMockManager(ctrl)

	router = mux.NewRouter()
	logger := slog.New(slog.NewTextHandler(&stubwriter.Writer{}, nil))

	err := Register(router, mockApartmentService, mockTokenManager, logger)
	assert.NoError(t, err)

	return ctrl, mockApartmentService, mockTokenManager, router
}
