package storage

import (
	"context"
	"fly_server/internal/domain"
	"sync"
	"time"
)

type Memory struct {
	mu     sync.RWMutex
	agents map[string]*domain.Agent
	tasks  map[string]*domain.Task
	photos map[string]*domain.Photo
}

func NewMemory() *Memory {
	return &Memory{
		agents: map[string]*domain.Agent{},
		tasks:  map[string]*domain.Task{},
		photos: map[string]*domain.Photo{},
	}
}

func (m *Memory) RegisterAgent(a domain.Agent) {

	m.mu.Lock()
	defer m.mu.Unlock()

	a.LastSeen = time.Now()

	copy := a
	m.agents[a.ID] = &copy
}

func (m *Memory) GetAgent(id string) (*domain.Agent, bool) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	ag, ok := m.agents[id]
	if !ok {
		return nil, false
	}

	copy := *ag
	return &copy, true
}

func (m *Memory) Heartbeat(id string) {

	m.mu.Lock()
	defer m.mu.Unlock()

	if ag, ok := m.agents[id]; ok {
		ag.LastSeen = time.Now()
	}
}

func (m *Memory) ListTasks() map[string]*domain.Task {

	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string]*domain.Task, len(m.tasks))

	for k, v := range m.tasks {

		copy := *v
		out[k] = &copy
	}

	return out
}

func (m *Memory) UpdateResult(t domain.Task) {

	m.mu.Lock()
	defer m.mu.Unlock()

	// if existing, ok := m.tasks[t.ID]; ok {
	// existing.Status = "done"
	// existing.Result = t.Result
	// }
}

func (m *Memory) CreateTask(t domain.Task) {

	// m.mu.Lock()
	// defer m.mu.Unlock()

	// copy := t
	// m.tasks[t.ID] = &copy
}

func (m *Memory) Save(ctx context.Context, photo domain.Photo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	copy := photo
	m.photos[photo.ID] = &copy
	return nil
}

func (m *Memory) Get(ctx context.Context, id string) (*domain.Photo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.photos[id], nil
}

func (m *Memory) GetByOwnerId(ctx context.Context, ownerId string) ([]*domain.Photo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.Photo
	for _, photo := range m.photos {
		if photo.DeviceID == ownerId {
			// создаем копию, чтобы не возвращать указатель на внутренние данные
			copy := *photo
			result = append(result, &copy)
		}
	}

	return result, nil
}

func (m *Memory) ListAgents() ([]*domain.Agent, error) {
	var agents []*domain.Agent
	for _, v := range m.agents {
		agents = append(agents, v)
	}
	return agents, nil
}
