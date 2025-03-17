package sova_sdk_go_test

import (
	"context"
	"encoding/base64"
	"testing"

	types "github.com/sova-network/sova-sdk-go/generated"
	"github.com/stretchr/testify/assert"

	sova "github.com/sova-network/sova-sdk-go"
)

func TestSuccessfulClient(t *testing.T) {
	client := sova.NewTestnetClient()

	privKeyB64 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	privateKey, err := base64.StdEncoding.DecodeString(privKeyB64)
	assert.NoError(t, err)

	ctx := context.Background()
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
	err = searcher.SubscribeByWorkchain(ctx, 0, func(mp *types.MempoolPacket) {
		t.Log("reseive mempool packet")
		mempoolCh <- mp
	})
	assert.NoError(t, err)

	// Subscribe to bundle results
	err = searcher.SubscribeBundleResult(ctx, func(br *types.BundleResult) {
		t.Log("reseive bundle result")
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
			assert.NoErrorf(t, err, "failed to send bundle: %v", err)
			assert.NotNilf(t, res, "response is nil")
			t.Logf("got bundle response: %v", res)

		// Wait for bundle result
		case res := <-resultCh:
			t.Logf("got bundle result")
			t.Logf("result: %v", res)
			assert.NotNil(t, res)
			return
		}
	}

}
