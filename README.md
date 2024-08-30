# Mevton SDK

Mevton SDK is a Go library for interacting with the Mevton MEV Block Engine and Searcher services.
It provides functionalities for authentication, streaming mempool transactions, subscribing to bundles, and sending bundles.

## Features

- **Authentication**: Authenticate using your private key to obtain access and refresh tokens.
- **Streaming Mempool Transactions**: Stream transactions from the client to the Mevton MEV Block Engine.
- **Subscribe to Bundles**: Subscribe to receive a stream of simulated and profitable bundles.
- **Send Bundles**: Send bundles to the Mevton MEV Block Engine for processing.

## Installation

Update git submodules to get the latest version of the SDK.

```bash 
  git submodule update --init --recursive
```

Generate protobuf files from the `proto` directory using the following command:

```bash
make protoc
or 
protoc -I=grpc/proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative grpc/proto/*
or
protoc --proto_path=grpc/proto --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative grpc/proto/*
```

Add dependencies to your Go project using the following command:

```bash
go get -u github.com/mevton-labs/mevton-sdk-go
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with any changes or enhancements. Follow these steps to contribute:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes and commit them with clear and descriptive messages.
4. Push your changes to your fork.
5. Open a pull request to the main repository.
