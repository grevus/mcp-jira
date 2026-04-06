package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const voyageDefaultURL = "https://api.voyageai.com/v1/embeddings"
const voyageModel = "voyage-3"
const voyageDimension = 1024

// VoyageEmbedder реализует Embedder через Voyage AI API.
type VoyageEmbedder struct {
	apiKey string
	url    string // переопределяется в тестах
	http   *http.Client
}

// NewVoyageEmbedder создаёт клиент.
// Если httpClient == nil — используется http.DefaultClient.
func NewVoyageEmbedder(apiKey string, httpClient *http.Client) *VoyageEmbedder {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &VoyageEmbedder{
		apiKey: apiKey,
		url:    voyageDefaultURL,
		http:   httpClient,
	}
}

// Name возвращает "voyage".
func (e *VoyageEmbedder) Name() string { return "voyage" }

// Dimension возвращает размерность векторов (1024).
func (e *VoyageEmbedder) Dimension() int { return voyageDimension }

// voyageRequest — тело запроса к Voyage AI.
type voyageRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// voyageEmbeddingItem — один элемент из data[].
type voyageEmbeddingItem struct {
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// voyageResponse — ответ Voyage AI.
type voyageResponse struct {
	Data []voyageEmbeddingItem `json:"data"`
}

// Embed выполняет один HTTP-запрос к Voyage AI и возвращает embeddings.
// Порядок результатов гарантирован через поле index из ответа.
// Batching на 128 — Task 28.
func (e *VoyageEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	payload := voyageRequest{
		Input: texts,
		Model: voyageModel,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("voyage: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("voyage: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("voyage: http do: %w", err)
	}
	defer resp.Body.Close()

	var voyageResp voyageResponse
	if err := json.NewDecoder(resp.Body).Decode(&voyageResp); err != nil {
		return nil, fmt.Errorf("voyage: decode response: %w", err)
	}

	// Восстанавливаем порядок по полю index из ответа.
	result := make([][]float32, len(texts))
	for _, item := range voyageResp.Data {
		result[item.Index] = item.Embedding
	}

	return result, nil
}
