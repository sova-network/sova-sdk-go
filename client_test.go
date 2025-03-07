package sova_sdk_go

import (
	"context"
	"testing"

	types "github.com/sova-network/sova-sdk-go/generated"
	"gotest.tools/assert"
)

func TestClient(t *testing.T) {
	Once.Do(RunServer)
	t.Run("NewMainnetClient", func(t *testing.T) {
		t.Skip("not released")
	})
	t.Run("NewTestnetClient", func(t *testing.T) {
		NewTestnetClient()
	})
	t.Run("NewCustomClient", func(t *testing.T) {
		NewCustomClient("url", nil, nil, nil)
	})
	t.Run("Authenticate", func(t *testing.T) {
		client:= NewCustomClient("localhost:50051", nil, nil, nil)
		privateKey := []byte{
			155, 202, 118, 43, 82, 100, 113, 150, 99, 21,
			45, 230, 88, 247, 193, 12, 92, 78, 191, 229,
			73, 191, 100, 156, 231, 41, 144, 54, 202, 199,
			75, 1,
		}

		token, err := client.Authenticate(context.Background(), privateKey)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if token == nil {
			t.Errorf("token is nil")
		}

		assert.Equal(t, token.Value, "access_token")  
	})
	t.Run("Searcher", func(t *testing.T) {
		client:= NewCustomClient("localhost:50051", nil, nil, nil)
		searcher, err := client.Searcher(context.Background())
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if searcher == nil {
			t.Error("searcher is nil")
		}

		searcher.SetAccessToken(&types.Token{Value: "access_token"})
		searcher.SubscribeByWorkchain(context.Background(), 0, func(mp *types.MempoolPacket) {
			t.Logf("got mempool packet: %v", mp)
		});

		searcher.SubscribeByAddress(context.Background(), []string{"addres1", "address2"}, func(mp *types.MempoolPacket) {
			t.Logf("got mempool packet: %v", mp)
		});
	})
}