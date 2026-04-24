package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/service"
)

// UserRepository の Mock
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepo) GetRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	password := "correct-password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	testUser := &model.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Roles:        []string{"employee"},
	}

	t.Run("成功: 正しいメールアドレスとパスワード", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		// mockRepo.FindByEmail が呼ばれたら testUser を返すように設定
		mockRepo.On("FindByEmail", ctx, "test@example.com").Return(testUser, nil)

		authService := service.NewAuthService(mockRepo)
		resp, err := authService.Login(ctx, "test@example.com", password)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, testUser.ID, resp.User.ID)
		assert.NotEmpty(t, resp.Token)
		
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: ユーザーが見つからない", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		mockRepo.On("FindByEmail", ctx, "unknown@example.com").Return(nil, nil)

		authService := service.NewAuthService(mockRepo)
		resp, err := authService.Login(ctx, "unknown@example.com", password)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "正しくありません")
		
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: パスワードが間違っている", func(t *testing.T) {
		mockRepo := new(mockUserRepo)
		mockRepo.On("FindByEmail", ctx, "test@example.com").Return(testUser, nil)

		authService := service.NewAuthService(mockRepo)
		resp, err := authService.Login(ctx, "test@example.com", "wrong-password")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "正しくありません")
		
		mockRepo.AssertExpectations(t)
	})
}
