package sova_sdk_go

import (
	"context"
	"errors"

	types "github.com/sova-network/sova-sdk-go/generated"
)

type SovaClient struct {
	url        string
	caPem      *string
	domainName *string
	authToken  *types.Token
}

func (s *SovaClient) NewMainnetClient() (*SovaClient, error) {
	return nil, errors.New("not released")
}

func NewTestnetClient() *SovaClient {
	return NewCustomClient(
		"https://testnet-engine.sova.finance:30010",
		"TESTNET_CA_PEM", // Replace with actual PEM
		"testnet-engine.sova.finance",
		nil,
	)
}

// NewCustomClient creates a new custom client
func NewCustomClient(url, caPem, domainName string, authToken *types.Token) *SovaClient {
	return &SovaClient{
		url:        url,
		caPem:      &caPem,
		domainName: &domainName,
		authToken:  authToken,
	}
}

// Authenticate performs authentication
func (c *SovaClient) Authenticate(ctx context.Context, privateKey []byte) (*types.Token, error) {
	auth, err := NewAuthClient(c.url, privateKey, c.caPem, c.domainName)
	if err != nil {
		return nil, err
	}

	err = auth.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	token := auth.AccessToken()

	c.authToken = token
	return token, nil
}

// Searcher returns a new searcher client
func (c *SovaClient) Searcher(ctx context.Context) (*SovaSearcher, error) {
	return NewSovaSearcherWithAccessToken(c.url, c.caPem, c.domainName, c.authToken)
}
