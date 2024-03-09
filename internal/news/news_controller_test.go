package news

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

type newsControllerTestSuite struct {
	router         *gin.Engine
	cfg            *config.Config
	newsService    *mocks.NewsService
	newsController domain.NewsController
}

func setupNewsControllerTestSuite(t *testing.T) newsControllerTestSuite {
	var us newsControllerTestSuite

	gin.SetMode(gin.TestMode)
	us.router = gin.Default()
	us.newsService = mocks.NewNewsService(t)
	us.cfg = &config.Config{
		Auth: config.Auth{
			Secret: "classting_test_secret",
		},
	}

	us.newsController = NewNewsController(us.newsService)
	RegisterRoutes(
		us.router, us.newsController,
		us.cfg,
	)

	return us
}

func Test_newsController_CreateNews(t *testing.T) {
	tests := []struct {
		name string
		body func() *bytes.Reader
		mock func(ts newsControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 올바른 소식 생성",
			body: func() *bytes.Reader {
				req := domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Title:    "올바른 소식 생성",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {
				ts.newsService.EXPECT().CreateNews(mock.Anything, domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Title:    "올바른 소식 생성",
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "FAIL - 타이틀 빈 문자열",
			body: func() *bytes.Reader {
				req := domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Title:    "",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
		{
			name: "FAIL - 학교 아이디 누락",
			body: func() *bytes.Reader {
				req := domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 0,
					Title:    "학교 아이디 누락",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodPost, "/news", tt.body())
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
			ts.newsService.AssertExpectations(t)
		})
	}
}

func Test_newsController_ListNews(t *testing.T) {
	tests := []struct {
		name  string
		query func() string
		mock  func(ts newsControllerTestSuite)
		code  int
	}{
		{
			name: "PASS - 전체 조회 (커서 미 입력)",
			query: func() string {
				params := url.Values{}
				params.Add("schoolID", "1")
				return params.Encode()
			},
			mock: func(ts newsControllerTestSuite) {
				ts.newsService.EXPECT().ListNews(mock.Anything, domain.ListNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   nil,
				}).Return(domain.ListNewsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
		{
			name: "PASS - 일부 조회 (커서 입력)",
			query: func() string {
				params := url.Values{}
				params.Add("cursor", "1")
				params.Add("schoolID", "1")
				return params.Encode()
			},
			mock: func(ts newsControllerTestSuite) {
				ts.newsService.EXPECT().ListNews(mock.Anything, domain.ListNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   pointer.Int(1),
				}).Return(domain.ListNewsResponse{}, nil).Once()
			},
			code: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodGet, "/news", nil)
			req.URL.RawQuery = tt.query()
			token, _ := user.CreateAccessToken(domain.User{
				Base: domain.Base{
					ID: 1,
				},
				Type: domain.UserUseTypeAdmin,
			}, ts.cfg.Auth.Secret, time.Now().UTC().Add(time.Hour*time.Duration(24)))
			req.Header.Set("Authorization", "Bearer "+token)

			// when
			rec := httptest.NewRecorder()
			t.Logf("Request URL: %s", req.URL.String())
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.newsService.AssertExpectations(t)
		})
	}
}

func Test_newsController_UpdateNews(t *testing.T) {
	tests := []struct {
		name string
		body func() *bytes.Reader
		mock func(ts newsControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 올바른 소식 수정",
			body: func() *bytes.Reader {
				req := domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "올바른 소식 수정",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {
				ts.newsService.EXPECT().UpdateNews(mock.Anything, domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "올바른 소식 수정",
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "FAIL - 빈 타이틀",
			body: func() *bytes.Reader {
				req := domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
		{
			name: "FAIL - 잘못 된 ID",
			body: func() *bytes.Reader {
				req := domain.UpdateNewsRequest{
					UserID: 1,
					ID:     0,
					Title:  "잘못 된 ID",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts newsControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodPut, "/news", tt.body())
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
			ts.newsService.AssertExpectations(t)
		})
	}
}

func Test_newsController_DeleteNews(t *testing.T) {
	tests := []struct {
		name string
		path func() string
		mock func(ts newsControllerTestSuite)
		code int
	}{
		{
			name: "PASS - 소식 삭제 성공",
			path: func() string {
				path, _ := url.JoinPath("/news", "100")
				return path
			},
			mock: func(ts newsControllerTestSuite) {
				ts.newsService.EXPECT().DeleteNews(mock.Anything, domain.DeleteNewsRequest{
					UserID: 1,
					ID:     100,
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "PASS - 잘못된 소식 ID 삭제 시도",
			path: func() string {
				path, _ := url.JoinPath("/news", "-1")
				return path
			},
			mock: func(ts newsControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodDelete, tt.path(), nil)
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
			ts.newsService.AssertExpectations(t)
		})
	}
}
