package school

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type schoolRepository struct {
	sqlDB *sql.DB
}

func NewSchoolRepository(sqlDB *sql.DB) *schoolRepository {
	return &schoolRepository{
		sqlDB: sqlDB,
	}
}

var _ domain.SchoolRepository = (*schoolRepository)(nil)

func (s schoolRepository) CreateSchool(ctx context.Context, school domain.School) (int, error) {
	const op cerrors.Op = "school/schoolRepository/CreateSchool"

	result, err := s.sqlDB.ExecContext(ctx, createSchoolQuery, school.UserID, school.Name, school.Region)
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return int(userID), nil
}

func (s schoolRepository) ListSchools(ctx context.Context, params domain.ListSchoolsParams) ([]domain.School, error) {
	const op cerrors.Op = "school/schoolRepository/ListSchools"

	var schools []domain.School

	query := fmt.Sprintf(listSchoolQuery,
		params.AndUserID(),
		params.AfterCursor(),
	)

	rows, err := s.sqlDB.QueryContext(ctx, query)
	if err != nil {
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}
	defer rows.Close()

	for rows.Next() {
		var school domain.School
		err := rows.Scan(
			&school.ID,
			&school.UserID,
			&school.Name,
			&school.Region,
		)
		if err != nil {
			return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
		}
		schools = append(schools, school)
	}

	return schools, nil
}

func (s schoolRepository) FindSchoolByNameAndRegion(ctx context.Context, params domain.FindSchoolByNameAndRegionParams) (*domain.School, error) {
	const op cerrors.Op = "school/schollRepository/FindSchoolByNameAndRegion"
	var school domain.School

	findSchoolByNameAndRegion := `SELECT id, user_id, name, region FROM schools WHERE name = ? AND region = ?`
	err := s.sqlDB.QueryRowContext(ctx, findSchoolByNameAndRegion, params.Name, params.Region).
		Scan(&school.ID, &school.UserID, &school.Name, &school.Region)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return &school, nil
}
