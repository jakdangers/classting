package user

import (
	"bytes"
	"classting/domain"
	"classting/mocks"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userControllerTestSuite struct {
	router         *gin.Engine
	userService    *mocks.UserService
	userController domain.UserController
}

func setupUserControllerTestSuite(t *testing.T) userControllerTestSuite {
	var us userControllerTestSuite

	gin.SetMode(gin.TestMode)
	us.router = gin.Default()
	us.userService = mocks.NewUserService(t)

	us.userController = NewUserController(us.userService)
	RegisterRoutes(
		us.router, us.userController,
	)

	return us
}

func Test_userController_CreateUser(t *testing.T) {
	tests := []struct {
		name  string
		input func() *bytes.Reader
		mock  func(ts userControllerTestSuite)
		code  int
	}{
		{
			name: "PASS - 올바른 유저네임",
			input: func() *bytes.Reader {
				req := domain.CreateUserRequest{
					UserName: "classting_user",
					Password: "classting",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts userControllerTestSuite) {
				ts.userService.EXPECT().CreateUser(mock.Anything, domain.CreateUserRequest{
					UserName: "classting_user",
					Password: "classting",
				}).Return(nil).Once()
			},
			code: http.StatusNoContent,
		},
		{
			name: "FAIL - 유저네임 빈 문자열",
			input: func() *bytes.Reader {
				req := domain.CreateUserRequest{
					UserName: "",
					Password: "classting",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts userControllerTestSuite) {
			},
			code: http.StatusBadRequest,
		},
		{
			name: "FAIL – 패스워드 빈 문자열",
			input: func() *bytes.Reader {
				req := domain.CreateUserRequest{
					UserName: "classting_user",
					Password: "",
				}
				jsonData, _ := json.Marshal(req)

				return bytes.NewReader(jsonData)
			},
			mock: func(ts userControllerTestSuite) {},
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupUserControllerTestSuite(t)
			tt.mock(ts)
			req, _ := http.NewRequest(http.MethodPost, "/users", tt.input())
			req.Header.Set("Content-Type", "application/json")

			// when
			rec := httptest.NewRecorder()
			ts.router.ServeHTTP(rec, req)

			// then
			assert.Equal(t, tt.code, rec.Code)
			ts.userService.AssertExpectations(t)
		})
	}
}
