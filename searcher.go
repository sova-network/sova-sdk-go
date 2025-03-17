package sova_sdk_go

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	types "github.com/sova-network/sova-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type SovaSearcher struct {
	client      types.SearcherServiceClient
	accessToken *types.Token
}

func NewSovaSearcher(searcherURL string, caCert *string, domainName *string) (*SovaSearcher, error) {
	return NewSovaSearcherWithAccessToken(searcherURL, caCert, domainName, nil)
}

// NewWithAccessToken function
func NewSovaSearcherWithAccessToken(url string, caPem *string, domainName *string, accessToken *types.Token) (*SovaSearcher, error) {
	var conn *grpc.ClientConn
	var err error

	if caPem != nil && domainName != nil {
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(*caPem)); !ok {
			return nil, errors.New("failed to parse CA certificate")
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs:    certPool,
			ServerName: *domainName,
		})

		conn, err = grpc.NewClient(url, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("failed to connect with TLS: %w", err)
		}
	} else {
		conn, err = grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
	}

	client := &SovaSearcher{
		client:      types.NewSearcherServiceClient(conn),
		accessToken: accessToken,
	}

	return client, nil
}

// Subscribe handles Mempool subscriptions
func (s *SovaSearcher) Subscribe(ctx context.Context, subscription *types.MempoolSubscription, onData func(*types.MempoolPacket)) error {
	ctx = s.addAuthorizationMetadata(ctx)

	stream, err := s.client.SubscribeMempool(ctx, subscription)
	if err != nil {
		return fmt.Errorf("stream subscription error: %v", err)
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

func (s *SovaSearcher) SubscribeByAddress(ctx context.Context, addresses []string, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_Addresses{
			Addresses: &types.AddressSubscriptionV0{
				Address: addresses,
			},
		},
	}, onData)
}

func (s *SovaSearcher) SubscribeByWorkchain(ctx context.Context, workchainId int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_Workchain{
			Workchain: &types.WorkchainSubscriptionV0{
				WorkchainId: workchainId,
			},
		}}, onData)
}

func (s *SovaSearcher) SubscribeByWorkchainShard(ctx context.Context, workchainId int32, shard []byte, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_WorkchainShard{
			WorkchainShard: &types.WorkchainShardSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
			},
		}}, onData)
}

func (s *SovaSearcher) SubscribeByExternalOutMsgBodyOpcode(ctx context.Context, workchainId int32, shard []byte, opcode int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_ExternalOutMessageBodyOpcode{
			ExternalOutMessageBodyOpcode: &types.ExternalOutMessageBodyOpcodeSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
				Opcode:      opcode,
			},
		}}, onData)
}

func (s *SovaSearcher) SubscribeByInternalMsgBodyOpcode(ctx context.Context, workchainId int32, shard []byte, opcode int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_InternalMessageBodyOpcode{
			InternalMessageBodyOpcode: &types.InternalMessageBodyOpcodeSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
				Opcode:      opcode,
			},
		}}, onData)
}

func (s *SovaSearcher) SetAccessToken(token *types.Token) {
	s.accessToken = token
}

func (s *SovaSearcher) SendBundle(ctx context.Context, bundle *types.Bundle) (*types.SendBundleResponse, error) {
	req := bundle

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.SendBundle(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send bundle: %v", err)
	}

	return resp, nil
}

func (s *SovaSearcher) GetTipAddresses(ctx context.Context) (*types.GetTipAddressesResponse, error) {
	req := &types.GetTipAddressesRequest{}

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.GetTipAddresses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tip addresses: %v", err)
	}

	return resp, nil
}

func (s *SovaSearcher) SubscribeBundleResult(ctx context.Context, onData func(*types.BundleResult)) error {
	req := &types.SubscribeBundleResultsRequest{}

	ctx = s.addAuthorizationMetadata(ctx)

	stream, err := s.client.SubscribeBundleResults(ctx, req)
	if err != nil {
		return fmt.Errorf("stream subscription error: %v", err)
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

func (s *SovaSearcher) addAuthorizationMetadata(ctx context.Context) context.Context {
	if s.accessToken == nil {
		return ctx
	}
	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", s.accessToken.GetValue()),
	})
	return metadata.NewOutgoingContext(ctx, md)
}
