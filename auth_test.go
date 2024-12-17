package mevton_sdk_go

import (
	"context"
	"net"
	"testing"

	types "github.com/mevton-labs/mevton-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MockAuthServiceClient struct {
	types.AuthServiceServer
}


func (receiver MockAuthServiceClient) GenerateAuthChallenge(context.Context, *types.GenerateAuthChallengeRequest) (*types.GenerateAuthChallengeResponse, error) {
	return &types.GenerateAuthChallengeResponse{
		Challenge: []byte("test_challenge"),
	}, nil
}

func (receiver MockAuthServiceClient) GenerateAuthTokens(context.Context, *types.GenerateAuthTokensRequest) (*types.GenerateAuthTokensResponse, error) {
	return &types.GenerateAuthTokensResponse{
		AccessToken:  &types.Token{Value: "access_token"},
		RefreshToken: &types.Token{Value: "refresh_token"},
	}, nil
}

func (receiver MockAuthServiceClient) RefreshAccessToken(context.Context, *types.RefreshAccessTokenRequest) (*types.RefreshAccessTokenResponse, error) {
	return &types.RefreshAccessTokenResponse{AccessToken: &types.Token{Value: "new_access_token"}}, nil
}

func RunServer() {
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	types.RegisterAuthServiceServer(s, &MockAuthServiceClient{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func TestAuthClient_Authenticate(t *testing.T) {
	RunServer()
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	privateKey := []byte{
		155, 202, 118, 43, 82, 100, 113, 150, 99, 21,
		45, 230, 88, 247, 193, 12, 92, 78, 191, 229,
		73, 191, 100, 156, 231, 41, 144, 54, 202, 199,
		75, 1,
	}

	client, err := NewAuthClient("127.0.0.1:50051", privateKey, nil, nil)
	if err != nil {
		t.Fatalf("Cannot create Auth Client: %v", err)
	}

	err = client.Authenticate()
	if err != nil {
		t.Fatalf("Cannot Authenticate: %v", err)
	}

	accessToken := client.AccessToken()
	expectedAccessToken := "access_token"
	if accessToken.Value != expectedAccessToken {
		t.Errorf("Expected %q, got %q", expectedAccessToken, accessToken.Value)
	}

	refreshToken := client.RefreshToken()
	expectedRefreshToken := "refresh_token"
	if refreshToken.Value != expectedRefreshToken {
		t.Errorf("Expected %q, got %q", expectedRefreshToken, refreshToken.Value)
	}

	err = client.RefreshAccessToken()
	if err != nil {
		t.Fatalf("Cannot Refresh Access Token: %v", err)
	}

	accessToken = client.AccessToken()
	expectedNewAccessToken := "new_access_token"
	if accessToken.Value != expectedNewAccessToken {
		t.Errorf("Expected %q, got %q", expectedNewAccessToken, accessToken.Value)
	}
}
