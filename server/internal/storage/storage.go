package storage

import "fly_server/internal/domain"

type Storage interface {
	RegisterAgent(domain.Agent)
	Heartbeat(id string)

	CreateTask(domain.Task)
	ListTasks() map[string]*domain.Task
	UpdateResult(domain.Task)

	ListAgents() ([]*domain.Agent, error)
}
