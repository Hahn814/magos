package grpc

import (
	"context"
	"os"

	"github.com/Hahn814/magos/internal/core/domain"
	magosagentpb "github.com/Hahn814/magos/proto/magos/v1/agent"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
)

type AgentHandler struct {
	magosagentpb.UnimplementedAgentServer
	repo domain.AgentRepository
}

func NewAgentHandler(repo domain.AgentRepository) *AgentHandler {
	return &AgentHandler{repo: repo}
}

func (h *AgentHandler) Describe(ctx context.Context, req *magostypespb.DescribeAgentRequest) (*magostypespb.DescribeAgentResponse, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &magostypespb.DescribeAgentResponse{}, err
	}

	// TODO: get agent ID from repo
	return &magostypespb.DescribeAgentResponse{Id: "foo", Hostname: hostname}, nil
}
