package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Hahn814/magos/cmd/api/app"
	magosagentpb "github.com/Hahn814/magos/proto/magos/v1/agent"
	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logLevel = new(slog.LevelVar) // INFO by default
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

var agents = []*magostypespb.Agent{}

type api struct {
	magosapipb.UnimplementedAPIServer
}

func (s *api) Hello(_ context.Context, in *magostypespb.HelloRequest) (*magostypespb.HelloResponse, error) {
	logger.Debug("recieved: %v", "message", in.GetName())
	return &magostypespb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *api) RegisterAgentServer(_ context.Context, in *magostypespb.RegisterAgentServerRequest) (*magostypespb.RegisterAgentServerResponse, error) {
	logger.Debug("register agent server", "agent", in)
	agents = append(agents, &magostypespb.Agent{Id: in.GetId(), Hostname: in.GetAddress()})
	return &magostypespb.RegisterAgentServerResponse{Success: true}, nil
}

func getAgentMetadata(address string) (*magostypespb.GetAgentResponse, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, dialOpts...)
	if err != nil {
		logger.Error("did not connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := magosagentpb.NewAgentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	r, err := client.Describe(ctx, &magostypespb.DescribeAgentRequest{})
	if err != nil {
		logger.Error("could not reigster agent", "error", err)
		os.Exit(2)
	}

	return &magostypespb.GetAgentResponse{Agent: &magostypespb.Agent{Id: r.GetId(), Hostname: r.GetHostname()}}, nil
}

func (s *api) GetAgent(_ context.Context, in *magostypespb.GetAgentRequest) (*magostypespb.GetAgentResponse, error) {
	logger.Debug("recieved", "request", in)

	for _, agent := range agents {
		if agent.GetId() == in.Agent.Id {
			response, err := getAgentMetadata(agent.GetHostname())
			if err != nil {
				logger.Error("could not get agent", "err", err)
				return response, err
			}

			return response, nil
		}
	}

	return nil, fmt.Errorf("no valid agent id provided")
}

func (s *api) GetAgents(_ context.Context, in *magostypespb.GetAgentsRequest) (*magostypespb.GetAgentsResponse, error) {
	logger.Debug("recieved", "request", in)
	return &magostypespb.GetAgentsResponse{Agents: agents}, nil
}

func main() {
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment

	serverError := make(chan error)
	serverReady := make(chan bool)
	addr := viper.GetString("api.addr")
	port := viper.GetInt("api.port")

	app.NewAPIClient(addr, port, serverError, serverReady)

	select {
	case <-serverReady:
		logger.Info("server ready")
	case err := <-serverError:
		logger.Error("error occurred starting API server", "error", err)
		os.Exit(1)
	case <-time.After(5 * time.Second):
		logger.Error("server did not start within the expected time")
		os.Exit(2)
	}

	select {}

}
