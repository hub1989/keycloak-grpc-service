package domain

import (
	"fmt"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"strings"
)

type Client struct {
	Id                                 string                 `json:"id"`
	ClientId                           string                 `json:"clientId"`
	Name                               string                 `json:"name"`
	RootUrl                            string                 `json:"rootUrl"`
	BaseUrl                            string                 `json:"baseUrl"`
	SurrogateAuthRequired              bool                   `json:"surrogateAuthRequired"`
	Enabled                            bool                   `json:"enabled"`
	AlwaysDisplayInConsole             bool                   `json:"alwaysDisplayInConsole"`
	ClientAuthenticatorType            string                 `json:"clientAuthenticatorType"`
	RedirectUris                       []string               `json:"redirectUris"`
	WebOrigins                         []interface{}          `json:"webOrigins"`
	NotBefore                          int                    `json:"notBefore"`
	BearerOnly                         bool                   `json:"bearerOnly"`
	ConsentRequired                    bool                   `json:"consentRequired"`
	StandardFlowEnabled                bool                   `json:"standardFlowEnabled"`
	ImplicitFlowEnabled                bool                   `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled          bool                   `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled             bool                   `json:"serviceAccountsEnabled"`
	PublicClient                       bool                   `json:"publicClient"`
	FrontchannelLogout                 bool                   `json:"frontchannelLogout"`
	Protocol                           string                 `json:"protocol"`
	Attributes                         map[string]interface{} `json:"attributes"`
	AuthenticationFlowBindingOverrides struct {
	} `json:"authenticationFlowBindingOverrides"`
	FullScopeAllowed          bool     `json:"fullScopeAllowed"`
	NodeReRegistrationTimeout int      `json:"nodeReRegistrationTimeout"`
	DefaultClientScopes       []string `json:"defaultClientScopes"`
	OptionalClientScopes      []string `json:"optionalClientScopes"`
	Access                    Access   `json:"access"`
}

func (c Client) ClientToGRpcResponse() user.ClientResponse {

	attributes := make(map[string]string)
	for key, value := range c.Attributes {
		x := fmt.Sprintf("%s", value)
		x = strings.ReplaceAll(x, "[", "")
		x = strings.ReplaceAll(x, "]", "")

		attributes[key] = x
	}

	return user.ClientResponse{
		Id:         c.Id,
		ClientId:   c.ClientId,
		Name:       c.Name,
		RootUrl:    c.RootUrl,
		WebUrl:     c.BaseUrl,
		Enabled:    c.Enabled,
		Attributes: attributes,
	}
}
