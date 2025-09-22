package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	magostypespb "github.com/Hahn814/magos/proto/magos/v1/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var logLevel = new(slog.LevelVar) // INFO by default
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

type api struct {
	magosapipb.UnimplementedAPIServer
}

func (s *api) Hello(_ context.Context, in *magostypespb.HelloRequest) (*magostypespb.HelloResponse, error) {
	logger.Debug("recieved: %v", "message", in.GetName())
	return &magostypespb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *api) RegisterAgentServer(_ context.Context, in *magostypespb.RegisterAgentServerRequest) (*magostypespb.RegisterAgentServerResponse, error) {
	logger.Debug("register agent server", "agent", in)
	return &magostypespb.RegisterAgentServerResponse{Address: in.GetAddress()}, nil
}

func main() {
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment

	viper.SetEnvPrefix("magos")
	viper.BindEnv("api.port")
	viper.SetDefault("api.port", 50051)
	port := viper.GetInt("api.port")
	addr := viper.GetString("api.addr")

	if addr == "" {
		addr = "0.0.0.0"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	magosapipb.RegisterAPIServer(s, &api{})

	configAttrs := slog.Group("configuration", "port", port, "addr", lis.Addr())
	logger.Info("Starting Magos API..", configAttrs)
	if err := s.Serve(lis); err != nil {
		logger.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
