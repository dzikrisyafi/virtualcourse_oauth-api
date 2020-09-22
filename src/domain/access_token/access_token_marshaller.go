package access_token

import "encoding/json"

type PasswordResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	UserID      int    `json:"user_id"`
}

type ClientCredentialsResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
}

func (at *AccessToken) Marshall() interface{} {
	if at.UserID == 0 {
		return ClientCredentialsResponse{
			AccessToken: at.AccessToken,
			TokenType:   at.TokenType,
			Expires:     at.Expires,
		}
	}

	atJson, _ := json.Marshal(at)
	var atResponse PasswordResponse
	json.Unmarshal(atJson, &atResponse)
	return atResponse
}
