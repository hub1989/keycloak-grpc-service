package domain

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"strings"
)

type UserRepresentation struct {
	Id                         string                 `json:"id,omitempty"`
	CreatedTimestamp           int64                  `json:"createdTimestamp,omitempty"`
	Username                   string                 `json:"username,omitempty"`
	Enabled                    bool                   `json:"enabled"`
	Totp                       bool                   `json:"totp,omitempty"`
	EmailVerified              bool                   `json:"emailVerified"`
	FirstName                  string                 `json:"firstName,omitempty"`
	LastName                   string                 `json:"lastName,omitempty"`
	Email                      string                 `json:"email,omitempty"`
	DisableableCredentialTypes []interface{}          `json:"disableableCredentialTypes,omitempty"`
	RequiredActions            []interface{}          `json:"requiredActions,omitempty"`
	NotBefore                  int                    `json:"notBefore,omitempty"`
	Credentials                []Credential           `json:"credentials,omitempty"`
	Access                     Access                 `json:"access,omitempty"`
	Attributes                 map[string]interface{} `json:"attributes,omitempty"`
	RealmRoles                 []string               `json:"realmRoles,omitempty"`
}

type Credential struct {
	Value     string `json:"value"`
	Type      string `json:"type"`
	Temporary bool   `json:"temporary"`
}

type Access struct {
	ManageGroupMembership bool `json:"manageGroupMembership"`
	View                  bool `json:"view"`
	MapRoles              bool `json:"mapRoles"`
	Impersonate           bool `json:"impersonate"`
	Manage                bool `json:"manage"`
}

type OIDCInfo struct {
	Subject             string   `json:"subject"`
	Iss                 string   `json:"iss"`
	Aud                 []string `json:"aud"`
	Sub                 string   `json:"sub"`
	Name                string   `json:"name"`
	GivenName           string   `json:"given_name"`
	FamilyName          string   `json:"family_name"`
	MiddleName          string   `json:"middle_name"`
	Nickname            string   `json:"nickname"`
	PreferredUsername   string   `json:"preferred_username"`
	Profile             string   `json:"profile"`
	Picture             string   `json:"picture"`
	Website             string   `json:"website"`
	Email               string   `json:"email"`
	EmailVerified       bool     `json:"email_verified"`
	Gender              string   `json:"gender"`
	Birthdate           string   `json:"birthdate"`
	Zoneinfo            string   `json:"zoneinfo"`
	Locale              string   `json:"locale"`
	PhoneNumber         string   `json:"phone_number"`
	PhoneNumberVerified bool     `json:"phone_number_verified"`
	UpdatedAt           int      `json:"updated_at"`
	ClaimsLocales       string   `json:"claims_locales"`
}

func (r UserRepresentation) UserToGRpcResponse() user.UserResponse {
	attributes := make(map[string]string)
	for key, value := range r.Attributes {
		x := fmt.Sprintf("%s", value)
		x = strings.ReplaceAll(x, "[", "")
		x = strings.ReplaceAll(x, "]", "")

		attributes[key] = x
	}

	phoneNumber := attributes["phoneNumber"]
	return user.UserResponse{
		Email:       r.Email,
		PhoneNumber: &wrappers.StringValue{Value: phoneNumber},
		GivenName:   &wrappers.StringValue{Value: r.FirstName},
		FamilyName:  &wrappers.StringValue{Value: r.LastName},
		Sub:         r.Id,
		Attributes:  attributes,
		Username:    r.Username,
		Roles:       nil,
	}
}

func UserGRpcRequestToUser(request *user.UserRequest) UserRepresentation {
	attributes := make(map[string]interface{})
	for key, value := range request.Attributes {
		attributes[key] = value
	}

	userRepresentation := UserRepresentation{
		Username:      request.Username,
		Enabled:       request.Enabled,
		EmailVerified: request.EmailVerified,
		Email:         request.Email.Value,
		Attributes:    attributes,
	}

	if request.PhoneNumber != nil && request.PhoneNumber.Value != "" {
		userRepresentation.Attributes["phoneNumber"] = request.PhoneNumber.Value
	}

	if request.FirstName != nil && request.FirstName.Value != "" {
		userRepresentation.FirstName = request.FirstName.Value
	}

	if request.LastName != nil && request.LastName.Value != "" {
		userRepresentation.LastName = request.LastName.Value
	}

	if request.Password != "" {
		credentials := []Credential{
			{
				Value:     request.Password,
				Type:      "password",
				Temporary: false,
			},
		}

		userRepresentation.Credentials = credentials
	}
	return userRepresentation
}

func (r UserRepresentation) UpdateUser(request *user.UpdateUserRequest) UserRepresentation {
	if request.PhoneNumber != nil && request.PhoneNumber.Value != "" {
		r.Attributes["phoneNumber"] = request.PhoneNumber.Value
	}

	if request.FirstName != nil && request.FirstName.Value != "" {
		r.FirstName = request.FirstName.Value
	}

	if request.LastName != nil && request.LastName.Value != "" {
		r.LastName = request.LastName.Value
	}

	if r.Enabled != request.Enabled {
		r.Enabled = request.Enabled
	}

	if r.EmailVerified != request.EmailVerified {
		r.EmailVerified = request.EmailVerified
	}

	if request.Attributes != nil {
		for key, value := range request.Attributes {
			r.Attributes[key] = value
		}
	}

	return r
}
