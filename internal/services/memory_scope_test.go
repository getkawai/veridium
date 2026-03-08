package services

import "testing"

func TestBuildMemoryScopeKey(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		sessionID string
		want      string
	}{
		{
			name:      "user and session",
			userID:    "u1",
			sessionID: "s1",
			want:      "u1:s1",
		},
		{
			name:      "missing user",
			userID:    "",
			sessionID: "s1",
			want:      "anon:s1",
		},
		{
			name:      "missing session",
			userID:    "u1",
			sessionID: "",
			want:      "u1:default",
		},
		{
			name:      "both missing",
			userID:    "",
			sessionID: "",
			want:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildMemoryScopeKey(tc.userID, tc.sessionID)
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
