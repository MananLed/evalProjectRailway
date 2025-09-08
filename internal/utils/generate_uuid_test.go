package utils

import "testing"

func TestGenerateUUID_NotEmpty(t *testing.T) {
	u := GenerateUUID()
	if u == [16]byte{} {
		t.Errorf("expected non-empty UUID, got zero value")
	}
}

func TestGenerateUUID_Unique(t *testing.T) {
	u1 := GenerateUUID()
	u2 := GenerateUUID()

	if u1 == u2 {
		t.Errorf("expected different UUIDs, got same: %v", u1)
	}
}