package subscription

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
	"k8s.io/utils/pointer"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type subscriptionControllerTestSuite struct {
	router                 *gin.Engine
	cfg                    *config.Config
	subscriptionService    *mocks.SubscriptionService
	subscriptionController domain.SubscriptionController
}

func setupSubscriptionControllerTestSuite(t *testing.T) subscriptionControllerTestSuite {
	var us subscriptionControllerTestSuite

	gin.SetMode(gin.TestMode)
	us.router = gin.Default()
	us.subscriptionService = mocks.NewSubscriptionService(t)
	us.cfg = &config.Config{
		Auth: config.Auth{
			Secret: "classting_test_secret",
		},
	}

	us.subscriptionController = NewSubscriptionController(us.subscriptionService)
	RegisterRoutes(
		us.router, us.subscriptionController,
		us.cfg,
	)

	return us
}

func Test_subscriptionController_CreateSubscription(t *testing.T) {
	tests := []struct {
		name string
		body func() *bytes.Reader
		mock func(ts subscriptionControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 올바른 구독 생성",
			body: func() *bytes.Reader {
				req := domain.CreateSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().CreateSubscription(mock.Anything, domain.CreateSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "FAIL - 학교 ID 누락",
			body: func() *bytes.Reader {
				req := domain.CreateSubscriptionRequest{
					UserID:   1,
					SchoolID: 0,
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts subscriptionControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodPost, "/subscriptions", tt.body())
			req.Header.Set("Content-Type", "application/json")
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeStudent,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.subscriptionService.AssertExpectations(t)
		})
	}
}

func Test_subscriptionController_ListSubscriptionSchools(t *testing.T) {
	tests := []struct {
		name  string
		query func() string
		mock  func(ts subscriptionControllerTestSuite)
		code  int
	}{
		{
			name: "PASS - 전체 조회 (커서 미 입력)",
			query: func() string {
				params := url.Values{}
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().ListSubscriptionSchools(mock.Anything, domain.ListSubscriptionSchoolsRequest{
					UserID: 1,
					Cursor: nil,
				}).Return(domain.ListSubscriptionSchoolsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
		{
			name: "PASS - 일부 조회 (커서 입력)",
			query: func() string {
				params := url.Values{}
				params.Add("cursor", "1")
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().ListSubscriptionSchools(mock.Anything, domain.ListSubscriptionSchoolsRequest{
					UserID: 1,
					Cursor: pointer.Int(1),
				}).Return(domain.ListSubscriptionSchoolsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodGet, "/subscriptions", nil)
			req.URL.RawQuery = tt.query()
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeStudent,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			t.Logf("Request URL: %s", req.URL.String())
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.subscriptionService.AssertExpectations(t)
		})
	}
}

func Test_subscriptionController_DeleteSubscription(t *testing.T) {
	tests := []struct {
		name string
		path func() string
		mock func(ts subscriptionControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 구독 취소 성공",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions", "1")
				return path
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().DeleteSubscription(mock.Anything, domain.DeleteSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "PASS - 잘못된 구독 ID 삭제 시도",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions", "-1")
				return path
			},
			mock: func(ts subscriptionControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodDelete, tt.path(), nil)
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeStudent,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.subscriptionService.AssertExpectations(t)
		})
	}
}

func Test_subscriptionController_ListSubscriptionSchoolNews(t *testing.T) {
	tests := []struct {
		name  string
		path  func() string
		query func() string
		mock  func(ts subscriptionControllerTestSuite)
		code  int
	}{
		{
			name: "PASS - 전체 조회 (커서 미 입력)",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions/news", "1")
				return path
			},
			query: func() string {
				params := url.Values{}
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().ListSubscriptionSchoolNews(mock.Anything, domain.ListSubscriptionSchoolNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   nil,
				}).Return(domain.ListSubscriptionSchoolNewsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
		{
			name: "PASS - 전체 조회 (커서 입력)",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions/news", "1")
				return path
			},
			query: func() string {
				params := url.Values{}
				params.Add("cursor", "1")
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {
				ts.subscriptionService.EXPECT().ListSubscriptionSchoolNews(mock.Anything, domain.ListSubscriptionSchoolNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   pointer.Int(1),
				}).Return(domain.ListSubscriptionSchoolNewsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
		{
			name: "PASS - 잘못 된 학교 ID 조회 시도",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions/news", "-1")
				return path
			},
			query: func() string {
				params := url.Values{}
				params.Add("cursor", "1")
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
		{
			name: "PASS - 잘못 된 커서 조회 시도",
			path: func() string {
				path, _ := url.JoinPath("/subscriptions/news", "1")
				return path
			},
			query: func() string {
				params := url.Values{}
				params.Add("cursor", "0")
				return params.Encode()
			},
			mock: func(ts subscriptionControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodGet, tt.path(), nil)
			req.URL.RawQuery = tt.query()
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeStudent,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			t.Logf("Request URL: %s", req.URL.String())
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.subscriptionService.AssertExpectations(t)
		})
	}
}
