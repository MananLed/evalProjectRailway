package utils

import (
	"context"
	"errors"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID) 
	if !ok{
		return nil, errors.New("user ID not found in context")
	}

	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return nil, errors.New("user role not found in context")
	}

	return &model.User{
		ID:   id,
		Role: model.UserRole(role),
	}, nil
}
