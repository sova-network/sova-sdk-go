package sova_sdk_go_test

import (
	"context"
	"encoding/base64"
	"testing"

	types "github.com/sova-network/sova-sdk-go/generated"
	"github.com/stretchr/testify/assert"

	sova "github.com/sova-network/sova-sdk-go"
)

func TestTestnetSuccessfulClient(t *testing.T) {
	ctx := context.Background()
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := sova.NewTestnetClient()

	privKeyB64 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	privateKey, err := base64.StdEncoding.DecodeString(privKeyB64)
	assert.NoError(t, err)

	// Authenticate
	token, err := client.Authenticate(ctx, privateKey)

	// Expect error: rpc error: code = Unauthenticated desc = Non-whitelisted public key
	assert.Error(t, err)
	assert.ErrorContains(t, err, "rpc error: code = Unauthenticated desc = Non-whitelisted public key")

	// Create a new searcher
	searcher, err := client.Searcher(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, searcher)
	searcher.SetAccessToken(token)

	mempoolCh := make(chan *types.MempoolPacket, 1)
	resultCh := make(chan *types.BundleResult, 1)

	// Subscribe to mempool packets
	err = searcher.SubscribeByWorkchain(subCtx, 0, func(mp *types.MempoolPacket) {
		t.Log("receive mempool packet")
		mempoolCh <- mp
	})
	assert.NoError(t, err)

	// Subscribe to bundle results
	err = searcher.SubscribeBundleResult(subCtx, func(br *types.BundleResult) {
		t.Log("receive bundle result")
		resultCh <- br
	})
	assert.NoError(t, err)

	for {
		select {
		// Wait for mempool packet
		case msg := <-mempoolCh:
			t.Logf("got mempool packet: %v", len(msg.ExternalMessages))

			// Send bundle
			res, err := searcher.SendBundle(ctx, &types.Bundle{
				Message: []*types.ExternalMessage{ // Correct type
					{
						Data: []byte(msg.ExternalMessages[0].Data),
					},
				},
				ExpirationNs: nil,
			})

			assert.Error(t, err)
			assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = Invalid bundle, tip amount is 0")
			assert.Nil(t, res)

			cancel()
			return
		}
	}
}

func TestMainnetSuccessfulClient(t *testing.T) {
	ctx := context.Background()
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := sova.NewMainnetClient()

	privKeyB64 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	privateKey, err := base64.StdEncoding.DecodeString(privKeyB64)
	assert.NoError(t, err)

	// Authenticate
	token, err := client.Authenticate(ctx, privateKey)

	// Expect error: rpc error: code = Unauthenticated desc = Non-whitelisted public key
	assert.Error(t, err)
	assert.ErrorContains(t, err, "rpc error: code = Unauthenticated desc = Non-whitelisted public key")

	// Create a new searcher
	searcher, err := client.Searcher(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, searcher)
	searcher.SetAccessToken(token)

	mempoolCh := make(chan *types.MempoolPacket, 1)
	resultCh := make(chan *types.BundleResult, 1)

	// Subscribe to mempool packets
	err = searcher.SubscribeByWorkchain(subCtx, 0, func(mp *types.MempoolPacket) {
		t.Log("receive mempool packet")
		mempoolCh <- mp
	})
	assert.NoError(t, err)

	// Subscribe to bundle results
	err = searcher.SubscribeBundleResult(subCtx, func(br *types.BundleResult) {
		t.Log("receive bundle result")
		resultCh <- br
	})
	assert.NoError(t, err)

	for {
		select {
		// Wait for mempool packet
		case msg := <-mempoolCh:
			t.Logf("got mempool packet: %v", len(msg.ExternalMessages))

			// Send bundle
			res, err := searcher.SendBundle(ctx, &types.Bundle{
				Message: []*types.ExternalMessage{ // Correct type
					{
						Data: []byte(msg.ExternalMessages[0].Data),
					},
				},
				ExpirationNs: nil,
			})

			assert.Error(t, err)
			assert.ErrorContains(t, err, "rpc error: code = InvalidArgument desc = Invalid bundle, tip amount is 0")
			assert.Nil(t, res)

			cancel()
			return
		}
	}
}
