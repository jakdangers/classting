package school

import (
	"classting/domain"
	"classting/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type schoolServiceTestSuite struct {
	schoolRepository *mocks.SchoolRepository
	userRepository   *mocks.UserRepository
	service          domain.SchoolService
}

func setupSchoolServiceTestSuite(t *testing.T) schoolServiceTestSuite {
	var us schoolServiceTestSuite

	us.schoolRepository = mocks.NewSchoolRepository(t)
	us.userRepository = mocks.NewUserRepository(t)
	us.service = NewSchoolService(us.schoolRepository, us.userRepository)

	return us
}

func Test_schoolService_CreateSchool(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.CreateSchoolRequest
	}

	var tests = []struct {
		name    string
		args    args
		mock    func(ts schoolServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 중복되지 않는 학교",
			args: args{
				ctx: context.Background(),
				req: domain.CreateSchoolRequest{
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				},
			},
			mock: func(ts schoolServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByNameAndRegion(mock.Anything, domain.FindSchoolByNameAndRegionParams{
					Name:   "클래스팅",
					Region: "서울",
				}).Return(nil, nil).Once()
				ts.schoolRepository.EXPECT().CreateSchool(mock.Anything, domain.School{
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				}).Return(1, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "FAIL - 중복된 학교 생성",
			args: args{
				ctx: context.Background(),
				req: domain.CreateSchoolRequest{
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				},
			},
			mock: func(ts schoolServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByNameAndRegion(mock.Anything, domain.FindSchoolByNameAndRegionParams{
					Name:   "클래스팅",
					Region: "서울",
				}).Return(&domain.School{}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.CreateSchool(tt.args.ctx, tt.args.req)

			// then
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
