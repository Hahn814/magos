up: build
	magosapi

build: protoc
	go build -o /home/magos/go/bin/magosapi ./cmd/api/
	go build -o /home/magos/go/bin/magosctl ./cmd/cli/
	go build -o /home/magos/go/bin/magosagent ./cmd/agent/

protoc:
	protoc -I=proto --go_out=proto/ --go_opt=paths=source_relative --go-grpc_out=proto/ --go-grpc_opt=paths=source_relative proto/magos/v1/types/types.proto
	protoc -I=proto --go_out=proto/ --go_opt=paths=source_relative --go-grpc_out=proto/ --go-grpc_opt=paths=source_relative proto/magos/v1/api/api.proto
	protoc -I=proto --go_out=proto/ --go_opt=paths=source_relative --go-grpc_out=proto/ --go-grpc_opt=paths=source_relative proto/magos/v1/agent/agent.proto
