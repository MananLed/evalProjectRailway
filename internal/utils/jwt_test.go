package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT_Success(t *testing.T) {
	token, err := GenerateJWT(GenerateUUID(), "passenger")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token string")
	}
}

func TestVerifyJwt_ValidToken(t *testing.T) {
	tokenStr, err := GenerateJWT(GenerateUUID(), "passenger")
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	token, err := VerifyJwt(tokenStr)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}
}

func TestVerifyJwt_InvalidSignature(t *testing.T) {
	tokenStr, _ := GenerateJWT(GenerateUUID(), "passenger")

	badToken := tokenStr[:len(tokenStr)-1] + "x"

	_, err := VerifyJwt(badToken)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}
}

func TestVerifyJwt_InvalidSigningMethod(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": GenerateUUID(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := VerifyJwt(tokenStr)
	if err == nil {
		t.Fatal("expected error for invalid signing method, got nil")
	}
}

func TestVerifyJwt_ExpiredToken(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": GenerateUUID(),
		"exp":     time.Now().Add(-time.Hour).Unix(), // expired
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	_, err := VerifyJwt(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestVerifyJwt_MalformedToken(t *testing.T) {
	_, err := VerifyJwt("not-a-token")
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}

	_, err = VerifyJwt("")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
