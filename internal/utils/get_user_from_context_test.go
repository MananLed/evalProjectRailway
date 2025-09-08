package utils

import (
	"context"
	"testing"

	"github.com/MananLed/evalProjectRailway/internal/model"
)

func TestGetUserFromContext_Success(t *testing.T) {

	id := GenerateUUID()

	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, id)
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RolePassenger))

	user, err := GetUserFromContext(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != id {
		t.Errorf("unexpected user values: %+v", user)
	}
	if user.Role != model.RolePassenger {
		t.Errorf("expected role Resident, got %s", user.Role)
	}
}

func TestGetUserFromContext_MissingID(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserRoleKey, string(model.RolePassenger))

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
}

func TestGetUserFromContext_MissingRole(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, GenerateUUID())

	_, err := GetUserFromContext(ctx)
	if err == nil {
		t.Fatal("expected error for missing role, got nil")
	}
}
