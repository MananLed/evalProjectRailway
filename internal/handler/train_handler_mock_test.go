package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MananLed/evalProjectRailway/internal/mocks"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func setupTrainHandler(t *testing.T) (*mocks.MockTrainServiceInterface, *mocks.MockUserServiceInterface, *TrainHandler, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	trainSvc := mocks.NewMockTrainServiceInterface(ctrl)
	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	h := NewTrainHandler(trainSvc, userSvc)
	return trainSvc, userSvc, h, ctrl
}

func TestAddTrain_Success(t *testing.T) {
	trainSvc, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	admin := &model.User{ID: uuid.New(), Role: model.RoleAdmin}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil)
	trainSvc.EXPECT().AddTrain(gomock.Any()).Return(nil)

	body := []byte(`{
		"name":"Express",
		"source":"Delhi",
		"destination":"Mumbai",
		"departuredate":"2025-09-10",
		"departuretime":"10:00",
		"arrivaldate":"2025-09-11",
		"arrivaltime":"08:00",
		"totalseats":100,
		"seatfare":500
	}`)
	req := httptest.NewRequest(http.MethodPost, "/trains", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.AddTrain(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}

func TestAddTrain_Unauthorized(t *testing.T) {
	_, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("unauthorized"))

	req := httptest.NewRequest(http.MethodPost, "/trains", bytes.NewBuffer([]byte(`{}`)))
	w := httptest.NewRecorder()

	h.AddTrain(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden, got %d", w.Code)
	}
}

func TestAddTrain_InvalidInput(t *testing.T) {
	trainSvc, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	admin := &model.User{ID: uuid.New(), Role: model.RoleAdmin}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil).AnyTimes()

	req := httptest.NewRequest(http.MethodPost, "/trains", bytes.NewBuffer([]byte(`invalid`)))
	w := httptest.NewRecorder()
	h.AddTrain(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}

	body := []byte(`{
		"name":"Test",
		"source":"A",
		"destination":"B",
		"departuredate":"2025-09-12",
		"departuretime":"10:00",
		"arrivaldate":"2025-09-12",
		"arrivaltime":"08:00",
		"totalseats":50,
		"seatfare":100
	}`)
	req2 := httptest.NewRequest(http.MethodPost, "/trains", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	h.AddTrain(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w2.Code)
	}

	body2 := []byte(`{
		"name":"Express",
		"source":"Delhi",
		"destination":"Mumbai",
		"departuredate":"2025-09-10",
		"departuretime":"10:00",
		"arrivaldate":"2025-09-11",
		"arrivaltime":"08:00",
		"totalseats":100,
		"seatfare":500
	}`)
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil).AnyTimes()
	trainSvc.EXPECT().AddTrain(gomock.Any()).Return(errors.New("db error"))

	req3 := httptest.NewRequest(http.MethodPost, "/trains", bytes.NewBuffer(body2))
	w3 := httptest.NewRecorder()
	h.AddTrain(w3, req3)

	if w3.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 InternalServerError, got %d", w3.Code)
	}
}

func TestDeleteTrain_Success(t *testing.T) {
	trainSvc, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	admin := &model.User{ID: uuid.New(), Role: model.RoleAdmin}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil)

	id := uuid.New()
	trainSvc.EXPECT().DeleteTrain(id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/trains/"+id.String(), nil)
	w := httptest.NewRecorder()

	h.DeleteTrain(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestDeleteTrain_InvalidID(t *testing.T) {
	_, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	admin := &model.User{ID: uuid.New(), Role: model.RoleAdmin}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil)

	req := httptest.NewRequest(http.MethodDelete, "/trains/not-a-uuid", nil)
	w := httptest.NewRecorder()

	h.DeleteTrain(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 BadRequest, got %d", w.Code)
	}
}

func TestViewTrain_Success(t *testing.T) {
	trainSvc, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New(), Role: model.RolePassenger}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(user, nil)

	trains := []model.Train{{ID: uuid.New(), Name: "Express"}}
	trainSvc.EXPECT().GetTrains(gomock.Any(), "Delhi", "Mumbai").Return(trains, nil)

	req := httptest.NewRequest(http.MethodGet, "/trains?date=2025-09-10&source=Delhi&destination=Mumbai", nil)
	w := httptest.NewRecorder()

	h.ViewTrain(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestViewTrain_InvalidDate(t *testing.T) {
	_, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New()}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/trains?date=09-10-2025", nil)
	w := httptest.NewRecorder()

	h.ViewTrain(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 BadRequest, got %d", w.Code)
	}
}

func TestViewTrain_ServiceError(t *testing.T) {
	trainSvc, userSvc, h, ctrl := setupTrainHandler(t)
	defer ctrl.Finish()

	user := &model.User{ID: uuid.New()}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(user, nil)

	trainSvc.EXPECT().GetTrains(gomock.Any(), "", "").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/trains", nil)
	w := httptest.NewRecorder()

	h.ViewTrain(w, req)

	if w.Code != http.StatusUnauthorized { 
		t.Errorf("expected 401 Unauthorized, got %d", w.Code)
	}
}
