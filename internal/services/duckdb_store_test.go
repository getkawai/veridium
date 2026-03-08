package services

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCreateVectorsTableQuery_NoDropTable(t *testing.T) {
	query := buildCreateVectorsTableQuery(768)
	normalized := strings.ToUpper(query)

	assert.NotContains(t, normalized, "DROP TABLE")
	assert.Contains(t, normalized, "CREATE TABLE IF NOT EXISTS VECTORS")
	assert.Contains(t, normalized, "EMBEDDING FLOAT[768]")
}

func TestBuildCreateHNSWIndexQuery_UsesConfigValues(t *testing.T) {
	query := buildCreateHNSWIndexQuery(&HNSWConfig{
		Metric:         "cosine",
		EfConstruction: 256,
		EfSearch:       128,
		M:              24,
	})
	normalized := strings.ToUpper(query)

	if !strings.Contains(normalized, "CREATE INDEX IF NOT EXISTS VEC_IDX") {
		t.Fatalf("query must create vec_idx idempotently, got: %s", query)
	}
	if !strings.Contains(normalized, "METRIC = 'COSINE'") {
		t.Fatalf("query must include metric from config, got: %s", query)
	}
	if !strings.Contains(normalized, "EF_CONSTRUCTION = 256") {
		t.Fatalf("query must include ef_construction from config, got: %s", query)
	}
	if !strings.Contains(normalized, "EF_SEARCH = 128") {
		t.Fatalf("query must include ef_search from config, got: %s", query)
	}
	if !strings.Contains(normalized, "M = 24") {
		t.Fatalf("query must include M from config, got: %s", query)
	}
}

func TestBuildCreateHNSWIndexQuery_NilConfigUsesDefault(t *testing.T) {
	query := buildCreateHNSWIndexQuery(nil)
	normalized := strings.ToUpper(query)

	if !strings.Contains(normalized, "METRIC = 'L2SQ'") {
		t.Fatalf("query must include default metric when config is nil, got: %s", query)
	}
	if !strings.Contains(normalized, "EF_CONSTRUCTION = 128") {
		t.Fatalf("query must include default ef_construction when config is nil, got: %s", query)
	}
	if !strings.Contains(normalized, "EF_SEARCH = 64") {
		t.Fatalf("query must include default ef_search when config is nil, got: %s", query)
	}
	if !strings.Contains(normalized, "M = 16") {
		t.Fatalf("query must include default M when config is nil, got: %s", query)
	}
}

func TestNormalizeHNSWConfig(t *testing.T) {
	defaults := DefaultHNSWConfig()

	tests := []struct {
		name     string
		input    *HNSWConfig
		expected *HNSWConfig
	}{
		{
			name:     "nil uses defaults",
			input:    nil,
			expected: defaults,
		},
		{
			name: "partial uses defaults for missing values",
			input: &HNSWConfig{
				Metric: "cosine",
			},
			expected: &HNSWConfig{
				Metric:         "cosine",
				EfConstruction: defaults.EfConstruction,
				EfSearch:       defaults.EfSearch,
				M:              defaults.M,
			},
		},
		{
			name: "complete config preserved",
			input: &HNSWConfig{
				Metric:         "ip",
				EfConstruction: 512,
				EfSearch:       300,
				M:              32,
			},
			expected: &HNSWConfig{
				Metric:         "ip",
				EfConstruction: 512,
				EfSearch:       300,
				M:              32,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeHNSWConfig(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("normalizeHNSWConfig() = %+v, want %+v", got, tc.expected)
			}
		})
	}
}

func TestParseEmbeddingDimension(t *testing.T) {
	dim, err := parseEmbeddingDimension("FLOAT[384]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dim != 384 {
		t.Fatalf("dim = %d, want 384", dim)
	}

	if _, err := parseEmbeddingDimension("DOUBLE[384]"); err == nil {
		t.Fatalf("expected error for unsupported type")
	}
}

