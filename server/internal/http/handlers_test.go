package http

import (
	"bytes"
	"encoding/json"
	"fly_server/internal/domain"
	"fly_server/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHandler() *Handler {

	store := storage.NewMemory()

	logger := log.New(io.Discard, "", 0)

	return &Handler{
		store: store,
		log:   logger,
	}
}

func TestRegisterAgent(t *testing.T) {

	h := newTestHandler()

	agent := domain.Agent{
		ID: "agent-1",
	}

	body, _ := json.Marshal(agent)

	req := httptest.NewRequest(
		http.MethodPost,
		"/agents/register",
		bytes.NewBuffer(body),
	)

	rec := httptest.NewRecorder()

	h.RegisterAgent(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestMemoryRegisterAgent(t *testing.T) {

	store := storage.NewMemory()

	store.RegisterAgent(domain.Agent{ID: "a1"})

	// TODO: проверить state
}
