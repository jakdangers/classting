package subscription

import (
	"classting/domain"
	"classting/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/utils/pointer"
	"testing"
)

type subscriptionServiceTestSuite struct {
	schoolRepository       *mocks.SchoolRepository
	newsRepository         *mocks.NewsRepository
	subscriptionRepository *mocks.SubscriptionRepository
	service                domain.SubscriptionService
}

func setupSubscriptionServiceTestSuite(t *testing.T) subscriptionServiceTestSuite {
	var us subscriptionServiceTestSuite

	us.schoolRepository = mocks.NewSchoolRepository(t)
	us.newsRepository = mocks.NewNewsRepository(t)
	us.subscriptionRepository = mocks.NewSubscriptionRepository(t)
	us.service = NewSubscriptionService(us.newsRepository, us.schoolRepository, us.subscriptionRepository)

	return us
}

func Test_subscriptionService_CreateSubscription(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.CreateSubscriptionRequest
	}

	var tests = []struct {
		name    string
		args    args
		mock    func(ts subscriptionServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 새로운 구독 생성",
			args: args{
				ctx: context.Background(),
				req: domain.CreateSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByID(mock.Anything, 1).Return(&domain.School{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				}, nil).Once()
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(nil, nil).Once()
				ts.subscriptionRepository.EXPECT().CreateSubscription(mock.Anything, domain.Subscription{
					UserID:   1,
					SchoolID: 1,
				}).Return(1, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "FAIL - 이미 구독한 학교",
			args: args{
				ctx: context.Background(),
				req: domain.CreateSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByID(mock.Anything, 1).Return(&domain.School{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				}, nil).Once()
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(&domain.Subscription{}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.CreateSubscription(tt.args.ctx, tt.args.req)

			// then
			ts.schoolRepository.AssertExpectations(t)
			ts.newsRepository.AssertExpectations(t)
			ts.subscriptionRepository.AssertExpectations(t)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_subscriptionService_ListSubscriptionSchools(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.ListSubscriptionSchoolsRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionServiceTestSuite)
		want    domain.ListSubscriptionSchoolsResponse
		wantErr bool
	}{
		{
			name: "PASS - 전체 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSubscriptionSchoolsRequest{
					UserID: 1,
					Cursor: nil,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().ListSubscriptionSchools(mock.Anything, domain.ListSubscriptionSchoolsParams{
					UserID: 1,
					Cursor: nil,
				}).Return([]domain.SubscriptionSchool{
					{
						Base: domain.Base{
							ID: 1,
						},
						SchoolID: 1,
						Name:     "클래스팅",
						Region:   "서울",
					},
				}, nil).Once()
			},
			want: domain.ListSubscriptionSchoolsResponse{
				SubscriptionSchools: []domain.SubscriptionSchoolDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 1,
						},
						SchoolID: 1,
						Name:     "클래스팅",
						Region:   "서울",
					},
				},
				Cursor: pointer.Int(1),
			},
			wantErr: false,
		},
		{
			name: "PASS - 페이징 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSubscriptionSchoolsRequest{
					UserID: 1,
					Cursor: pointer.Int(1),
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().ListSubscriptionSchools(mock.Anything, domain.ListSubscriptionSchoolsParams{
					UserID: 1,
					Cursor: pointer.Int(1),
				}).Return([]domain.SubscriptionSchool{
					{
						Base: domain.Base{
							ID: 2,
						},
						SchoolID: 1,
						Name:     "클래스팅",
						Region:   "서울",
					},
				}, nil).Once()
			},
			want: domain.ListSubscriptionSchoolsResponse{
				SubscriptionSchools: []domain.SubscriptionSchoolDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 2,
						},
						SchoolID: 1,
						Name:     "클래스팅",
						Region:   "서울",
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
			ts := setupSubscriptionServiceTestSuite(t)
			tt.mock(ts)

			// when
			got, err := ts.service.ListSubscriptionSchools(tt.args.ctx, tt.args.req)

			// then
			ts.newsRepository.AssertExpectations(t)
			ts.schoolRepository.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_subscriptionService_DeleteSubscription(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.DeleteSubscriptionRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 정상 구독 취소",
			args: args{
				ctx: context.Background(),
				req: domain.DeleteSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(&domain.Subscription{
					Base: domain.Base{
						ID: 1,
					},
					UserID:   1,
					SchoolID: 1,
				}, nil).Once()
				ts.subscriptionRepository.EXPECT().DeleteSubscription(mock.Anything, 1).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "PASS - 이미 취소된 구독",
			args: args{
				ctx: context.Background(),
				req: domain.DeleteSubscriptionRequest{
					UserID:   1,
					SchoolID: 1,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(nil, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.DeleteSubscription(tt.args.ctx, tt.args.req)

			// then
			ts.newsRepository.AssertExpectations(t)
			ts.schoolRepository.AssertExpectations(t)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_subscriptionService_ListSubscriptionSchoolNews(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.ListSubscriptionSchoolNewsRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionServiceTestSuite)
		want    domain.ListSubscriptionSchoolNewsResponse
		wantErr bool
	}{
		{
			name: "PASS - 구독한 학교의 소식 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSubscriptionSchoolNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   nil,
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(&domain.Subscription{
					Base: domain.Base{
						ID: 1,
					},
					UserID:   1,
					SchoolID: 1,
				}, nil).Once()
				ts.newsRepository.EXPECT().ListNews(mock.Anything, domain.ListNewsParams{
					SchoolID: pointer.Int(1),
				}).Return([]domain.News{
					{
						Base: domain.Base{
							ID: 1,
						},
						SchoolID: 1,
						UserID:   2,
						Title:    "구독한 뉴스",
					},
				}, nil).Once()
			},
			want: domain.ListSubscriptionSchoolNewsResponse{
				SubscriptionSchoolNews: []domain.SubscriptionSchoolNewsDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 1,
						},
						SchoolID: 1,
						Title:    "구독한 뉴스",
					},
				},
				Cursor: pointer.Int(1),
			},
			wantErr: false,
		},
		{
			name: "PASS - 페이징 구독한 학교의 소식 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSubscriptionSchoolNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   pointer.Int(1),
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(&domain.Subscription{
					Base: domain.Base{
						ID: 1,
					},
					UserID:   1,
					SchoolID: 1,
				}, nil).Once()
				ts.newsRepository.EXPECT().ListNews(mock.Anything, domain.ListNewsParams{
					SchoolID: pointer.Int(1),
					Cursor:   pointer.Int(1),
				}).Return([]domain.News{
					{
						Base: domain.Base{
							ID: 2,
						},
						SchoolID: 1,
						UserID:   2,
						Title:    "구독한 뉴스",
					},
				}, nil).Once()
			},
			want: domain.ListSubscriptionSchoolNewsResponse{
				SubscriptionSchoolNews: []domain.SubscriptionSchoolNewsDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 2,
						},
						SchoolID: 1,
						Title:    "구독한 뉴스",
					},
				},
				Cursor: pointer.Int(2),
			},
			wantErr: false,
		},
		{
			name: "FAIL - 구독하지 않은 학교의 소식 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListSubscriptionSchoolNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   pointer.Int(1),
				},
			},
			mock: func(ts subscriptionServiceTestSuite) {
				ts.subscriptionRepository.EXPECT().FindSubscriptionByUserIDAndSchoolID(mock.Anything, domain.FindSubscriptionByUserIDAndSchoolIDParams{
					UserID:   1,
					SchoolID: 1,
				}).Return(nil, nil).Once()
			},
			want:    domain.ListSubscriptionSchoolNewsResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionServiceTestSuite(t)
			tt.mock(ts)

			// when
			got, err := ts.service.ListSubscriptionSchoolNews(tt.args.ctx, tt.args.req)

			// then
			ts.newsRepository.AssertExpectations(t)
			ts.schoolRepository.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
