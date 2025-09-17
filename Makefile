up: build
	magosapi

build: protoc
	go build -o /home/magos/go/bin/magosapi ./cmd/api/
	go build -o /home/magos/go/bin/magosctl ./cmd/cli/

protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/magos/v1/daemon.proto
