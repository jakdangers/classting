package subscription

import (
	"classting/domain"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type subscriptionRepositoryTestSuite struct {
	sqlDB                  *sql.DB
	sqlMock                sqlmock.Sqlmock
	subscriptionRepository domain.SubscriptionRepository
}

func setupSubscriptionRepositoryTestSuite() subscriptionRepositoryTestSuite {
	var us subscriptionRepositoryTestSuite

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	us.sqlDB = mockDB
	us.sqlMock = mock

	if err != nil {
		panic(err)
	}

	us.subscriptionRepository = NewSubscriptionRepository(mockDB)

	return us
}

func Test_subscriptionRepository_CreateSubscription(t *testing.T) {
	type args struct {
		ctx  context.Context
		news domain.Subscription
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionRepositoryTestSuite)
		want    int
		wantErr bool
	}{
		{
			name: "PASS - 학교 소식 생성",
			args: args{
				ctx: context.Background(),
				news: domain.Subscription{
					SchoolID: 1,
					UserID:   1,
				},
			},
			mock: func(ts subscriptionRepositoryTestSuite) {
				ts.sqlMock.ExpectExec(`INSERT INTO subscriptions`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionRepositoryTestSuite()
			tt.mock(ts)

			// when
			got, err := ts.subscriptionRepository.CreateSubscription(tt.args.ctx, tt.args.news)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_subscriptionRepository_ListSubscriptionSchools(t *testing.T) {
	type args struct {
		ctx    context.Context
		params domain.ListSubscriptionSchoolsParams
	}

	createDate := time.Now()
	updateDate := time.Now()

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionRepositoryTestSuite)
		want    []domain.SubscriptionSchool
		wantErr bool
	}{
		{
			name: "PASS - 조건 없는 구독 조회",
			args: args{
				ctx: context.Background(),
				params: domain.ListSubscriptionSchoolsParams{
					UserID: 1,
				},
			},
			mock: func(ts subscriptionRepositoryTestSuite) {
				query := `SELECT (.+) FROM schools JOIN subscriptions ON schools.id = subscriptions.school_id`
				columns := []string{"subscriptions.id", "subscriptions.create_date", "subscriptions.update_date", "schools.school_id", "schools.name", "schools.region"}
				rows := sqlmock.NewRows(columns).AddRow(1, createDate, updateDate, 1, "클래스팅 초등학교", "서울")
				ts.sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want: []domain.SubscriptionSchool{
				{
					Base: domain.Base{
						ID:         1,
						CreateDate: createDate,
						UpdateDate: updateDate,
					},
					SchoolID: 1,
					Name:     "클래스팅 초등학교",
					Region:   "서울",
				},
			},
			wantErr: false,
		},
		{
			name: "PASS - 페이징 구독 조회",
			args: args{
				ctx: context.Background(),
				params: domain.ListSubscriptionSchoolsParams{
					UserID: 1,
				},
			},
			mock: func(ts subscriptionRepositoryTestSuite) {
				query := `SELECT (.+) FROM schools JOIN subscriptions ON schools.id = subscriptions.school_id`
				columns := []string{"subscriptions.id", "subscriptions.create_date", "subscriptions.update_date", "schools.school_id", "schools.name", "schools.region"}
				rows := sqlmock.NewRows(columns).AddRow(1, createDate, updateDate, 1, "클래스팅 초등학교", "서울")
				ts.sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want: []domain.SubscriptionSchool{
				{
					Base: domain.Base{
						ID:         1,
						CreateDate: createDate,
						UpdateDate: updateDate,
					},
					SchoolID: 1,
					Name:     "클래스팅 초등학교",
					Region:   "서울",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionRepositoryTestSuite()
			tt.mock(ts)

			// whenR
			got, err := ts.subscriptionRepository.ListSubscriptionSchools(tt.args.ctx, tt.args.params)

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

func Test_subscriptionRepository_DeleteSubscription(t *testing.T) {
	type args struct {
		ctx            context.Context
		subscriptionID int
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts subscriptionRepositoryTestSuite)
		wantErr bool
	}{
		{
			name: "PASS - 구독 취소",
			args: args{
				ctx:            context.Background(),
				subscriptionID: 1,
			},
			mock: func(ts subscriptionRepositoryTestSuite) {
				ts.sqlMock.ExpectExec("DELETE FROM subscriptions").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSubscriptionRepositoryTestSuite()
			tt.mock(ts)

			// whenR
			err := ts.subscriptionRepository.DeleteSubscription(tt.args.ctx, tt.args.subscriptionID)

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
