package domain

import (
	"fmt"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
)

type GroupOverview struct {
	Id        string        `json:"id,omitempty"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	SubGroups []interface{} `json:"subGroups"`
}

type Group struct {
	Id          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Attributes  map[string]string `json:"attributes"`
	RealmRoles  []string          `json:"realmRoles"`
	ClientRoles struct {
	} `json:"clientRoles"`
	SubGroups []interface{} `json:"subGroups"`
	Access    Access        `json:"access"`
}

func (o GroupOverview) GroupOverviewToGRpcResponse() user.GroupResponse {
	return user.GroupResponse{
		Id:   o.Id,
		Name: o.Name,
	}
}

func (g Group) GroupToGRpcResponse() user.GroupResponse {
	return user.GroupResponse{
		Id:   g.Id,
		Name: g.Name,
	}
}

func GroupGRpcRequestToRequest(request *user.GroupRequest) GroupOverview {
	name := request.Name.Value

	return GroupOverview{
		Name: name,
		Path: fmt.Sprintf("/%s", name),
	}
}
