package news

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
	"k8s.io/utils/pointer"
)

type newsService struct {
	newsRepository   domain.NewsRepository
	schoolRepository domain.SchoolRepository
}

func NewNewsService(
	newsRepository domain.NewsRepository,
	schoolRepository domain.SchoolRepository,
) *newsService {
	return &newsService{
		newsRepository:   newsRepository,
		schoolRepository: schoolRepository,
	}
}

var _ domain.NewsService = (*newsService)(nil)

func (s newsService) CreateNews(ctx context.Context, req domain.CreateNewsRequest) error {
	const op cerrors.Op = "news/service/createNews"

	school, err := s.schoolRepository.FindSchoolByID(ctx, req.SchoolID)
	if err != nil {
		return err
	}
	if school == nil {
		return cerrors.E(op, cerrors.Invalid, "해당 학교가 존재하지 않습니다.")
	}
	if school.UserID != req.UserID {
		return cerrors.E(op, cerrors.Permission, "해당 학교에 대한 권한이 없습니다.")
	}

	_, err = s.newsRepository.CreateNews(ctx, domain.News{
		SchoolID: req.SchoolID,
		UserID:   req.UserID,
		Title:    req.Title,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s newsService) ListNews(ctx context.Context, req domain.ListNewsRequest) (domain.ListNewsResponse, error) {
	const op cerrors.Op = "news/service/ListNews"

	news, err := s.newsRepository.ListNews(ctx, domain.ListNewsParams{
		UserID:   pointer.Int(req.UserID),
		SchoolID: pointer.Int(req.SchoolID),
		Cursor:   req.Cursor,
	})
	if err != nil {
		return domain.ListNewsResponse{}, cerrors.E(op, cerrors.Internal, err, "소식을 조회하는 중에 에러가 발생했습니다.")
	}

	var newsDTOS []domain.NewsDTO
	for _, n := range news {
		newsDTOS = append(newsDTOS, domain.NewsDTOFrom(n))
	}

	var cursor *int
	if len(newsDTOS) > 0 {
		cursor = &newsDTOS[len(newsDTOS)-1].ID
	}

	return domain.ListNewsResponse{
		News:   newsDTOS,
		Cursor: cursor,
	}, nil
}

func (s newsService) UpdateNews(ctx context.Context, req domain.UpdateNewsRequest) error {
	const op cerrors.Op = "news/service/UpdateNews"

	news, err := s.newsRepository.FindNewsByID(ctx, req.ID)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}
	if news == nil {
		return cerrors.E(op, cerrors.NotExist, "소식을 찾을 수 없습니다.")
	}
	if news.UserID != req.UserID {
		return cerrors.E(op, cerrors.Permission, "소식을 수정할 권한이 없습니다.")
	}

	news.Title = req.Title

	if err := s.newsRepository.UpdateNews(ctx, *news); err != nil {
		return cerrors.E(op, cerrors.Internal, err, "소식을 수정하는 중에 에러가 발생했습니다.")
	}

	return nil
}

func (s newsService) DeleteNews(ctx context.Context, req domain.DeleteNewsRequest) error {
	const op cerrors.Op = "news/service/DeleteNews"

	news, err := s.newsRepository.FindNewsByID(ctx, req.ID)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}
	if news == nil {
		return cerrors.E(op, cerrors.NotExist, "소식을 찾을 수 없습니다.")
	}
	if news.UserID != req.UserID {
		return cerrors.E(op, cerrors.Permission, "소식을 삭제할 권한이 없습니다.")
	}

	if err := s.newsRepository.DeleteNews(ctx, req.ID); err != nil {
		return cerrors.E(op, cerrors.Internal, err, "소식을 삭제하는 중에 에러가 발생했습니다.")
	}

	return nil
}
