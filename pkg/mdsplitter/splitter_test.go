package mdsplitter

import (
	"reflect"
	"testing"
)

func TestSplitter_Split(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		input  string
		want   []Chunk
	}{
		{
			name: "basic header splitting with trim",
			config: &Config{
				Headers: map[string]string{
					"#":   "h1",
					"##":  "h2",
					"###": "h3",
				},
				TrimHeaders: true,
			},
			input: "# Header1\n\n```code1\ncode2\ncode3\n```\n## Header2\n\nContent1\n\n### Header3\n\nContent2\n\n## Header4\n\nContent3",
			want: []Chunk{
				{
					Content: "\n```code1\ncode2\ncode3\n```",
					Metadata: map[string]string{
						"h1": "Header1",
					},
				},
				{
					Content: "\nContent1\n",
					Metadata: map[string]string{
						"h1": "Header1",
						"h2": "Header2",
					},
				},
				{
					Content: "\nContent2\n",
					Metadata: map[string]string{
						"h1": "Header1",
						"h2": "Header2",
						"h3": "Header3",
					},
				},
				{
					Content: "\nContent3",
					Metadata: map[string]string{
						"h1": "Header1",
						"h2": "Header4",
					},
				},
			},
		},
		{
			name: "header splitting without trim",
			config: &Config{
				Headers: map[string]string{
					"##":  "h2",
					"###": "h3",
				},
				TrimHeaders: false,
			},
			input: "## Section 1\n\nSome content here.\n\n### Subsection 1.1\n\nMore content.\n\n## Section 2\n\nFinal content.",
			want: []Chunk{
				{
					Content: "## Section 1\n\nSome content here.\n",
					Metadata: map[string]string{
						"h2": "Section 1",
					},
				},
				{
					Content: "### Subsection 1.1\n\nMore content.\n",
					Metadata: map[string]string{
						"h2": "Section 1",
						"h3": "Subsection 1.1",
					},
				},
				{
					Content: "## Section 2\n\nFinal content.",
					Metadata: map[string]string{
						"h2": "Section 2",
					},
				},
			},
		},
		{
			name: "code blocks are preserved",
			config: &Config{
				Headers: map[string]string{
					"##": "h2",
				},
				TrimHeaders: false,
			},
			input: "## Code Example\n\n```go\n## This is not a header\nfunc main() {}\n```\n\n## Next Section\n\nContent",
			want: []Chunk{
				{
					Content: "## Code Example\n\n```go\n## This is not a header\nfunc main() {}\n```\n",
					Metadata: map[string]string{
						"h2": "Code Example",
					},
				},
				{
					Content: "## Next Section\n\nContent",
					Metadata: map[string]string{
						"h2": "Next Section",
					},
				},
			},
		},
		{
			name: "empty input",
			config: &Config{
				Headers: map[string]string{
					"##": "h2",
				},
				TrimHeaders: false,
			},
			input: "",
			want: []Chunk{
				{
					Content:  "",
					Metadata: map[string]string{},
				},
			},
		},
		{
			name: "no headers in text",
			config: &Config{
				Headers: map[string]string{
					"##": "h2",
				},
				TrimHeaders: false,
			},
			input: "Just some plain text\nwithout any headers\nat all.",
			want: []Chunk{
				{
					Content:  "Just some plain text\nwithout any headers\nat all.",
					Metadata: map[string]string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitter, err := New(tt.config)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got := splitter.Split(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() mismatch\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestNew_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "empty headers",
			config: &Config{
				Headers: map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "invalid header character",
			config: &Config{
				Headers: map[string]string{
					"##*": "h2",
				},
			},
			wantErr: true,
		},
		{
			name: "valid config",
			config: &Config{
				Headers: map[string]string{
					"##": "h2",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
