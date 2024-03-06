package user

import (
	"classting/domain"
	"classting/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type userServiceTestSuite struct {
	userRepository *mocks.UserRepository
	service        domain.UserService
}

func setupUserServiceTestSuite(t *testing.T) userServiceTestSuite {
	var us userServiceTestSuite

	us.userRepository = mocks.NewUserRepository(t)
	us.service = NewUserService(us.userRepository)

	return us
}

func Test_userService_CreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.CreateUserRequest
	}

	var tests = []struct {
		name    string
		args    args
		mock    func(ts userServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 중복되지 않는 유저아이디",
			args: args{
				ctx: context.Background(),
				req: domain.CreateUserRequest{
					UserName: "unique_user_name",
					Password: "classting",
					UserType: domain.UserUseTypeStudent,
				},
			},
			mock: func(ts userServiceTestSuite) {
				ts.userRepository.EXPECT().FindUserByUserName(mock.Anything, "unique_user_name").Return(nil, nil).Once()
				ts.userRepository.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
					return user.UserName == "unique_user_name" && compareHashAndPassword("classting", user.Password)
				})).Return(1, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "FAIL - 중복된 유저 아이디",
			args: args{
				ctx: context.Background(),
				req: domain.CreateUserRequest{
					UserName: "duplicate_user_name",
					Password: "classting",
					UserType: domain.UserUseTypeStudent,
				},
			},
			mock: func(ts userServiceTestSuite) {
				ts.userRepository.EXPECT().FindUserByUserName(mock.Anything, "duplicate_user_name").Return(&domain.User{}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupUserServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.CreateUser(tt.args.ctx, tt.args.req)

			// then
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_validateUserName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "PASS - 유효한 유저아이디", input: "exist_user_name", want: "exist_user_name", wantErr: false},
		{name: "FAIL - 빈 문자열 유저아이디", input: "", want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateUserName(tt.input)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
