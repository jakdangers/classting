package school

import (
	"classting/domain"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
	"testing"
)

type schoolRepositoryTestSuite struct {
	sqlDB            *sql.DB
	sqlMock          sqlmock.Sqlmock
	schoolRepository domain.SchoolRepository
}

func setupSchoolRepositoryTestSuite() schoolRepositoryTestSuite {
	var us schoolRepositoryTestSuite

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	us.sqlDB = mockDB
	us.sqlMock = mock

	if err != nil {
		panic(err)
	}

	us.schoolRepository = NewSchoolRepository(mockDB)

	return us
}

func Test_schoolRepository_CreateSchool(t *testing.T) {
	type args struct {
		ctx    context.Context
		school domain.School
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts schoolRepositoryTestSuite)
		want    int
		wantErr bool
	}{
		{
			name: "PASS - 학교 페이지 생성",
			args: args{
				ctx: context.Background(),
				school: domain.School{
					UserID: 1,
					Name:   "클래스팅학교",
					Region: "서울",
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				ts.sqlMock.ExpectExec("INSERT INTO schools").
					WithArgs(1, "클래스팅학교", "서울").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "PASS - 중복 된 학교 생성",
			args: args{
				ctx: context.Background(),
				school: domain.School{
					UserID: 1,
					Name:   "클래스팅학교",
					Region: "서울",
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				ts.sqlMock.ExpectExec("INSERT INTO schools").
					WithArgs(1, "클래스팅학교", "서울").
					WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolRepositoryTestSuite()
			tt.mock(ts)

			// when
			got, err := ts.schoolRepository.CreateSchool(tt.args.ctx, tt.args.school)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_schoolRepository_FindSchoolByNameAndRegion(t *testing.T) {
	type args struct {
		ctx    context.Context
		params domain.FindSchoolByNameAndRegionParams
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts schoolRepositoryTestSuite)
		want    *domain.School
		wantErr bool
	}{
		{
			name: "PASS - 존재하는 지역, 학교명으로 조회",
			args: args{
				ctx: context.Background(),
				params: domain.FindSchoolByNameAndRegionParams{
					Name:   "클래스팅학교",
					Region: "서울",
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				query := "SELECT id, user_id, name, region FROM schools"
				columns := []string{"id", "user_id", "name", "region"}
				rows := sqlmock.NewRows(columns).AddRow(1, 1, "클래스팅학교", "서울")
				ts.sqlMock.ExpectQuery(query).WithArgs("클래스팅학교", "서울").WillReturnRows(rows)
			},
			want: &domain.School{
				Base: domain.Base{
					ID: 1,
				},
				UserID: 1,
				Name:   "클래스팅학교",
				Region: "서울",
			},
			wantErr: false,
		},
		{
			name: "PASS - 존재하지 않는 지역, 학교명으로 조회",
			args: args{
				ctx: context.Background(),
				params: domain.FindSchoolByNameAndRegionParams{
					Name:   "클래스팅학교",
					Region: "서울",
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				query := "SELECT id, user_id, name, region FROM schools"
				ts.sqlMock.ExpectQuery(query).WithArgs("클래스팅학교", "서울").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolRepositoryTestSuite()
			tt.mock(ts)

			// when
			got, err := ts.schoolRepository.FindSchoolByNameAndRegion(tt.args.ctx, tt.args.params)

			// then
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Equalf(t, tt.wantErr, err != nil, err.Error())
			}
		})
	}
}

func Test_schoolRepository_ListSchools(t *testing.T) {
	type args struct {
		ctx    context.Context
		params domain.ListSchoolsParams
	}

	tests := []struct {
		name    string
		args    args
		mock    func(ts schoolRepositoryTestSuite)
		want    []domain.School
		wantErr bool
	}{
		{
			name: "PASS - 전체 조회",
			args: args{
				ctx:    context.Background(),
				params: domain.ListSchoolsParams{},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				query := `SELECT id, user_id, name, region FROM schools`
				columns := []string{"id", "user_id", "name", "region"}
				rows := sqlmock.NewRows(columns).AddRow(1, 1, "클래스팅", "서울")
				ts.sqlMock.ExpectQuery(query).WillReturnRows(rows)
			},
			want: []domain.School{
				{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				},
			},
			wantErr: false,
		},
		{
			name: "PASS - 전체 조회 (특정 유저 아이디 입력)",
			args: args{
				ctx: context.Background(),
				params: domain.ListSchoolsParams{
					UserID: pointer.Int(1),
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				query := `SELECT id, user_id, name, region FROM schools`
				columns := []string{"id", "user_id", "name", "region"}
				rows := sqlmock.NewRows(columns).AddRow(1, 1, "클래스팅", "서울")
				ts.sqlMock.ExpectQuery(query).WillReturnRows(rows)
			},
			want: []domain.School{
				{
					Base: domain.Base{
						ID: 1,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				},
			},
			wantErr: false,
		},
		{
			name: "PASS - 전체 조회 (특정 유저 아이디 입력, 커서 입력)",
			args: args{
				ctx: context.Background(),
				params: domain.ListSchoolsParams{
					UserID: pointer.Int(1),
					Cursor: pointer.Int(1),
				},
			},
			mock: func(ts schoolRepositoryTestSuite) {
				query := `SELECT id, user_id, name, region FROM schools`
				columns := []string{"id", "user_id", "name", "region"}
				rows := sqlmock.NewRows(columns).AddRow(2, 1, "클래스팅", "서울")
				ts.sqlMock.ExpectQuery(query).WillReturnRows(rows)
			},
			want: []domain.School{
				{
					Base: domain.Base{
						ID: 2,
					},
					UserID: 1,
					Name:   "클래스팅",
					Region: "서울",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ts := setupSchoolRepositoryTestSuite()
			tt.mock(ts)

			// whenR
			got, err := ts.schoolRepository.ListSchools(tt.args.ctx, tt.args.params)

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
