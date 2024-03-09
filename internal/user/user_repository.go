package user

import (
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"context"
	"database/sql"
	"errors"
)

type userRepository struct {
	sqlDB *sql.DB
}

func NewUserRepository(sqlDB *sql.DB) *userRepository {
	return &userRepository{
		sqlDB: sqlDB,
	}
}

var _ domain.UserRepository = (*userRepository)(nil)

func (u userRepository) CreateUser(ctx context.Context, user domain.User) (int, error) {
	const op cerrors.Op = "user/userRepository/createUser"

	result, err := u.sqlDB.ExecContext(ctx, createUserQuery, user.UserName, user.Password, user.Type)
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return int(userID), nil
}

func (u userRepository) FindUserByUserName(ctx context.Context, mobileID string) (*domain.User, error) {
	const op cerrors.Op = "user/userRepository/findByUserID"
	var user domain.User

	err := u.sqlDB.QueryRowContext(ctx, findUserByUserNameQuery, mobileID).
		Scan(&user.ID, &user.UserName, &user.Password, &user.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return &user, nil
}
