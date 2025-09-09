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
	mockSvc.EXPECT().IsEmailUnique("john1@example.com").Return(false)
	mockSvc.EXPECT().SignUp(gomock.Any()).Return(nil)

	body := []byte(`{"firstname":"John","lastname":"Doe","email":"john@example.com","password":"Password@123","mobile":"9876543210"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}

	body2 := []byte(`{bjljnjkh}`)
	req2 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body2))
	w2 := httptest.NewRecorder()

	h.SignUp(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w2.Code)
	}

	body3 := []byte(`{"firstname":"John","lastname":"Doe","email":"johnexample.com","password":"Password@123","mobile":"9876543210"}`)
	req3 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body3))
	w3 := httptest.NewRecorder()

	h.SignUp(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w.Code)
	}

	body4 := []byte(`{"firstname":"John","lastname":"Doe","email":"john1@example.com","password":"Password@123","mobile":"9876543210"}`)
	req4 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body4))
	w4 := httptest.NewRecorder()

	h.SignUp(w4, req4)

	if w4.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w4.Code)
	}
}

func TestUserHandler_SignUpFailed(t *testing.T){
	mockSvc, h, ctrl := setup(t)
	defer ctrl.Finish()

	mockSvc.EXPECT().IsEmailUnique("john@example.com").Return(true)
	mockSvc.EXPECT().SignUp(gomock.Any()).Return(errors.New("sjflds"))

	body := []byte(`{"firstname":"John","lastname":"Doe","email":"john@example.com","password":"Password@123","mobile":"9876543210"}`)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected StatusInternalServerError, got %d", w.Code)
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

	body2 := []byte(`gygjkkk`)
	req2 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body2))
	w2 := httptest.NewRecorder()

	h.Login(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w2.Code)
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

	user1 := &model.User{ID: uuid.New(), Email: "ctx1@example.com"}
	ctx1 := context.WithValue(context.Background(), "user", *user1)
	mockSvc.EXPECT().GetUserByID(ctx1).Return(user1, errors.New("dlfjld"))

	req1 := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx1)
	w1 := httptest.NewRecorder()

	h.ViewProfile(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected StatusUnauthorized, got %d", w1.Code)
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

	update := map[string]string{"firstname":"nin", "middlename":"ja", "lastname":"aa", "email": "new@example.com", "mobile":"9878889989"}
	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	user1 := &model.User{ID: uuid.New(), Email: "old@example.com"}
	ctx1 := context.WithValue(context.Background(), "user", *user1)

	mockSvc.EXPECT().GetUserByID(ctx1).Return(user, errors.New("dlsfjsl"))

	update1 := map[string]string{"firstname":"nin", "middlename":"ja", "lastname":"aa", "email": "new@example.com", "mobile":"9878889989"}
	body1, _ := json.Marshal(update1)
	req1 := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body1)).WithContext(ctx1)
	w1 := httptest.NewRecorder()

	h.UpdateProfile(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected StatusUnauthorized, got %d", w1.Code)
	}

	user2 := &model.User{ID: uuid.New(), Email: "old@example.com"}
	ctx2 := context.WithValue(context.Background(), "user", *user2)

	mockSvc.EXPECT().GetUserByID(ctx2).Return(user2, nil)
	mockSvc.EXPECT().IsEmailUnique("new@example.com").Return(true)
	mockSvc.EXPECT().UpdateProfile(gomock.Any()).Return(errors.New("flsjf"))

	update2 := map[string]string{"firstname":"nin", "middlename":"ja", "lastname":"aa", "email": "new@example.com", "mobile":"9878889989"}
	body2, _ := json.Marshal(update2)
	req2 := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body2)).WithContext(ctx2)
	w2 := httptest.NewRecorder()

	h.UpdateProfile(w2, req2)

	if w2.Code != http.StatusInternalServerError {
		t.Errorf("expected StatusInternalServerError, got %d", w2.Code)
	}

	user3 := &model.User{ID: uuid.New(), Email: "old@example.com"}
	ctx3 := context.WithValue(context.Background(), "user", *user3)

	mockSvc.EXPECT().GetUserByID(ctx3).Return(user, nil)
	mockSvc.EXPECT().IsEmailUnique("new11@example.com").Return(false)

	update3 := map[string]string{"firstname":"nin", "middlename":"ja", "lastname":"aa", "email": "new11@example.com", "mobile":"9878889989"}
	body3, _ := json.Marshal(update3)
	req3 := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(body3)).WithContext(ctx3)
	w3 := httptest.NewRecorder()

	h.UpdateProfile(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w3.Code)
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

	user1 := &model.User{ID: uuid.New(), Email: "changepass@example.com"}
	ctx1 := context.WithValue(context.Background(), "user", *user1)

	mockSvc.EXPECT().GetUserByID(ctx1).Return(user, errors.New("sldjfl"))
	// mockSvc.EXPECT().ChangePassword(user, "oldPass", "newPass").Return(nil)

	body1 := []byte(`{"oldPassword":"oldPass","newPassword":"newPass"}`)
	req1 := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewBuffer(body1)).WithContext(ctx1)
	w1 := httptest.NewRecorder()

	h.ChangePassword(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected StatusUnauthorized, got %d", w1.Code)
	}

	user2 := &model.User{ID: uuid.New(), Email: "changepass@example.com"}
	ctx2 := context.WithValue(context.Background(), "user", *user2)

	mockSvc.EXPECT().GetUserByID(ctx2).Return(user2, nil)
	mockSvc.EXPECT().ChangePassword(user2, "oldPass", "newPass").Return(errors.New("slfjls"))

	body2 := []byte(`{"oldPassword":"oldPass","newPassword":"newPass"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewBuffer(body2)).WithContext(ctx2)
	w2 := httptest.NewRecorder()

	h.ChangePassword(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected StatusUnauthorized, got %d", w2.Code)
	}

	user3 := &model.User{ID: uuid.New(), Email: "changepass@example.com"}
	ctx3 := context.WithValue(context.Background(), "user", *user3)

	mockSvc.EXPECT().GetUserByID(ctx3).Return(user3, nil)
	// mockSvc.EXPECT().ChangePassword(user2, "oldPass", "newPass").Return(errors.New("slfjls"))

	body3 := []byte(`dfsflsfls`)
	req3 := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewBuffer(body3)).WithContext(ctx3)
	w3 := httptest.NewRecorder()

	h.ChangePassword(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %d", w3.Code)
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

	user1 := &model.User{ID: uuid.New(), Email: "delete1@example.com"}
	ctx1 := context.WithValue(context.Background(), "user", *user1)

	mockSvc.EXPECT().GetUserByID(ctx1).Return(user1, errors.New("ldfjls"))

	req1 := httptest.NewRequest(http.MethodDelete, "/delete", nil).WithContext(ctx1)
	w1 := httptest.NewRecorder()

	h.DeleteProfile(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected StatusUnauthorized, got %d", w1.Code)
	}

	user2 := &model.User{ID: uuid.New(), Email: "delete2@example.com"}
	ctx2 := context.WithValue(context.Background(), "user", *user2)

	mockSvc.EXPECT().GetUserByID(ctx2).Return(user2, nil)
	mockSvc.EXPECT().DeleteProfile(user2.ID).Return(errors.New("fsldjfsdl"))

	req2 := httptest.NewRequest(http.MethodDelete, "/delete", nil).WithContext(ctx2)
	w2 := httptest.NewRecorder()

	h.DeleteProfile(w2, req2)

	if w2.Code != http.StatusInternalServerError {
		t.Errorf("expected StatusInternalServerError, got %d", w2.Code)
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
