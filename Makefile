submodules:
	git submodule update --init --recursive

lint:
	golangci-lint run --out-format=tab -n
test:
	go test -v -cover -race ./...

protoc:
	protoc --proto_path=grpc/proto --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative grpc/proto/*

.PHONY: lint test protoc submodules
