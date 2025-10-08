package usecase

import (
	"context"

	"github.com/Hahn814/magos/internal/core/domain"
)

type GetAgentUseCase struct {
	repo domain.AgentRepository
}

func NewGetAgentUseCase(repo domain.AgentRepository) *GetAgentUseCase {
	return &GetAgentUseCase{repo: repo}
}

func (uc *GetAgentUseCase) GetAgentByID(ID string) (*domain.Agent, error) {
	return uc.repo.FindByID(context.Background(), ID)
}
