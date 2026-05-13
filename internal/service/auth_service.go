package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type AuthService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(fullName, email, password string) error {
	existing, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}

	if existing != nil {
		return domain.ErrBadRequest("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return domain.ErrInternal("could not hash password")
	}

	user := domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "customer",
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", domain.ErrInternal("something went wrong")
	}
	if user == nil {
		return "", domain.ErrBadRequest("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrBadRequest("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", domain.ErrInternal("could not create token")
	}

	return signed, nil
}
