package domain

import "context"

type Agent struct {
	ID string
}

type AgentRepository interface {
	Create(ctx context.Context, agent Agent) error
	FindByID(ctx context.Context, id string) (*Agent, error)
}
