package http

import (
	"encoding/json"
	"fly_server/internal/domain"
	"fly_server/internal/infra/sshworker"
	"fly_server/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store     storage.Storage
	log       *log.Logger
	sshWorker *sshworker.Client
}

var ssh = sshworker.New(
	"192.168.0.100",
	"user",
	"/home/prometheus/.ssh/id_rsa",
	"/home/user/kopters/src/drone",
)

func NewHandler(store storage.Storage, logger *log.Logger) *Handler {
	return &Handler{
		store:     store,
		log:       logger,
		sshWorker: ssh,
	}
}

// RegisterAgent godoc
//
//	@Summary		Register agent
//	@Description	Register new drone agent
//	@Tags			agents
//	@Accept			json
//	@Produce		json
//	@Param			agent	body		domain.Agent	true	"Agent data"
//	@Success		200		{object}	domain.Agent
//	@Router			/agents/register [post]
func (h *Handler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	var a domain.Agent

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	h.store.RegisterAgent(a)

	json.NewEncoder(w).Encode(a)
}

// Heartbeat godoc
//
//	@Summary	Agent heartbeat
//	@Tags		agents
//	@Param		id	path	string	true	"Agent ID"
//	@Success	200
//	@Router		/agents/{id}/heartbeat [post]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing agent id", http.StatusBadRequest)
		return
	}

	h.store.Heartbeat(id)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SubmitResult(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var t domain.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Println("submit result decode error:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.store.UpdateResult(t)

	w.WriteHeader(http.StatusOK)
}

// CreateTask godoc
//
//	@Summary	Create task
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		task	body		domain.Task	true	"Task"
//	@Success	200		{object}	domain.Task
//	@Router		/tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t domain.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Println("create task decode error:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// 1. сохранить в БД
	h.store.CreateTask(t)

	// 2. сериализуем сразу в JSON
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		h.log.Println("marshal error:", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	// 3. отправляем на дрон (без локального файла)
	err = h.sshWorker.SendTaskAndExecute(data)
	if err != nil {
		h.log.Println("ssh execute error:", err)
		http.Error(w, "ssh error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {

	tasks := h.store.ListTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// ListAgents godoc
// @Summary      List all agents
// @Description  Returns a list of all registered drone agents with their details and last seen timestamp
// @Tags         agents
// @Produce      json
// @Success      200  {array}   domain.Agent
// @Failure      500  {string}  string "Internal error"
// @Router       /agents [get]
func (h *Handler) ListAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.store.ListAgents()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}
