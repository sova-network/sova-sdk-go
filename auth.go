package mevton_sdk_go

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/mevton-labs/mevton-sdk-go/generated"
	"google.golang.org/grpc"
)

// AuthClient wraps the AuthServiceClient and provides convenience methods.
type AuthClient struct {
	client       authpb.AuthServiceClient
	key          ed25519.PrivateKey
	refreshToken *authpb.Token
	accessToken  *authpb.Token
}

// NewAuthClient creates a new AuthClient instance.
func NewAuthClient(authUrl string, privateKey []byte, caPem *[]byte, domainName *string) (*AuthClient, error) {
	var conn *grpc.ClientConn
	var err error

	if caPem != nil && domainName != nil {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(*caPem)) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		// Configure TLS with CA certificate and domain name
		creds := credentials.NewClientTLSFromCert(certPool, *domainName)
		conn, err = grpc.NewClient(authUrl, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("failed to connect with TLS: %v", err)
		}
	} else {
		conn, err = grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to searcher service: %v", err)
		}
	}

	client := authpb.NewAuthServiceClient(conn)
	privKey := ed25519.NewKeyFromSeed(privateKey)
	return &AuthClient{
		client:       client,
		key:          privKey,
		accessToken:  nil,
		refreshToken: nil,
	}, nil
}

func (a *AuthClient) Authenticate() error {
	ctx := context.Background()
	pubkey := a.key.Public().(ed25519.PublicKey)
	challengeReq := &authpb.GenerateAuthChallengeRequest{
		Pubkey: pubkey,
	}

	challengeResp, err := a.client.GenerateAuthChallenge(ctx, challengeReq)
	if err != nil {
		return fmt.Errorf("failed to generate auth challenge: %v", err)
	}

	challenge := challengeResp.GetChallenge()

	signedChallenge := ed25519.Sign(a.key[:], challenge)
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
