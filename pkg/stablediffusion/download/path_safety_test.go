package download

import (
	"path/filepath"
	"testing"
)

func TestIsSafeExtractPath(t *testing.T) {
	dest := filepath.Join("data", "libraries", "stablediffusion", "temp")

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "dest_root_is_allowed",
			path: dest,
			want: true,
		},
		{
			name: "child_path_is_allowed",
			path: filepath.Join(dest, "libgosd-fallback.so"),
			want: true,
		},
		{
			name: "path_traversal_is_blocked",
			path: filepath.Join(dest, "..", "outside"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSafeExtractPath(dest, tt.path)
			if got != tt.want {
				t.Fatalf("isSafeExtractPath(%q, %q) = %v, want %v", dest, tt.path, got, tt.want)
			}
		})
	}
}
