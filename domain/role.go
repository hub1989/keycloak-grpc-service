package domain

import user "github.com/hub1989/keycloak-protobuf/golang/keycloak"

type Role struct {
	Id                 string      `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	ScopeParamRequired interface{} `json:"scopeParamRequired"`
	Composite          bool        `json:"composite"`
	Composites         interface{} `json:"composites"`
	ClientRole         bool        `json:"clientRole"`
	ContainerId        string      `json:"containerId"`
	Attributes         interface{} `json:"attributes"`
}

func (r Role) RoleToGRpcResponse() user.RoleResponse {
	return user.RoleResponse{
		Id:   r.Id,
		Name: r.Name,
	}
}
