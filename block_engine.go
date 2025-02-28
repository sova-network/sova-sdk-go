package mevton_sdk_go

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	types "github.com/sova-network/sova-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type BlockEngineClient struct {
	client      types.BlockEngineValidatorClient
	accessToken *types.Token
}

func NewBlockEngine(blockEngineURL string, accessToken *types.Token) (*BlockEngineClient, error) {
	conn, err := grpc.NewClient(blockEngineURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to block engine: %v", err)
	}
	client := types.NewBlockEngineValidatorClient(conn)
	return &BlockEngineClient{
		client:      client,
		accessToken: accessToken,
	}, nil
}

func (b *BlockEngineClient) StreamMempool(ctx context.Context, stream <-chan *types.MempoolPacket) error {
	ctx = b.addAuthorizationMetadata(ctx)

	clientStream, err := b.client.StreamMempool(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client stream: %v", err)
	}

	for packet := range stream {
		if err := clientStream.Send(packet); err != nil {
			return fmt.Errorf("failed to send packet: %v", err)
		}
	}

	if err := clientStream.CloseSend(); err != nil {
		return fmt.Errorf("failed to close stream: %v", err)
	}

	return nil
}

func (b *BlockEngineClient) SubscribeBundles(ctx context.Context, onData func(*types.Bundle)) error {
	ctx = b.addAuthorizationMetadata(ctx)

	stream, err := b.client.SubscribeBundles(ctx, &types.SubscribeBundlesRequest{})
	if err != nil {
		return fmt.Errorf("failed to subscribe to bundles: %v", err)
	}

	go func() {
		for {
			bundle, err := stream.Recv()
			if err != nil {
				// TODO handle error
				if !errors.Is(err, grpc.ErrServerStopped) {
					log.Printf("stream error: %v", err)
				}
				if err == io.EOF {
					log.Println("stream closed by server")
					return
				}
				if err != nil {
					log.Printf("stream error: %v", err)
					return
				}
				return
			}
			onData(bundle)
		}
	}()

	return nil
}

func (b *BlockEngineClient) addAuthorizationMetadata(ctx context.Context) context.Context {
	if b.accessToken == nil {
		return ctx
	}
	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", b.accessToken.GetValue()),
	})
	return metadata.NewOutgoingContext(ctx, md)
}
