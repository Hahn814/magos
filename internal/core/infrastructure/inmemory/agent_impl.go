package infrastructure

import (
	"context"
	"errors"

	"github.com/Hahn814/magos/internal/core/domain"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
)

type AgentRepositoryImpl struct {
	_    domain.AgentRepository
	data []domain.Agent
}

func (r *AgentRepositoryImpl) CreateAgent(ctx context.Context, agent domain.Agent) error {
	r.data = append(r.data, agent)
	return nil
}

func (r *AgentRepositoryImpl) GetAgentByID(ctx context.Context, ID string) (*domain.Agent, error) {
	for _, agent := range r.data {
		if agent.ID == ID {
			return &agent, nil
		}
	}

	return &domain.Agent{}, ErrAgentNotFound
}
