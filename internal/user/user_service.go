package user

import (
	"classting/config"
	"classting/domain"
	cerrors "classting/pkg/cerrors"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userService struct {
	userRepository domain.UserRepository
	cfg            *config.Config
}

func NewUserService(
	userRepository domain.UserRepository,
	cfg *config.Config,
) *userService {
	return &userService{
		userRepository: userRepository,
		cfg:            cfg,
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

func (us userService) LoginUser(ctx context.Context, req domain.LoginUserRequest) (domain.LoginUserResponse, error) {
	const op cerrors.Op = "user/service/loginUser"

	userName, err := validateUserName(req.UserName)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	user, err := us.userRepository.FindUserByUserName(ctx, userName)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}
	if user == nil {
		return domain.LoginUserResponse{}, cerrors.E(op, cerrors.Invalid, "아이디 또는 비밀번호를 확인해주세요.")
	}

	if !compareHashAndPassword(req.Password, user.Password) {
		return domain.LoginUserResponse{}, cerrors.E(op, cerrors.Invalid, "아이디 또는 비밀번호를 확인해주세요.")
	}

	creationTime := time.Now().UTC()
	expirationTime := creationTime.Add(time.Hour * time.Duration(us.cfg.ExpiryHours))

	accessToken, err := createAccessToken(*user, us.cfg.Auth.Secret, expirationTime)
	if err != nil {
		return domain.LoginUserResponse{}, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return domain.LoginUserResponse{
		AccessToken: accessToken,
		ExpiresIn:   expirationTime.Unix(),
	}, nil
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

func createAccessToken(user domain.User, secret string, exp time.Time) (accessToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"userType": user.UseType,
		"exp":      exp.Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
