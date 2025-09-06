package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/response"
	"github.com/MananLed/evalProjectRailway/internal/service"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService service.UserServiceInterface
}

func NewUserHandler(us service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Firstname  string `json:"firstname"`
		Middlename string `json:"middlename"`
		Lastname   string `json:"lastname"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Mobile     string `json:"mobile"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	if !utils.ValidateEmail(req.Email) || !utils.ValidateMobileNumber(req.Mobile) || !utils.ValidatePassword(req.Password) {
		logger.LogToFile("Validation error")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	if !h.UserService.IsEmailUnique(req.Email) {
		logger.LogToFile("Email not unique")
		response.ErrorResponse(w, http.StatusBadRequest, "Email cannot be used", 1001)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password", 1010)
		return
	}

	user := model.User{
		ID:         utils.GenerateUUID(),
		Firstname:  req.Firstname,
		Middlename: req.Middlename,
		Lastname:   req.Lastname,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Mobile:     req.Mobile,
		Role:       model.RolePassenger,
	}

	err = h.UserService.SignUp(user)

	if err != nil {
		logger.LogToFile("Error creating user")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error creating user", 1006)
		return
	}

	logger.LogToFile("User created successfully")
	response.SuccessResponse(w, nil, "User created successfully", http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid Input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid Input", 1001)
		return
	}

	user, err := h.UserService.Login(req.Email, req.Password)

	if err != nil {
		logger.LogToFile("Invalid email or password")
		response.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password", 1005)
		return
	}

	var jwtTokenString string

	jwtTokenString, err = utils.GenerateJWT(user.ID, string(user.Role))

	if err != nil {
		logger.LogToFile("Error generating the token")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error genrating the token", 1006)
		return
	}

	response.SuccessResponse(w, map[string]interface{}{"token": jwtTokenString}, "Token generated successfully", http.StatusCreated)
}

func (h *UserHandler) ViewProfile(w http.ResponseWriter, r *http.Request) {

	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user id not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	response.SuccessResponse(w, user, "Users Retrived Successfully", http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	user, err := h.UserService.GetUserByID(r.Context())
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	var req struct {
		Firstname  string `json:"firstName"`
		Middlename string `json:"middleName"`
		Lastname   string `json:"lastName"`
		Email      string `json:"email"`
		Mobile     string `json:"mobile"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	if req.Firstname != "" {
		user.Firstname = req.Firstname
	}
	if req.Middlename != "" {
		user.Middlename = req.Middlename
	}
	if req.Lastname != "" {
		user.Lastname = req.Lastname
	}
	if req.Email != "" {
		if !h.UserService.IsEmailUnique(req.Email) {
		logger.LogToFile("Email not unique")
		response.ErrorResponse(w, http.StatusBadRequest, "Email cannot be used", 1001)
		return
		}
		user.Email = req.Email
	}
	if req.Mobile != "" {
		user.Mobile = req.Mobile
	}
	err = h.UserService.UpdateProfile(*user)

	if err != nil {
		logger.LogToFile("Failed to update user: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user", 1011)
		return
	}

	logger.LogToFile("User updated successfully")
	response.SuccessResponse(w, nil, "User updated successfully", http.StatusOK)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	err = h.UserService.ChangePassword(user, req.OldPassword, req.NewPassword)
	if err != nil {
		logger.LogToFile("unauthorized user changing password")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}
	response.SuccessResponse(w, nil, "User password updated successfully", http.StatusOK)
}

func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	err = h.UserService.DeleteProfile(user.ID)
	if err != nil {
		logger.LogToFile("Error deleting user")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting user", 1010)
		return
	}

	response.SuccessResponse(w, nil, "Profile deleted successfully", http.StatusOK)
}
