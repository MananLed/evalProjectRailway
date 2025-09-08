package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	users map[uuid.UUID]model.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[uuid.UUID]model.User),
	}
}

func (m *MockUserRepo) IsEmailUnique(email string) bool {
	for _, u := range m.users {
		if u.Email == email {
			return false
		}
	}
	return true
}

func (m *MockUserRepo) AddUser(user model.User) error {
	if !m.IsEmailUnique(user.Email) {
		return errors.New("email already exists")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) GetUserByEmailAndPassword(email, password string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email && u.Password == password {
			return &u, nil
		}
	}
	return nil, errors.New("invalid credentials")
}

func (m *MockUserRepo) GetUserByID(id uuid.UUID) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return &u, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepo) UpdateUser(user model.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepo) ChangePassword(id uuid.UUID, hashedPassword string) error {
	if u, ok := m.users[id]; ok {
		u.Password = hashedPassword
		m.users[id] = u
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepo) DeleteUserByID(id uuid.UUID) error {
	if _, ok := m.users[id]; ok {
		delete(m.users, id)
		return nil
	}
	return errors.New("user not found")
}

func TestUserService_SignUpAndLogin(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	user := model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "secret",
	}

	err := service.SignUp(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.Login("test@example.com", "secret")
	if err != nil {
		t.Errorf("expected successful login, got error %v", err)
	}

	_, err = service.Login("test@example.com", "wrongpass")
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	user := model.User{
		ID:       uuid.New(),
		Email:    "change@example.com",
		Password: string(hashed),
	}
	repo.users[user.ID] = user

	err := service.ChangePassword(&user, "oldpass", "newpass")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = service.ChangePassword(&user, "wrongpass", "anotherpass")
	if err == nil {
		t.Errorf("expected error for wrong current password, got nil")
	}
}

func TestUserService_DeleteProfile(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	user := model.User{ID: uuid.New(), Email: "delete@example.com", Password: "pass"}
	repo.users[user.ID] = user

	err := service.DeleteProfile(user.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = service.DeleteProfile(user.ID)
	if err == nil {
		t.Errorf("expected error for deleting non-existing user, got nil")
	}
}

func TestUserService_IsEmailUnique(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	user := model.User{ID: uuid.New(), Email: "unique@example.com", Password: "pass"}
	repo.users[user.ID] = user

	if service.IsEmailUnique("unique@example.com") {
		t.Errorf("expected email to be not unique, got unique")
	}

	if !service.IsEmailUnique("other@example.com") {
		t.Errorf("expected email to be unique, got not unique")
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	user := model.User{ID: uuid.New(), Email: "context@example.com", Password: "pass"}
	repo.users[user.ID] = user

	ctx := context.WithValue(context.Background(), utils.UserIDKey, user.ID)
	ctx = context.WithValue(ctx, utils.UserRoleKey, string(model.RolePassenger))
	got, err := service.GetUserByID(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if got == nil || got.ID != user.ID {
		t.Errorf("expected user %v, got %v", user.ID, got)
	}

	ctx = context.Background()
	_, err = service.GetUserByID(ctx)
	if err == nil {
		t.Errorf("expected error when no user in context, got nil")
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	user := model.User{ID: uuid.New(), Email: "update@example.com", Password: "pass"}
	repo.users[user.ID] = user

	updated := model.User{ID: user.ID, Email: "updated@example.com", Password: "newpass"}
	err := service.UpdateProfile(updated)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	got, _ := repo.GetUserByID(user.ID)
	if got.Email != "updated@example.com" {
		t.Errorf("expected email updated@example.com, got %v", got.Email)
	}

	nonExistent := model.User{ID: uuid.New(), Email: "ghost@example.com"}
	err = service.UpdateProfile(nonExistent)
	if err == nil {
		t.Errorf("expected error for updating non-existing user, got nil")
	}
}