func TestParseHNSWConfigFromCreateIndexSQL(t *testing.T) {
	sql := `
		CREATE INDEX IF NOT EXISTS vec_idx ON vectors
		USING HNSW (embedding)
		WITH (
			metric = 'cosine',
			ef_construction = 256,
			ef_search = 128,
			M = 24
		)
	`

	cfg, err := parseHNSWConfigFromCreateIndexSQL(sql)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := &HNSWConfig{
		Metric:         "cosine",
		EfConstruction: 256,
		EfSearch:       128,
		M:              24,
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Fatalf("parseHNSWConfigFromCreateIndexSQL() = %+v, want %+v", cfg, expected)
	}

	cfg, err = parseHNSWConfigFromCreateIndexSQL("CREATE INDEX vec_idx ON vectors USING HNSW (embedding)")
	if err != nil {
		t.Fatalf("unexpected error for SQL without explicit WITH options: %v", err)
	}

	defaults := DefaultHNSWConfig()
	if !reflect.DeepEqual(cfg, defaults) {
		t.Fatalf("parseHNSWConfigFromCreateIndexSQL() defaults = %+v, want %+v", cfg, defaults)
	}
}

func TestDuckDBStoreInitRuntimeChecks(t *testing.T) {
	t.Run("verifyEmbeddingDimension", func(t *testing.T) {
		cases := []struct {
			name       string
			initDim    int
			verifyDim  int
			expectErr  bool
			errContain string
		}{
			{name: "matching dimension passes", initDim: 384, verifyDim: 384, expectErr: false},
			{name: "mismatched dimension fails", initDim: 384, verifyDim: 768, expectErr: true, errContain: "embedding dimension mismatch"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				store, err := NewDuckDBStoreWithConfig("", tc.initDim, &HNSWConfig{
					Metric:         "l2sq",
					EfConstruction: 128,
					EfSearch:       64,
					M:              16,
				})
				require.NoError(t, err)
				t.Cleanup(func() { _ = store.Close() })

				err = store.verifyEmbeddingDimension(tc.verifyDim)
				if tc.expectErr {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tc.errContain)
					return
				}
				require.NoError(t, err)
			})
		}
	})

	t.Run("ensureHNSWIndex", func(t *testing.T) {
		cases := []struct {
			name      string
			initial   *HNSWConfig
			updated   *HNSWConfig
			expectCfg *HNSWConfig
		}{
			{
				name: "create with initial config",
				initial: &HNSWConfig{
					Metric:         "l2sq",
					EfConstruction: 128,
					EfSearch:       64,
					M:              16,
				},
				updated: nil,
				expectCfg: &HNSWConfig{
					Metric:         "l2sq",
					EfConstruction: 128,
					EfSearch:       64,
					M:              16,
				},
			},
			{
				name: "drift triggers rebuild",
				initial: &HNSWConfig{
					Metric:         "l2sq",
					EfConstruction: 128,
					EfSearch:       64,
					M:              16,
				},
				updated: &HNSWConfig{
					Metric:         "cosine",
					EfConstruction: 256,
					EfSearch:       128,
					M:              24,
				},
				expectCfg: &HNSWConfig{
					Metric:         "cosine",
					EfConstruction: 256,
					EfSearch:       128,
					M:              24,
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				store, err := NewDuckDBStoreWithConfig("", 384, tc.initial)
				require.NoError(t, err)
				t.Cleanup(func() { _ = store.Close() })

				// Ensure runtime path has at least one row to keep index metadata populated in-memory.
				err = store.UpsertVector(context.Background(), "id-1", "file-1", make([]float32, 384))
				require.NoError(t, err)

				if tc.updated != nil {
					store.config = normalizeHNSWConfig(tc.updated)
					require.NoError(t, store.ensureHNSWIndex())
				}

				sqlText, exists, err := store.getExistingHNSWIndexSQL()
				require.NoError(t, err)
				require.True(t, exists, "vec_idx must exist")
				require.NotEmpty(t, strings.TrimSpace(sqlText))

				metric, metricKnown, err := store.getExistingHNSWIndexMetric()
				require.NoError(t, err)
				require.True(t, metricKnown, "vec_idx metric must be available via pragma_hnsw_index_info")
				assert.Equal(t, strings.ToLower(tc.expectCfg.Metric), strings.ToLower(metric))
			})
		}
	})

	t.Run("metric mapping used by search paths", func(t *testing.T) {
		cases := []struct {
			metric       string
			expectMetric DistanceMetric
		}{
			{metric: "l2sq", expectMetric: DistanceEuclidean},
			{metric: "cosine", expectMetric: DistanceCosine},
			{metric: "ip", expectMetric: DistanceInnerProduct},
			{metric: "unknown", expectMetric: DistanceEuclidean},
		}

		for i, tc := range cases {
			t.Run(fmt.Sprintf("case_%d_%s", i, tc.metric), func(t *testing.T) {
				assert.Equal(t, tc.expectMetric, metricNameToDistanceMetric(tc.metric))
			})
		}
	})
}
