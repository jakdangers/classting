package school

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
)

type schoolService struct {
	userRepository   domain.UserRepository
	schoolRepository domain.SchoolRepository
}

func NewSchoolService(
	schoolRepository domain.SchoolRepository,
	userRepository domain.UserRepository,
) *schoolService {
	return &schoolService{
		userRepository:   userRepository,
		schoolRepository: schoolRepository,
	}
}

var _ domain.SchoolService = (*schoolService)(nil)

func (s schoolService) CreateSchool(ctx context.Context, req domain.CreateSchoolRequest) error {
	const op cerrors.Op = "school/service/createSchool"

	school, err := s.schoolRepository.FindSchoolByNameAndRegion(ctx, domain.FindSchoolByNameAndRegionParams{
		Name:   req.Name,
		Region: req.Region,
	})
	if err != nil {
		return err
	}
	if school != nil {
		return cerrors.E(op, cerrors.Invalid, "이미 사용중인 지역, 학교명입니다.")
	}

	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	_, err = s.schoolRepository.CreateSchool(ctx, domain.School{
		UserID: req.UserID,
		Name:   req.Name,
		Region: req.Region,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s schoolService) ListSchools(ctx context.Context, req domain.ListSchoolsRequest) (domain.ListSchoolsResponse, error) {
	const op cerrors.Op = "school/service/ListSchools"

	schools, err := s.schoolRepository.ListSchools(ctx, domain.ListSchoolsParams{
		UserID: req.UserID,
		Cursor: req.Cursor,
	})
	if err != nil {
		return domain.ListSchoolsResponse{}, cerrors.E(op, cerrors.Internal, err, "학교를 조회하는 중에 에러가 발생했습니다.")
	}

	var schoolDTOs []domain.SchoolDTO
	for _, school := range schools {
		schoolDTOs = append(schoolDTOs, domain.SchoolDTOFrom(school))
	}

	var cursor *int
	if len(schoolDTOs) > 0 {
		cursor = &schoolDTOs[len(schoolDTOs)-1].ID
	}

	return domain.ListSchoolsResponse{
		Schools: schoolDTOs,
		Cursor:  cursor,
	}, nil
}
