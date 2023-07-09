package keycloak

import (
	"fmt"
	"net/http"
	"os"
)

type ClientCredentials struct {
	ClientId     string
	ClientSecret string
}
type Configuration interface {
	GetUserEndpoint() string
	GetBaseUrl() string
	GetRealm() string
	GetClientCredentials() ClientCredentials
	GetGroupEndpoint() string
	GetClient() *http.Client
}

type DefaultKeycloakConfiguration struct {
	BaseURL string
	Realm   string
	*http.Client
}

func (d DefaultKeycloakConfiguration) GetUserEndpoint() string {
	return fmt.Sprintf("%s/admin/realms/%s/users", d.BaseURL, d.Realm)
}

func (d DefaultKeycloakConfiguration) GetBaseUrl() string {
	return d.BaseURL
}

func (d DefaultKeycloakConfiguration) GetRealm() string {
	return d.Realm
}

func (d DefaultKeycloakConfiguration) GetClientCredentials() ClientCredentials {
	return ClientCredentials{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func (d DefaultKeycloakConfiguration) GetGroupEndpoint() string {
	return fmt.Sprintf("%s/admin/realms/%s/groups", d.BaseURL, d.Realm)
}

func (d DefaultKeycloakConfiguration) GetClient() *http.Client {
	return d.Client
}
