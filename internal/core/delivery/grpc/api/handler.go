package api

import (
	"context"

	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
)

type APIHandler struct {
	magosapipb.UnimplementedAPIServer
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) GetAgent(ctx context.Context, req *magostypespb.GetAgentRequest) (*magostypespb.GetAgentResponse, error) {
	return &magostypespb.GetAgentResponse{}, nil
}

func (h *APIHandler) GetAgents(ctx context.Context, req *magostypespb.GetAgentsRequest) (*magostypespb.GetAgentsResponse, error) {
	return &magostypespb.GetAgentsResponse{}, nil
}

func (h *APIHandler) Hello(ctx context.Context, req *magostypespb.HelloRequest) (*magostypespb.HelloResponse, error) {
	return &magostypespb.HelloResponse{}, nil
}

func (h *APIHandler) RegisterAgentServer(ctx context.Context, req *magostypespb.RegisterAgentServerRequest) (*magostypespb.RegisterAgentServerResponse, error) {
	return &magostypespb.RegisterAgentServerResponse{}, nil
}
