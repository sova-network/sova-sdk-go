package mevton_sdk_go

import (
	"context"
	"crypto/ed25519"
	"fmt"

	authpb "github.com/mevton-labs/mevton-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient wraps the AuthServiceClient and provides convenience methods.
type AuthClient struct {
	client       authpb.AuthServiceClient
	key          ed25519.PrivateKey
	refreshToken *authpb.Token
	accessToken  *authpb.Token
}

// NewAuthClient creates a new AuthClient instance.

// TODO use conn as a parameter; conn *grpc.ClientConn
func NewAuthClient(authUrl string, privateKey []byte) (*AuthClient, error) {
	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to block engine: %v", err)
	}
	client := authpb.NewAuthServiceClient(conn)
	key := ed25519.PrivateKey(privateKey)
	return &AuthClient{
		client:       client,
		key:          key,
		accessToken:  nil,
		refreshToken: nil,
	}, nil
}

func (a *AuthClient) Authenticate() error {
	ctx := context.Background()
	privKey := a.key
	pubkey := a.key.Public().(ed25519.PublicKey)
	challengeReq := &authpb.GenerateAuthChallengeRequest{
		Pubkey: pubkey,
	}

	challengeResp, err := a.client.GenerateAuthChallenge(ctx, challengeReq)
	if err != nil {
		return fmt.Errorf("failed to generate auth challenge: %v", err)
	}

	challenge := challengeResp.GetChallenge()
	signedChallenge := ed25519.Sign(privKey, challenge)

	tokenReq := &authpb.GenerateAuthTokensRequest{
		Challenge:       challenge,
		SignedChallenge: signedChallenge,
	}

	tokenResp, err := a.client.GenerateAuthTokens(ctx, tokenReq)
	if err != nil {
		return fmt.Errorf("failed to generate auth tokens: %v", err)
	}

	a.accessToken = tokenResp.GetAccessToken()
	a.refreshToken = tokenResp.GetRefreshToken()

	return nil
}

func (a *AuthClient) RefreshAccessToken() error {
	ctx := context.Background()
	if a.refreshToken != nil {
		refreshReq := &authpb.RefreshAccessTokenRequest{
			RefreshToken: a.refreshToken.GetValue(),
		}
		tokenResp, err := a.client.RefreshAccessToken(ctx, refreshReq)
		if err != nil {
			return fmt.Errorf("failed to refresh access token: %v", err)
		}
		a.accessToken = tokenResp.GetAccessToken()
		return nil
	}

	return fmt.Errorf("refresh token is required")
}

func (a *AuthClient) AccessToken() *authpb.Token {
	return a.accessToken
}

func (a *AuthClient) RefreshToken() *authpb.Token {
	return a.refreshToken
}
