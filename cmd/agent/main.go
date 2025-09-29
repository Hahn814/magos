package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"time"

	magosagentpb "github.com/Hahn814/magos/proto/magos/v1/agent"
	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logLevel = new(slog.LevelVar) // INFO by default
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
var id = randomString(10)

type agent struct {
	magosagentpb.UnimplementedAgentServer
}

type api struct {
	magosapipb.UnimplementedAPIServer
	Address string
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func (s *agent) Hello(_ context.Context, in *magostypespb.HelloRequest) (*magostypespb.HelloResponse, error) {
	logger.Debug("recieved: %v", "message", in.GetName())
	return &magostypespb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *agent) Describe(_ context.Context, in *magostypespb.DescribeAgentRequest) (*magostypespb.DescribeAgentResponse, error) {
	logger.Debug("recieved", "request", in)
	hostname, err := os.Hostname()
	if err != nil {
		return &magostypespb.DescribeAgentResponse{}, err
	}

	return &magostypespb.DescribeAgentResponse{Id: id, Hostname: hostname}, nil
}

func (s *api) registerAgentServer(_ context.Context, in *magostypespb.RegisterAgentServerRequest) {

	// TODO: Bind retry policy params to environment
	retryPolicy := fmt.Sprintf(`{
		"methodConfig": [
			{
			"name": [
				{
				"service": "api.API",
				"method": "RegisterAgentServer"
				}
			],
			"timeout": "%ds",
			"retryPolicy": {
				"maxAttempts": 5,
				"initialBackoff": "10s",
				"maxBackoff": "180s",
				"backoffMultiplier": 6,
				"retryableStatusCodes": [
				"UNAVAILABLE",
				"DEADLINE_EXCEEDED"
				]
			}
			}
		]
	}`, int(5*time.Minute.Seconds()))

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithMaxCallAttempts(5),
	}
	conn, err := grpc.NewClient(s.Address, dialOpts...)
	if err != nil {
		logger.Error("did not connect", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := magosapipb.NewAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	r, err := client.RegisterAgentServer(ctx, in)
	if err != nil {
		logger.Error("could not reigster agent", "error", err)
		os.Exit(2)
	}

	logger.Info("registered agent", "success", r.GetSuccess())
}

func findOpenPrivatePort(minPort int, maxPort int) (int, error) {
	for port := minPort; port <= maxPort; port++ {
		logger.Debug("checking", "port", port)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}

	}

	// TODO: make this port range configurable
	return 0, fmt.Errorf("could not find open port in the %d-%d range", minPort, maxPort)
}

func main() {
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment

	agentReady := make(chan bool)
	agentError := make(chan error)

	go func() {
		// TODO: refactor to be reciever function provided by the api and agent types
		addr := "0.0.0.0"
		port, err := findOpenPrivatePort(49152, 65535)
		if err != nil {
			logger.Error("Failed to find open port", "error", err)
			os.Exit(1)
		}

		agent_addr := fmt.Sprintf("%s:%d", addr, port)
		lis, err := net.Listen("tcp", agent_addr)
		if err != nil {
			logger.Error("Failed to listen", "error", err)
			os.Exit(1)
		}
		s := grpc.NewServer()
		magosagentpb.RegisterAgentServer(s, &agent{})

		go func() {
			configAttrs := slog.Group("configuration", "port", port, "addr", lis.Addr())
			logger.Info("Starting Magos agent service..", configAttrs)

			// TODO: get api server addr from env
			// TODO: get api server port from env
			api := &api{
				Address: fmt.Sprintf("%s:%d", addr, 50051),
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			api.registerAgentServer(ctx, &magostypespb.RegisterAgentServerRequest{Address: agent_addr, Id: id})
			if err := s.Serve(lis); err != nil {
				logger.Error("Failed to serve", "error", err)
				agentError <- err
			}

			agentReady <- true // TODO: replace with actual readiness test
		}()

	}()

	select {
	case <-agentReady:
		logger.Info("agent is ready")
	case err := <-agentError:
		logger.Error("failed to start agent", "error", err)
		os.Exit(1)
	case <-time.After(time.Hour):
		logger.Error("failed to start agent after time.")
		os.Exit(1)
	}

	select {}

}
