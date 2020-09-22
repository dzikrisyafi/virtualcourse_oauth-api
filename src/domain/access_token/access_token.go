package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/dzikrisyafi/kursusvirtual_utils-go/crypto_utils"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/rest_errors"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant_type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() rest_errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return rest_errors.NewBadRequestError("invalid grant_type parameter")
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
	ClientID    string `json:"client_id"`
}

func (at *AccessToken) Validate() rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.NewBadRequestError("invalid access token id")
	}

	if at.UserID <= 0 {
		return rest_errors.NewBadRequestError("invalid user id")
	}

	if at.ClientID == "" {
		return rest_errors.NewBadRequestError("invalid client id")
	}

	if at.TokenType == "" {
		return rest_errors.NewBadRequestError("invalid token type")
	}

	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError("invalid expiration time")
	}

	return nil
}

func GetNewAccessToken(userID int, clientID string) AccessToken {
	return AccessToken{
		UserID:    userID,
		ClientID:  crypto_utils.GetSha256(clientID),
		TokenType: "Bearer",
		Expires:   int(time.Now().UTC().Add(expirationTime * time.Hour).Unix()),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(int64(at.Expires), 0).Before(time.Now().UTC())
}

func (at *AccessToken) GenerateByUserID() {
	at.AccessToken = crypto_utils.GetSha256(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}

func (at *AccessToken) GenerateByClientID() {
	at.AccessToken = crypto_utils.GetSha256(fmt.Sprintf("at-%s-%d-ran", at.ClientID, at.Expires))
}
