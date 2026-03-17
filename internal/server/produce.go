package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/StreamPeek/StreamPeek/internal/kafka"
)

// ProduceHandler handles requests to produce records to Kafka.
type ProduceHandler struct {
	client *kafka.Client
}

// NewProduceHandler creates a new handle.
func NewProduceHandler(client *kafka.Client) *ProduceHandler {
	return &ProduceHandler{
		client: client,
	}
}

// RegisterRoutes registers the produce routes onto the provided router.
func (h *ProduceHandler) RegisterRoutes(r chi.Router) {
	r.Post("/topics/{topic}/record", h.HandleProduce)
}

// ProduceResponse represents the JSON output of a successful produce request.
type ProduceResponse struct {
	Topic      string `json:"topic"`
	Partition  int32  `json:"partition"`
	Offset     int64  `json:"offset"`
	BrokerID   int32  `json:"broker_id"`
	DurationMs int64  `json:"duration_ms"`
	ReadURL    string `json:"read_url"`
}

// ErrorResponse represents a JSON error.
type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *ProduceHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

// HandleProduce processes the POST /topics/{topic}/record request.
func (h *ProduceHandler) HandleProduce(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")
	if topic == "" {
		h.writeError(w, http.StatusBadRequest, "topic parameter is required")
		return
	}

	keyStr := r.URL.Query().Get("key")
	var key []byte
	if keyStr != "" {
		key = []byte(keyStr)
	}

	// Read body as pure JSON raw value
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		h.writeError(w, http.StatusBadRequest, "empty request body")
		return
	}

	// Validate it's actually JSON if we want to be strict, though technically
	// we just pass it to Kafka. We'll do a quick check to fail fast on malformed requests.
	if !json.Valid(body) {
		h.writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	// Extract headers
	var kafkaHeaders []kgo.RecordHeader
	for headerKey, values := range r.Header {
		if strings.HasPrefix(strings.ToLower(headerKey), "x-kafka-") {
			k := headerKey[len("x-kafka-"):] // strip prefix
			if len(values) > 0 {
				kafkaHeaders = append(kafkaHeaders, kgo.RecordHeader{
					Key:   k,
					Value: []byte(values[0]), // take first value
				})
			}
		}
	}

	// Produce
	res, err := h.client.ProduceSync(r.Context(), topic, key, body, kafkaHeaders)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to produce message: %v", err))
		return
	}

	// Build success response
	resp := ProduceResponse{
		Topic:      res.Topic,
		Partition:  res.Partition,
		Offset:     res.Offset,
		BrokerID:   res.BrokerID,
		DurationMs: res.DurationMs,
		ReadURL:    fmt.Sprintf("/topics/%s/record?partition=%d&offset=%d", res.Topic, res.Partition, res.Offset),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Just log, we've already written the status
		fmt.Printf("failed to encode response: %v\n", err)
	}
}
