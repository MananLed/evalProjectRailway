package utils

import "testing"

func TestValidateMobileNumber(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"9876543210", true},
		{"1234567890", false},
		{"98765", false},
		{" 9876543210 ", true},
	}

	for _, tt := range tests {
		got := ValidateMobileNumber(tt.input)
		if got != tt.want {
			t.Errorf("ValidateMobileNumber(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name+tag@domain.co", true},
		{"invalid-email", false},
		{"missing@domain", false},
	}

	for _, tt := range tests {
		got := ValidateEmail(tt.input)
		if got != tt.want {
			t.Errorf("ValidateEmail(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"short1!", false},
		{"alllowercasepassword123!", true},
		{"NoSpecialCharacter123", false},
		{"nouppercasebutvalid123!", true},
		{"missingdigit!!!!", false},
	}

	for _, tt := range tests {
		got := ValidatePassword(tt.input)
		if got != tt.want {
			t.Errorf("ValidatePassword(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
