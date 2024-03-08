package news

import (
	"classting/domain"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
	"testing"
	"time"
)

type newsRepositoryTestSuite struct {
	sqlDB          *sql.DB
	sqlMock        sqlmock.Sqlmock
	newsRepository domain.NewsRepository
}

func setupNewsRepositoryTestSuite() newsRepositoryTestSuite {
	var us newsRepositoryTestSuite

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	us.sqlDB = mockDB
	us.sqlMock = mock

	if err != nil {
		panic(err)
	}

	us.newsRepository = NewNewsRepository(mockDB)

	return us
}

func Test_newsRepository_CreateNews(t *testing.T) {
	type args struct {
		ctx  context.Context
		news domain.News
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsRepositoryTestSuite)
		want    int
		wantErr bool
	}{
		{
			name: "PASS - 학교 소식 생성",
			args: args{
				ctx: context.Background(),
				news: domain.News{
					SchoolID: 1,
					UserID:   1,
					Title:    "클래스팅 새소식",
				},
			},
			mock: func(ts newsRepositoryTestSuite) {
				ts.sqlMock.ExpectExec(`INSERT INTO news`).
					WithArgs(1, 1, "클래스팅 새소식").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsRepositoryTestSuite()
			tt.mock(ts)

			// when
			got, err := ts.newsRepository.CreateNews(tt.args.ctx, tt.args.news)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsRepository_ListNews(t *testing.T) {
	type args struct {
		ctx    context.Context
		params domain.ListNewsParams
	}

	createDate := time.Now()
	updateDate := time.Now()

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsRepositoryTestSuite)
		want    []domain.News
		wantErr bool
	}{
		{
			name: "PASS - 조건 없는 소식 조회 성공",
			args: args{
				ctx: context.Background(),
				params: domain.ListNewsParams{
					UserID: 1,
				},
			},
			mock: func(ts newsRepositoryTestSuite) {
				query := `SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news`
				columns := []string{"id", "create_date", "update_date", "delete_date", "school_id", "user_id", "title"}
				rows := sqlmock.NewRows(columns).AddRow(100, createDate, updateDate, nil, 1, 1, "클래스팅 새소식")
				ts.sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want: []domain.News{
				{
					Base: domain.Base{
						ID:         100,
						CreateDate: createDate,
						UpdateDate: updateDate,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "클래스팅 새소식",
				},
			},
			wantErr: false,
		},
		{
			name: "PASS - 페이징 소식 조회 성공",
			args: args{
				ctx: context.Background(),
				params: domain.ListNewsParams{
					UserID: 1,
					Cursor: pointer.Int(1),
				},
			},
			mock: func(ts newsRepositoryTestSuite) {
				query := `SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news`
				columns := []string{"id", "create_date", "update_date", "delete_date", "school_id", "user_id", "title"}
				rows := sqlmock.NewRows(columns).AddRow(100, createDate, updateDate, nil, 1, 1, "클래스팅 새소식")
				ts.sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want: []domain.News{
				{
					Base: domain.Base{
						ID:         100,
						CreateDate: createDate,
						UpdateDate: updateDate,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "클래스팅 새소식",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsRepositoryTestSuite()
			tt.mock(ts)

			// whenR
			got, err := ts.newsRepository.ListNews(tt.args.ctx, tt.args.params)

			// then
			assert.Equal(t, tt.want, got)
			if ts.sqlMock.ExpectationsWereMet() != nil {
				t.Errorf("there were unfulfilled expectations: %s", ts.sqlMock.ExpectationsWereMet())
			}
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsRepository_FindNewsByID(t *testing.T) {
	type args struct {
		ctx    context.Context
		newsID int
	}

	createDate := time.Now()
	updateDate := time.Now()

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsRepositoryTestSuite)
		want    *domain.News
		wantErr bool
	}{
		{
			name: "PASS - 존재하는 소식 조회",
			args: args{
				ctx:    context.Background(),
				newsID: 1,
			},
			mock: func(ts newsRepositoryTestSuite) {
				query := "SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news"
				columns := []string{"id", "create_date", "update_date", "delete_date", "school_id", "user_id", "title"}
				rows := sqlmock.NewRows(columns).AddRow(1, createDate, updateDate, nil, 1, 1, "클래스팅 소식")
				ts.sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want: &domain.News{
				Base: domain.Base{
					ID:         1,
					CreateDate: createDate,
					UpdateDate: updateDate,
				},
				SchoolID: 1,
				UserID:   1,
				Title:    "클래스팅 소식",
			},
			wantErr: false,
		},
		{
			name: "PASS - 존재하지 소식 조회",
			args: args{
				ctx:    context.Background(),
				newsID: 7777,
			},
			mock: func(ts newsRepositoryTestSuite) {
				query := "SELECT id, create_date, update_date, delete_date, school_id, user_id, title FROM news"
				ts.sqlMock.ExpectQuery(query).WithArgs(7777).WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsRepositoryTestSuite()
			tt.mock(ts)

			// when
			got, err := ts.newsRepository.FindNewsByID(tt.args.ctx, tt.args.newsID)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsRepository_UpdateNews(t *testing.T) {
	type args struct {
		ctx  context.Context
		news domain.News
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsRepositoryTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 소식 수정",
			args: args{
				ctx: context.Background(),
				news: domain.News{
					Base: domain.Base{
						ID: 1,
					},
					SchoolID: 1,
					UserID:   1,
					Title:    "소식 수정",
				},
			},
			mock: func(ts newsRepositoryTestSuite) {
				ts.sqlMock.ExpectExec("UPDATE news").
					WithArgs("소식 수정", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsRepositoryTestSuite()
			tt.mock(ts)

			// when
			err := ts.newsRepository.UpdateNews(tt.args.ctx, tt.args.news)

			// then
			if ts.sqlMock.ExpectationsWereMet() != nil {
				t.Errorf("there were unfulfilled expectations: %s", ts.sqlMock.ExpectationsWereMet())
			}
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_newsRepository_DeleteNews(t *testing.T) {
	type args struct {
		ctx    context.Context
		newsID int
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts newsRepositoryTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 소식 삭제",
			args: args{
				ctx:    context.Background(),
				newsID: 1,
			},
			mock: func(ts newsRepositoryTestSuite) {
				ts.sqlMock.ExpectExec("UPDATE news").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupNewsRepositoryTestSuite()
			tt.mock(ts)

			// when
			err := ts.newsRepository.DeleteNews(tt.args.ctx, tt.args.newsID)

			// then
			if ts.sqlMock.ExpectationsWereMet() != nil {
				t.Errorf("there were unfulfilled expectations: %s", ts.sqlMock.ExpectationsWereMet())
			}
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}
