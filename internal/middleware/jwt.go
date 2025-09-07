package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MananLed/evalProjectRailway/internal/response"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			response.ErrorResponse(w, http.StatusUnauthorized, "Missing token", 1002)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token format", 1002)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		token, err := utils.VerifyJwt(tokenString)
		if err != nil {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		userID, ok := claims["user_id"].(uuid.UUID)
		if !ok {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}
		
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
