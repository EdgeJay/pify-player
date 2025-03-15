package utils

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "zero length",
			size: 0,
		},
		{
			name: "length of 10",
			size: 10,
		},
		{
			name: "length of 32",
			size: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRandomString(tt.size)
			if len(got) != tt.size {
				t.Errorf("GenerateRandomString() length = %v, want %v", len(got), tt.size)
			}

			// Test that all characters are valid
			const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			for _, c := range got {
				valid := false
				for _, l := range letters {
					if c == l {
						valid = true
						break
					}
				}
				if !valid {
					t.Errorf("GenerateRandomString() contains invalid character: %c", c)
				}
			}
		})
	}
}
