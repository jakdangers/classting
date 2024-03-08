package school

import (
	"classting/domain"
	"classting/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/utils/pointer"
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

func Test_schoolService_ListSchools(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.ListSchoolsRequest
	}

	var tests = []struct {
		name    string
		args    args
		mock    func(ts schoolServiceTestSuite)
		want    domain.ListSchoolsResponse
		wantErr bool
	}{
		{
			name: "PASS - 전체 학교 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSchoolsRequest{
					UserID: nil,
					Cursor: nil,
				},
			},
			mock: func(ts schoolServiceTestSuite) {
				ts.schoolRepository.EXPECT().ListSchools(mock.Anything, domain.ListSchoolsParams{
					UserID: nil,
					Cursor: nil,
				}).Return([]domain.School{
					{
						Base: domain.Base{
							ID: 1,
						},
						UserID: 1,
						Name:   "클래스팅",
						Region: "서울",
					},
				}, nil).Once()
			},
			want: domain.ListSchoolsResponse{
				Schools: []domain.SchoolDTO{
					{
						ID:     1,
						Name:   "클래스팅",
						Region: "서울",
					},
				},
				Cursor: pointer.Int(1),
			},
			wantErr: false,
		},
		{
			name: "PASS - 특정 관리자의 학교 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSchoolsRequest{
					UserID: pointer.Int(1),
					Cursor: nil,
				},
			},
			mock: func(ts schoolServiceTestSuite) {
				ts.schoolRepository.EXPECT().ListSchools(mock.Anything, domain.ListSchoolsParams{
					UserID: pointer.Int(1),
					Cursor: nil,
				}).Return([]domain.School{
					{
						Base: domain.Base{
							ID: 1,
						},
						UserID: 1,
						Name:   "클래스팅",
						Region: "서울",
					},
				}, nil).Once()
			},
			want: domain.ListSchoolsResponse{
				Schools: []domain.SchoolDTO{
					{
						ID:     1,
						Name:   "클래스팅",
						Region: "서울",
					},
				},
				Cursor: pointer.Int(1),
			},
			wantErr: false,
		},
		{
			name: "PASS - 특정 커서 이상의 학교 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSchoolsRequest{
					UserID: nil,
					Cursor: pointer.Int(1),
				},
			},
			mock: func(ts schoolServiceTestSuite) {
				ts.schoolRepository.EXPECT().ListSchools(mock.Anything, domain.ListSchoolsParams{
					UserID: nil,
					Cursor: pointer.Int(1),
				}).Return([]domain.School{
					{
						Base: domain.Base{
							ID: 2,
						},
						UserID: 1,
						Name:   "클래스팅",
						Region: "서울",
					},
				}, nil).Once()
			},
			want: domain.ListSchoolsResponse{
				Schools: []domain.SchoolDTO{
					{
						ID:     2,
						Name:   "클래스팅",
						Region: "서울",
					},
				},
				Cursor: pointer.Int(2),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolServiceTestSuite(t)
			tt.mock(ts)

			// when
			got, err := ts.service.ListSchools(tt.args.ctx, tt.args.req)

			// then
			ts.schoolRepository.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
