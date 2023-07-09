package controller

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hub1989/keycloak-grpc-service/domain"
	"github.com/hub1989/keycloak-grpc-service/keycloak"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupController struct {
	user.UnimplementedGroupServiceServer
	keycloak.CredentialService
	keycloak.GroupService
}

func (g GroupController) CreateGroup(ctx context.Context, in *user.GroupRequest) (*user.GroupResponse, error) {
	if in.Name == nil || in.Name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "group name cannot be nil or empty")
	}

	request := domain.GroupGRpcRequestToRequest(in)

	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = g.GroupService.CreateGroup(ctx, request, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.GroupResponse{
		Name: request.Name,
	}, err
}

func (g GroupController) GetGroupsInRealm(ctx context.Context, in *wrappers.StringValue) (*user.GroupsResponse, error) {
	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := g.GroupService.GetGroupsInRealm(ctx, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var groups []*user.GroupResponse

	for _, group := range resp {
		g := group.GroupOverviewToGRpcResponse()
		groups = append(groups, &g)
	}

	return &user.GroupsResponse{Groups: groups}, nil
}

func (g GroupController) GetGroupById(ctx context.Context, in *wrappers.StringValue) (*user.GroupResponse, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "group name cannot be nil or empty")
	}

	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := g.GroupService.GetGroupById(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	group := resp.GroupToGRpcResponse()
	return &group, nil
}

func (g GroupController) DeleteGroup(ctx context.Context, in *wrappers.StringValue) (*empty.Empty, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "group name cannot be nil or empty")
	}

	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = g.GroupService.DeleteGroup(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GroupController) GetGroupMembers(ctx context.Context, in *wrappers.StringValue) (*user.UsersResponse, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "group name cannot be nil or empty")
	}

	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := g.GroupService.GetGroupMembers(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var users []*user.UserResponse

	for _, user := range resp {
		u := user.UserToGRpcResponse()
		users = append(users, &u)
	}

	return &user.UsersResponse{Users: users}, nil
}

func (g GroupController) AddRoleToGroup(ctx context.Context, in *user.RoleGroupRequest) (*empty.Empty, error) {

	if in.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group id cannot be nil or empty")
	}

	if in.Role.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "role id cannot be nil or empty")
	}

	if in.Role.Name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "role name cannot be nil or empty")
	}

	token, err := g.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = g.GroupService.AddRoleToGroup(ctx, domain.Role{Name: in.Role.Name.Value}, in.GroupId, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
