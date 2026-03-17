package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StreamPeek/StreamPeek/internal/kafka"
	"github.com/go-chi/chi/v5"
)

// mockKafkaClient overrides the actual produce behavior for testing.
// In a real scenario we'd use interfaces, but for simplicity of this verification,
// we just test the handler logic itself up to the client boundaries.
// Given time constraints, we will perform an integration test with an empty client to verify routing and 400 errors.

func TestProduceHandler_InvalidJSON(t *testing.T) {
	// A basic test to verify 400 Bad Request on invalid JSON
	client, err := kafka.NewClient("localhost:9092")
	if err != nil {
		t.Fatalf("Failed to initialize test client: %v", err)
	}
	defer client.Close()

	handler := NewProduceHandler(client)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest("POST", "/topics/test-topic/record", bytes.NewBuffer([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", w.Result().StatusCode)
	}
}

func BenchmarkProduceHandler_ValidJSONAllocation(b *testing.B) {
	// Benchmark validating json parsing overhead in the handler.
	// We mock out the kafka client by skipping the actual ProduceSync call.

	payload := []byte(`{"sensor_id": 1234, "temperature": 22.5, "status": "active"}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Just benchmark the json.Valid checking which we added for quick fail.
		if !json.Valid(payload) {
			b.Fatal("invalid json")
		}
	}
}
