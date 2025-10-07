package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	magosapipb "github.com/Hahn814/magos/proto/magos/v1/api"
	"google.golang.org/grpc"
)

var logLevel = new(slog.LevelVar)
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

type apiConfig struct {
	magosapipb.UnimplementedAPIServer
}

func NewAPIClient(host string, port int, errorChannel chan error, serverReady chan bool) {

	go func() {
		address := fmt.Sprintf("%s:%d", host, port)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			errorChannel <- err
		}

		grpcServer := grpc.NewServer()
		magosapipb.RegisterAPIServer(grpcServer, &apiConfig{})
		go func() {
			err := grpcServer.Serve(lis)
			if err != nil {
				errorChannel <- err
			}
		}()

		serverReady <- true
	}()
}
