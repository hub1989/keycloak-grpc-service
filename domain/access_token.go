package domain

import (
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
)

type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

func (r AccessTokenResponse) AccessTokenToGRpcResponse() user.AccessTokenResponse {
	return user.AccessTokenResponse{
		AccessToken:      r.AccessToken,
		ExpiresIn:        int32(r.ExpiresIn),
		RefreshExpiresIn: int32(r.RefreshExpiresIn),
		RefreshToken:     r.RefreshToken,
		TokenType:        r.TokenType,
		IdToken:          r.IdToken,
		NotBeforePolicy:  int32(r.NotBeforePolicy),
		SessionState:     r.SessionState,
		Scope:            r.Scope,
	}
}
