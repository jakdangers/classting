package subscription

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
	"k8s.io/utils/pointer"
)

type subscriptionService struct {
	newsRepository         domain.NewsRepository
	schoolRepository       domain.SchoolRepository
	subscriptionRepository domain.SubscriptionRepository
}

func NewSubscriptionService(
	newsRepository domain.NewsRepository,
	schoolRepository domain.SchoolRepository,
	subscriptionRepository domain.SubscriptionRepository,
) *subscriptionService {
	return &subscriptionService{
		newsRepository:         newsRepository,
		schoolRepository:       schoolRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

var _ domain.SubscriptionService = (*subscriptionService)(nil)

func (s subscriptionService) CreateSubscription(ctx context.Context, req domain.CreateSubscriptionRequest) error {
	const op cerrors.Op = "subscription/service/CreateSubscription"

	school, err := s.schoolRepository.FindSchoolByID(ctx, req.SchoolID)
	if err != nil {
		return err
	}
	if school == nil {
		return cerrors.E(op, cerrors.Invalid, "해당 학교가 존재하지 않습니다.")
	}

	subscription, err := s.subscriptionRepository.FindSubscriptionByUserIDAndSchoolID(ctx, domain.FindSubscriptionByUserIDAndSchoolIDParams{
		UserID:   req.UserID,
		SchoolID: req.SchoolID,
	})
	if err != nil {
		return err
	}
	if subscription != nil {
		return cerrors.E(op, cerrors.Invalid, "이미 구독한 학교입니다.")
	}

	_, err = s.subscriptionRepository.CreateSubscription(ctx, domain.Subscription{
		SchoolID: req.SchoolID,
		UserID:   req.UserID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s subscriptionService) ListSubscriptionSchools(ctx context.Context, req domain.ListSubscriptionSchoolsRequest) (domain.ListSubscriptionSchoolsResponse, error) {
	const op cerrors.Op = "subscription/service/ListSubscriptionSchools"

	subscriptionSchools, err := s.subscriptionRepository.ListSubscriptionSchools(ctx, domain.ListSubscriptionSchoolsParams{
		UserID: req.UserID,
		Cursor: req.Cursor,
	})
	if err != nil {
		return domain.ListSubscriptionSchoolsResponse{}, cerrors.E(op, cerrors.Internal, err, "구독한 학교를 조회하는 중에 에러가 발생했습니다.")
	}

	var subscriptionSchoolsDTOS []domain.SubscriptionSchoolDTO
	for _, n := range subscriptionSchools {
		subscriptionSchoolsDTOS = append(subscriptionSchoolsDTOS, domain.SubscriptionSchoolDTOFrom(n))
	}

	var cursor *int
	if len(subscriptionSchoolsDTOS) > 0 {
		cursor = &subscriptionSchoolsDTOS[len(subscriptionSchoolsDTOS)-1].ID
	}

	return domain.ListSubscriptionSchoolsResponse{
		SubscriptionSchools: subscriptionSchoolsDTOS,
		Cursor:              cursor,
	}, nil
}

func (s subscriptionService) ListSubscriptionSchoolNews(ctx context.Context, req domain.ListSubscriptionSchoolNewsRequest) (domain.ListSubscriptionSchoolNewsResponse, error) {
	const op cerrors.Op = "subscription/service/ListSubscriptionSchoolNews"

	subscription, err := s.subscriptionRepository.FindSubscriptionByUserIDAndSchoolID(ctx, domain.FindSubscriptionByUserIDAndSchoolIDParams{
		UserID:   req.UserID,
		SchoolID: req.SchoolID,
	})
	if err != nil {
		return domain.ListSubscriptionSchoolNewsResponse{}, err
	}
	if subscription == nil {
		return domain.ListSubscriptionSchoolNewsResponse{}, cerrors.E(op, cerrors.Invalid, "구독한 학교가 아닙니다.")
	}

	news, err := s.newsRepository.ListNews(ctx, domain.ListNewsParams{
		SchoolID: pointer.Int(req.SchoolID),
		Cursor:   req.Cursor,
	})
	if err != nil {
		return domain.ListSubscriptionSchoolNewsResponse{}, cerrors.E(op, cerrors.Internal, err, "소식을 조회하는 중에 에러가 발생했습니다.")
	}

	var newsDTOS []domain.SubscriptionSchoolNewsDTO
	for _, n := range news {
		newsDTOS = append(newsDTOS, domain.SubscriptionSchoolNewsDTOFrom(n))
	}

	var cursor *int
	if len(newsDTOS) > 0 {
		cursor = &newsDTOS[len(newsDTOS)-1].ID
	}

	return domain.ListSubscriptionSchoolNewsResponse{
		SubscriptionSchoolNews: newsDTOS,
		Cursor:                 cursor,
	}, nil
}

func (s subscriptionService) DeleteSubscription(ctx context.Context, req domain.DeleteSubscriptionRequest) error {
	const op cerrors.Op = "subscription/service/DeleteSubscription"

	subscription, err := s.subscriptionRepository.FindSubscriptionByUserIDAndSchoolID(ctx, domain.FindSubscriptionByUserIDAndSchoolIDParams{
		UserID:   req.UserID,
		SchoolID: req.SchoolID,
	})
	if err != nil {
		return err
	}
	if subscription == nil {
		return cerrors.E(op, cerrors.Invalid, "구독한 학교가 아닙니다.")
	}
	if subscription.UserID != req.UserID {
		return cerrors.E(op, cerrors.Permission, "구독한 학교를 삭제할 권한이 없습니다.")
	}

	if err := s.subscriptionRepository.DeleteSubscription(ctx, subscription.ID); err != nil {
		return cerrors.E(op, cerrors.Internal, err, "구독한 학교를 삭제하는 중에 에러가 발생했습니다.")
	}

	return nil
}
