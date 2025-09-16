package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	legionpb "github.com/Hahn814/legion/proto/legion/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type legionServer struct {
	legionpb.UnimplementedLegionServer
}

func main() {
	var logLevel = new(slog.LevelVar) // INFO by default
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	logLevel.Set(slog.LevelDebug) // TODO: bind log level to environment

	viper.SetEnvPrefix("legion")
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
	legionpb.RegisterLegionServer(s, &legionServer{})

	configAttrs := slog.Group("configuration", "port", port, "addr", lis.Addr())
	logger.Info("Starting legion..", configAttrs)
	if err := s.Serve(lis); err != nil {
		logger.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
