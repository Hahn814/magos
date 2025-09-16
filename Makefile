up: build
	legion

build: protoc
	go build -o /home/legion/go/bin/legion ./cmd/legion/
	go build -o /home/legion/go/bin/legionctl ./cmd/legionctl/

protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/legion/v1/legion.proto
