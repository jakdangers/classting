package user

import (
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"context"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var (
	mobileIDPattern     = regexp.MustCompile(`^(010-\d{4}-\d{4}|010\d{8})$`)
	hasHyphenPattern    = regexp.MustCompile(`-`)
	removeHyphenPattern = regexp.MustCompile(`-`)
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(
	userRepository domain.UserRepository,
) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

var _ domain.UserService = (*userService)(nil)

func (us userService) CreateUser(ctx context.Context, req domain.CreateUserRequest) error {
	const op cerrors.Op = "user/service/createUser"

	userName, err := validateUserName(req.UserName)
	if err != nil {
		return err
	}

	user, err := us.userRepository.FindUserByUserName(ctx, userName)
	if err != nil {
		return err
	}
	if user != nil {
		return cerrors.E(op, cerrors.Invalid, "이미 사용중인 유저아이디입니다.")
	}

	hashedPassword, err := hashPasswordWithSalt(req.Password)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	_, err = us.userRepository.CreateUser(ctx, domain.User{
		UserName: userName,
		Password: hashedPassword,
		UseType:  req.UserType,
	})
	if err != nil {
		return err
	}

	return nil
}

func hashPasswordWithSalt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func compareHashAndPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validateUserName(userName string) (string, error) {
	const op cerrors.Op = "user/service/validateUserName"

	if userName == "" {
		return "", cerrors.E(op, cerrors.Invalid, "잘못된 아이디입니다.")
	}

	return userName, nil
}
