package embed

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAIEmbedder_HappyPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/embeddings", r.URL.Path)
		require.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body struct {
			Input      []string `json:"input"`
			Model      string   `json:"model"`
			Dimensions int      `json:"dimensions"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, []string{"hello", "world"}, body.Input)
		require.Equal(t, "text-embedding-3-small", body.Model)
		require.Equal(t, 1024, body.Dimensions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]any{
			"data": []map[string]any{
				{"index": 0, "embedding": []float32{0.1, 0.2, 0.3, 0.4}},
				{"index": 1, "embedding": []float32{0.5, 0.6, 0.7, 0.8}},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	e := NewOpenAIEmbedder("test-key", nil)
	e.baseURL = ts.URL + "/v1/embeddings"

	result, err := e.Embed(t.Context(), []string{"hello", "world"})
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.InDeltaSlice(t, []float32{0.1, 0.2, 0.3, 0.4}, result[0], 1e-6)
	require.InDeltaSlice(t, []float32{0.5, 0.6, 0.7, 0.8}, result[1], 1e-6)
}

func TestOpenAIEmbedder_Name(t *testing.T) {
	e := NewOpenAIEmbedder("key", nil)
	require.Equal(t, "openai", e.Name())
}

func TestOpenAIEmbedder_Dimension(t *testing.T) {
	e := NewOpenAIEmbedder("key", nil)
	require.Equal(t, 1024, e.Dimension())
}
