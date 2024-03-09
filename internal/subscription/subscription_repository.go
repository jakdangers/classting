package subscription

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type subscriptionRepository struct {
	sqlDB *sql.DB
}

func NewSubscriptionRepository(sqlDB *sql.DB) *subscriptionRepository {
	return &subscriptionRepository{
		sqlDB: sqlDB,
	}
}

var _ domain.SubscriptionRepository = (*subscriptionRepository)(nil)

func (n subscriptionRepository) CreateSubscription(ctx context.Context, subscription domain.Subscription) (int, error) {
	const op cerrors.Op = "subscription/subscriptionRepository/CreateSubscription"

	result, err := n.sqlDB.ExecContext(ctx, createSubscriptionQuery, subscription.SchoolID, subscription.UserID)
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	subscriptionID, err := result.LastInsertId()
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return int(subscriptionID), nil
}

func (n subscriptionRepository) ListSubscriptionSchools(ctx context.Context, params domain.ListSubscriptionSchoolsParams) ([]domain.SubscriptionSchool, error) {
	const op cerrors.Op = "subscription/subscriptionRepository/ListSubscriptionSchools"

	var subscriptionSchools []domain.SubscriptionSchool

	query := fmt.Sprintf(listSubscriptionSchoolsQuery, params.AfterCursor())

	rows, err := n.sqlDB.QueryContext(ctx, query, params.UserID)
	if err != nil {
		fmt.Println(err)
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}
	defer rows.Close()

	for rows.Next() {
		var subscriptionSchool domain.SubscriptionSchool
		err := rows.Scan(
			&subscriptionSchool.ID,
			&subscriptionSchool.CreateDate,
			&subscriptionSchool.UpdateDate,
			&subscriptionSchool.SchoolID,
			&subscriptionSchool.Name,
			&subscriptionSchool.Region,
		)
		if err != nil {
			return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
		}
		subscriptionSchools = append(subscriptionSchools, subscriptionSchool)
	}

	return subscriptionSchools, nil
}

func (n subscriptionRepository) DeleteSubscription(ctx context.Context, subscriptionID int) error {
	const op cerrors.Op = "subscription/subscriptionRepository/DeleteSubscription"

	_, err := n.sqlDB.ExecContext(ctx, deleteSubscriptionQuery, subscriptionID)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return nil
}

func (n subscriptionRepository) FindSubscriptionByUserIDAndSchoolID(ctx context.Context, params domain.FindSubscriptionByUserIDAndSchoolIDParams) (*domain.Subscription, error) {
	const op cerrors.Op = "subscription/subscriptionRepository/FindSubscriptionByUserIDAndSchoolID"

	var subscription domain.Subscription

	err := n.sqlDB.QueryRowContext(ctx, findSubscriptionByUserIDAndSchoolIDQuery, params.UserID, params.SchoolID).Scan(
		&subscription.ID,
		&subscription.CreateDate,
		&subscription.UpdateDate,
		&subscription.SchoolID,
		&subscription.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return &subscription, nil
}
