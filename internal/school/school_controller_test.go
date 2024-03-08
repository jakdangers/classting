package school

import (
	"bytes"
	"classting/config"
	"classting/domain"
	"classting/internal/user"
	"classting/mocks"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type schoolControllerTestSuite struct {
	router           *gin.Engine
	cfg              *config.Config
	schoolService    *mocks.SchoolService
	schoolController domain.SchoolController
}

func setupSchoolControllerTestSuite(t *testing.T) schoolControllerTestSuite {
	var us schoolControllerTestSuite

	gin.SetMode(gin.TestMode)
	us.router = gin.Default()
	us.schoolService = mocks.NewSchoolService(t)
	us.cfg = &config.Config{
		Auth: config.Auth{
			Secret: "classting_test_secret",
		},
	}

	us.schoolController = NewSchoolController(us.schoolService, us.cfg)
	RegisterRoutes(
		us.router, us.schoolController,
		us.cfg,
	)

	return us
}

func Test_schoolController_CreateSchool(t *testing.T) {
	tests := []struct {
		name string
		body func() *bytes.Reader
		mock func(ts schoolControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 올바른 학교명",
			body: func() *bytes.Reader {
				req := domain.CreateSchoolRequest{
					Name:   "클래스팅",
					Region: "서울",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts schoolControllerTestSuite) {
				ts.schoolService.EXPECT().CreateSchool(mock.Anything, domain.CreateSchoolRequest{
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "PASS - 빈 문자열 학교명",
			body: func() *bytes.Reader {
				req := domain.CreateSchoolRequest{
					Name:   "",
					Region: "서울",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts schoolControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
		{
			name: "PASS - 빈 문자열 지역",
			body: func() *bytes.Reader {
				req := domain.CreateSchoolRequest{
					Name:   "클래스팅",
					Region: "",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts schoolControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodPost, "/schools", tt.body())
			req.Header.Set("Content-Type", "application/json")
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeAdmin,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.schoolService.AssertExpectations(t)
		})
	}
}
