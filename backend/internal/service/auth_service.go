package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/amical/routine-design/backend/internal/middleware"
	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
)

// AuthService は認証ビジネスロジックを提供する。
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService はAuthServiceを生成する。
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login はメールアドレスとパスワードで認証し、JWTを発行する。
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		slog.Warn("ログイン失敗: ユーザーが見つかりません", "email", email)
		return nil, fmt.Errorf("メールアドレスまたはパスワードが正しくありません: %w", ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.Warn("ログイン失敗: パスワード不一致", "email", email)
		return nil, fmt.Errorf("メールアドレスまたはパスワードが正しくありません: %w", ErrUnauthorized)
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func generateToken(user *model.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Roles:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-key"
	}
	return token.SignedString([]byte(secret))
}
