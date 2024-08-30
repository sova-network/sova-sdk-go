package mevton_sdk_go

import (
	"context"
	"fmt"
	"log"

	types "github.com/mevton-labs/mevton-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Searcher struct {
	client      types.SearcherServiceClient
	accessToken *types.Token
}

// TODO use conn as a parameter; conn *grpc.ClientConn
func NewSearcher(searcherURL string) (*Searcher, error) {
	conn, err := grpc.NewClient(searcherURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to searcher service: %v", err)
	}

	client := types.NewSearcherServiceClient(conn)

	return &Searcher{
		client:      client,
		accessToken: nil,
	}, nil
}

func (s *Searcher) SetAccessToken(token *types.Token) {
	s.accessToken = token
}

func (s *Searcher) SubscribeMempool(ctx context.Context, addresses []string, onData func(*types.MempoolPacket)) error {
	req := &types.MempoolSubscription{
		Addresses: &types.AddressSubscriptionV0{
			Address: addresses,
		},
	}

	if len(addresses) == 0 {
		req.Addresses = nil
	}

	ctx = s.addAuthorizationMetadata(ctx)

	stream, err := s.client.SubscribeMempool(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to subscribe to mempool: %v", err)
	}

	go func() {
		for {
			packet, err := stream.Recv()
			if err != nil {
				log.Printf("stream error: %v", err)
				return
			}
			onData(packet)
		}
	}()

	return nil
}

func (s *Searcher) SendBundle(ctx context.Context, bundle *types.Bundle) (*types.SendBundleResponse, error) {
	req := bundle

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.SendBundle(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send bundle: %v", err)
	}

	return resp, nil
}

func (s *Searcher) GetTipAddresses(ctx context.Context) (*types.GetTipAddressesResponse, error) {
	req := &types.GetTipAddressesRequest{}

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.GetTipAddresses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tip addresses: %v", err)
	}

	return resp, nil
}

func (s *Searcher) addAuthorizationMetadata(ctx context.Context) context.Context {
	if s.accessToken == nil {
		return ctx
	}
	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", s.accessToken.GetValue()),
	})
	return metadata.NewOutgoingContext(ctx, md)
}
