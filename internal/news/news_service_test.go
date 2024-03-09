package news

import (
	"classting/domain"
	"classting/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/utils/pointer"
	"testing"
)

type newsServiceTestSuite struct {
	schoolRepository *mocks.SchoolRepository
	newsRepository   *mocks.NewsRepository
	service          domain.NewsService
}

func setupNewsServiceTestSuite(t *testing.T) newsServiceTestSuite {
	var us newsServiceTestSuite

	us.schoolRepository = mocks.NewSchoolRepository(t)
	us.newsRepository = mocks.NewNewsRepository(t)
	us.service = NewNewsService(us.newsRepository, us.schoolRepository)

	return us
}

func Test_newsService_CreateNews(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.CreateNewsRequest
	}

	var tests = []struct {
		name    string
		args    args
		mock    func(ts newsServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 새로운 소식 발행",
			args: args{
				ctx: context.Background(),
				req: domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Title:    "클래스팅 소식",
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByID(mock.Anything, 1).Return(&domain.School{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				}, nil).Once()
				ts.newsRepository.EXPECT().CreateNews(mock.Anything, domain.News{
					SchoolID: 1,
					UserID:   1,
					Title:    "클래스팅 소식",
				}).Return(1, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "FAIL - 타인의 학교 소식 발행",
			args: args{
				ctx: context.Background(),
				req: domain.CreateNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Title:    "클래스팅 소식",
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.schoolRepository.EXPECT().FindSchoolByID(mock.Anything, 1).Return(&domain.School{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 2,
					Name:   "클래스팅",
					Region: "서울",
				}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.CreateNews(tt.args.ctx, tt.args.req)

			// then
			ts.schoolRepository.AssertExpectations(t)
			ts.newsRepository.AssertExpectations(t)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsService_ListNews(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.ListNewsRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsServiceTestSuite)
		want    domain.ListNewsResponse
		wantErr bool
	}{
		{
			name: "PASS - 전체 조회",
			args: args{
				ctx: context.Background(),
				req: domain.ListNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   nil,
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().ListNews(mock.Anything, domain.ListNewsParams{
					UserID:   pointer.Int(1),
					SchoolID: pointer.Int(1),
					Cursor:   nil,
				}).Return([]domain.News{
					{
						Base: domain.Base{
							ID: 1,
						},
						SchoolID: 1,
						UserID:   1,
						Title:    "클래스팅 소식",
					},
				}, nil).Once()
			},
			want: domain.ListNewsResponse{
				News: []domain.NewsDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 1,
						},
						SchoolID: 1,
						Title:    "클래스팅 소식",
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
				req: domain.ListNewsRequest{
					UserID:   1,
					SchoolID: 1,
					Cursor:   pointer.Int(1),
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().ListNews(mock.Anything, domain.ListNewsParams{
					UserID:   pointer.Int(1),
					SchoolID: pointer.Int(1),
					Cursor:   pointer.Int(1),
				}).Return([]domain.News{
					{
						Base: domain.Base{
							ID: 2,
						},
						SchoolID: 1,
						UserID:   1,
						Title:    "클래스팅 소식",
					},
				}, nil).Once()
			},
			want: domain.ListNewsResponse{
				News: []domain.NewsDTO{
					{
						BaseDTO: domain.BaseDTO{
							ID: 2,
						},
						SchoolID: 1,
						Title:    "클래스팅 소식",
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
			ts := setupNewsServiceTestSuite(t)
			tt.mock(ts)

			// when
			got, err := ts.service.ListNews(tt.args.ctx, tt.args.req)

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

func Test_newsService_UpdateNews(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.UpdateNewsRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 정상 수정",
			args: args{
				ctx: context.Background(),
				req: domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "타이틀 수정",
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(&domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "타이틀 원본",
				}, nil).Once()
				ts.newsRepository.EXPECT().UpdateNews(mock.Anything, domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "타이틀 수정",
				}).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "PASS - 요청 된 소식이 없음",
			args: args{
				ctx: context.Background(),
				req: domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "타이틀 수정",
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(nil, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "PASS - 요청 된 소식의 작성자가 아님",
			args: args{
				ctx: context.Background(),
				req: domain.UpdateNewsRequest{
					UserID: 1,
					ID:     1,
					Title:  "타이틀 수정",
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(&domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   7777,
					Title:    "소유주가 다름",
				}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.UpdateNews(tt.args.ctx, tt.args.req)

			// then
			ts.newsRepository.AssertExpectations(t)
			ts.schoolRepository.AssertExpectations(t)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsService_DeleteNews(t *testing.T) {
	type args struct {
		ctx context.Context
		req domain.DeleteNewsRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsServiceTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 정상 삭제",
			args: args{
				ctx: context.Background(),
				req: domain.DeleteNewsRequest{
					UserID: 1,
					ID:     1,
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(&domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "삭제할 소식",
				}, nil).Once()
				ts.newsRepository.EXPECT().DeleteNews(mock.Anything, 1).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "FAIL - 요청 된 소식이 없음",
			args: args{
				ctx: context.Background(),
				req: domain.DeleteNewsRequest{
					UserID: 1,
					ID:     1,
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(nil, nil).Once()
			},
			wantErr: true,
		},
		{
			name: "FAIL - 요청 된 소식의 작성자가 아님",
			args: args{
				ctx: context.Background(),
				req: domain.DeleteNewsRequest{
					UserID: 1,
					ID:     1,
				},
			},
			mock: func(ts newsServiceTestSuite) {
				ts.newsRepository.EXPECT().FindNewsByID(mock.Anything, 1).Return(&domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   7777,
					Title:    "소식의 작성자가 달라요",
				}, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsServiceTestSuite(t)
			tt.mock(ts)

			// when
			err := ts.service.DeleteNews(tt.args.ctx, tt.args.req)

			// then
			ts.newsRepository.AssertExpectations(t)
			ts.schoolRepository.AssertExpectations(t)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
