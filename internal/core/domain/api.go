package domain

import (
	"context"
)

type APIRepository interface {
	GetAgentByID(ctx context.Context, req, ID string) (*Agent, error)
	GetAgents(ctx context.Context) ([]*Agent, error)
	Hello(ctx context.Context) (string, error)
	RegisterAgentServer(ctx context.Context, agent Agent) error
}
