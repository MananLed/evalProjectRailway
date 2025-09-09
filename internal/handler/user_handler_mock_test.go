package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MananLed/evalProjectRailway/internal/mocks"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*mocks.MockUserServiceInterface, *UserHandler, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockSvc := mocks.NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(mockSvc)
	return mockSvc, h, ctrl
}

func TestUserHandler_SignUp(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().IsEmailUnique("john@example.com").Return(true)
	mockSvc.EXPECT().SignUp(gomock.Any()).Return(nil)

	body := []byte(`{"firstname":"John","lastname":"Doe","email":"john@example.com","password":"Password@123","mobile":"9876543210"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}

func TestUserHandler_Login(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Email: "john@example.com", Role: model.RolePassenger}
	mockSvc.EXPECT().Login("john@example.com", "Password@123").Return(user, nil)

	body := []byte(`{"email":"john@example.com","password":"Password@123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}

func TestUserHandler_ViewProfile(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Email: "ctx@example.com"}
	ctx := context.WithValue(context.Background(), "user", *user)
	mockSvc.EXPECT().GetUserByID(ctx).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	h.ViewProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Email: "old@example.com"}
	ctx := context.WithValue(context.Background(), "user", *user)

	mockSvc.EXPECT().GetUserByID(ctx).Return(user, nil)
	mockSvc.EXPECT().IsEmailUnique("new@example.com").Return(true)
	mockSvc.EXPECT().UpdateProfile(gomock.Any()).Return(nil)

	update := map[string]string{"email": "new@example.com"}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestUserHandler_ChangePassword(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Email: "changepass@example.com"}
	ctx := context.WithValue(context.Background(), "user", *user)

	mockSvc.EXPECT().GetUserByID(ctx).Return(user, nil)
	mockSvc.EXPECT().ChangePassword(user, "oldPass", "newPass").Return(nil)

	body := []byte(`{"oldPassword":"oldPass","newPassword":"newPass"}`)
	req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewBuffer(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	h.ChangePassword(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestUserHandler_DeleteProfile(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Email: "delete@example.com"}
	ctx := context.WithValue(context.Background(), "user", *user)

	mockSvc.EXPECT().GetUserByID(ctx).Return(user, nil)
	mockSvc.EXPECT().DeleteProfile(user.ID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	h.DeleteProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestUserHandler_Login_Failure(t *testing.T) {
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().Login("wrong@example.com", "badPass").Return(nil, errors.New("invalid"))

	body := []byte(`{"email":"wrong@example.com","password":"badPass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", w.Code)
	}
}
