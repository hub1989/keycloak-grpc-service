package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hub1989/keycloak-grpc-service/keycloak"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientController struct {
	user.UnimplementedClientServiceServer
	keycloak.ClientService
	keycloak.CredentialService
}

func (c ClientController) GetClients(ctx context.Context, in *empty.Empty) (*user.ClientsResponse, error) {
	token, err := c.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	clients, err := c.ClientService.GetClients(ctx, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response []*user.ClientResponse

	for _, client := range clients {
		gRpcResponse := client.ClientToGRpcResponse()
		response = append(response, &gRpcResponse)
	}

	return &user.ClientsResponse{Clients: response}, nil
}

func (c ClientController) GetClientByClientId(ctx context.Context, in *wrappers.StringValue) (*user.ClientResponse, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "clientId cannot be nil or empty")
	}
	token, err := c.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	client, err := c.ClientService.GetClientByClientId(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := client.ClientToGRpcResponse()
	return &resp, nil
}

func (c ClientController) GetClientById(ctx context.Context, in *wrappers.StringValue) (*user.ClientResponse, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be nil or empty")
	}
	token, err := c.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	client, err := c.ClientService.GetClientById(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := client.ClientToGRpcResponse()
	return &resp, nil
}

func (c ClientController) GetClientsByIds(ctx context.Context, in *user.StringsRequest) (*user.ClientsResponse, error) {
	if in == nil || len(in.Requests) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ids cannot be nil or empty")
	}
	var ids []string

	for _, request := range in.Requests {
		ids = append(ids, request)
	}
	token, err := c.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	clients, err := c.ClientService.GetClientsByIds(ctx, ids, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response []*user.ClientResponse

	for _, client := range clients {
		gRpcResponse := client.ClientToGRpcResponse()
		response = append(response, &gRpcResponse)
	}

	return &user.ClientsResponse{Clients: response}, nil
}

func (c ClientController) GetClientsByClientIds(ctx context.Context, in *user.StringsRequest) (*user.ClientsResponse, error) {
	if in == nil || len(in.Requests) == 0 {
		return nil, status.Error(codes.InvalidArgument, "clientIds cannot be nil or empty")
	}
	var ids []string

	for _, request := range in.Requests {
		ids = append(ids, request)
	}
	token, err := c.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	clients, err := c.ClientService.GetClientsByClientIds(ctx, ids, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response []*user.ClientResponse

	for _, client := range clients {
		gRpcResponse := client.ClientToGRpcResponse()
		response = append(response, &gRpcResponse)
	}

	return &user.ClientsResponse{Clients: response}, nil
}
