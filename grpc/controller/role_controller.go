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

type RoleController struct {
	user.UnimplementedRoleServiceServer
	keycloak.RoleService
	keycloak.CredentialService
}

func (r RoleController) AssignRoleToUser(ctx context.Context, in *user.UserRoleRequest) (*empty.Empty, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "userId cannot be nil or empty")
	}

	if in.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role cannot be nil")
	}

	if in.Role.Name == nil || in.Role.Name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "role name cannot be nil or empty")
	}

	if in.Role.Id == nil || in.Role.Id.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "role id cannot be nil or empty")
	}

	role := domain.Role{
		Name: in.Role.Name.Value,
		Id:   in.Role.Id.Value,
	}

	token, err := r.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = r.RoleService.AssignRoleToUser(ctx, in.UserId, []domain.Role{role}, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (r RoleController) GetUserRoles(ctx context.Context, in *wrappers.StringValue) (*user.RolesResponse, error) {
	if in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "userId cannot be nil or empty")
	}

	token, err := r.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	roles, err := r.RoleService.GetUserRoles(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response []*user.RoleResponse
	for _, role := range roles {
		r := role.RoleToGRpcResponse()
		response = append(response, &r)
	}

	return &user.RolesResponse{Roles: response}, nil
}

func (r RoleController) GetAvailableRoles(ctx context.Context, in *wrappers.StringValue) (*user.RolesResponse, error) {
	if in == nil || in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "userId cannot be nil or empty")
	}

	token, err := r.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	roles, err := r.RoleService.GetAvailableRoles(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response []*user.RoleResponse
	for _, role := range roles {
		r := role.RoleToGRpcResponse()
		response = append(response, &r)
	}

	return &user.RolesResponse{Roles: response}, nil
}

func (r RoleController) RemoveRoleFromUser(ctx context.Context, in *user.UserRoleRequest) (*empty.Empty, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "userId cannot be nil or empty")
	}

	if in.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role cannot be nil or empty")
	}

	roleRequest := in.Role
	if roleRequest.Id.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "role id cannot be nil or empty")
	}

	token, err := r.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	role := domain.Role{
		Id: roleRequest.Id.Value,
	}

	if roleRequest.Name != nil && roleRequest.Name.Value != "" {
		role.Name = roleRequest.Name.Value
	}

	err = r.RoleService.RemoveRoleFromUser(ctx, in.UserId, []domain.Role{role}, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (r RoleController) CreateRole(ctx context.Context, in *user.RoleRequest) (*empty.Empty, error) {
	if in.Name == nil {
		return nil, status.Error(codes.InvalidArgument, "role cannot be nil or empty")
	}

	if in.Name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "role cannot be nil or empty")
	}

	role := domain.Role{
		Name: in.Name.Value,
	}

	token, err := r.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = r.RoleService.CreateRole(ctx, role, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
