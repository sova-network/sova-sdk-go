package mevton_sdk_go

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"

	types "github.com/mevton-labs/mevton-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type MevtonSearcher struct {
	client      types.SearcherServiceClient
	accessToken *types.Token
}

func NewMevtonSearcher(searcherURL string, caCert *[]byte, domainName *string) (*MevtonSearcher, error) {
	var opts []grpc.DialOption

	if caCert != nil && domainName != nil {
		// Set up TLS using the provided CA certificate and domain name
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(*caCert)) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig := &tls.Config{
			RootCAs:    certPool,
			ServerName: *domainName,
		}

		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(searcherURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to searcher service: %v", err)
	}
	client := types.NewSearcherServiceClient(conn)

	return &MevtonSearcher{
		client:      client,
		accessToken: nil,
	}, nil
}

// Subscribe handles Mempool subscriptions
func (s *MevtonSearcher) Subscribe(ctx context.Context, subscription *types.MempoolSubscription, onData func(*types.MempoolPacket)) error {
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

func (s *MevtonSearcher) SubscribeByAddress(ctx context.Context, addresses []string, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_Addresses{
			Addresses: &types.AddressSubscriptionV0{
				Address: addresses,
			},
		},
	}, onData)
}

func (s *MevtonSearcher) SubscribeByWorkchain(ctx context.Context, workchainId int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_Workchain{
			Workchain: &types.WorkchainSubscriptionV0{
				WorkchainId: workchainId,
			},
		}}, onData)
}

func (s *MevtonSearcher) SubscribeByWorkchainShard(ctx context.Context, workchainId int32, shard []byte, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_WorkchainShard{
			WorkchainShard: &types.WorkchainShardSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
			},
		}}, onData)
}

func (s *MevtonSearcher) SubscribeByExternalOutMsgBodyOpcode(ctx context.Context, workchainId int32, shard []byte, opcode int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_ExternalOutMessageBodyOpcode{
			ExternalOutMessageBodyOpcode: &types.ExternalOutMessageBodyOpcodeSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
				Opcode:      opcode,
			},
		}}, onData)
}

func (s *MevtonSearcher) SubscribeByInternalMsgBodyOpcode(ctx context.Context, workchainId int32, shard []byte, opcode int32, onData func(*types.MempoolPacket)) error {
	return s.Subscribe(ctx, &types.MempoolSubscription{
		Subscription: &types.MempoolSubscription_InternalMessageBodyOpcode{
			InternalMessageBodyOpcode: &types.InternalMessageBodyOpcodeSubscriptionV0{
				WorkchainId: workchainId,
				Shard:       shard,
				Opcode:      opcode,
			},
		}}, onData)
}

func (s *MevtonSearcher) SetAccessToken(token *types.Token) {
	s.accessToken = token
}

func (s *MevtonSearcher) SendBundle(ctx context.Context, bundle *types.Bundle) (*types.SendBundleResponse, error) {
	req := bundle

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.SendBundle(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send bundle: %v", err)
	}

	return resp, nil
}

func (s *MevtonSearcher) GetTipAddresses(ctx context.Context) (*types.GetTipAddressesResponse, error) {
	req := &types.GetTipAddressesRequest{}

	ctx = s.addAuthorizationMetadata(ctx)

	resp, err := s.client.GetTipAddresses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tip addresses: %v", err)
	}

	return resp, nil
}

func (s *MevtonSearcher) addAuthorizationMetadata(ctx context.Context) context.Context {
	if s.accessToken == nil {
		return ctx
	}
	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", s.accessToken.GetValue()),
	})
	return metadata.NewOutgoingContext(ctx, md)
}
