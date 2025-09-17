package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	magospb "github.com/Hahn814/magos/proto/magos/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var logLevel = new(slog.LevelVar) // INFO by default
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

type agent struct {
	magospb.UnimplementedAgentServer
}

func (s *agent) Hello(_ context.Context, in *magospb.HelloRequest) (*magospb.HelloResponse, error) {
	logger.Debug("recieved: %v", "message", in.GetName())
	return &magospb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

func main() {
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment

	viper.SetEnvPrefix("magos")
	viper.BindEnv("port")
	viper.SetDefault("port", 50051)
	port := viper.GetInt("port")
	addr := viper.GetString("addr")

	if addr == "" {
		addr = "0.0.0.0"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	magospb.RegisterAgentServer(s, &agent{})

	configAttrs := slog.Group("configuration", "port", port, "addr", lis.Addr())
	logger.Info("Starting Magos agent service..", configAttrs)
	if err := s.Serve(lis); err != nil {
		logger.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
