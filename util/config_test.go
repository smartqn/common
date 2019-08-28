package util

import "testing"

func TestCurDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// {"normal", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurDir(); got != tt.want {
				t.Errorf("CurDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
