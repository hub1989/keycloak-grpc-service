package controller

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/hub1989/keycloak-grpc-service/domain"
	"github.com/hub1989/keycloak-grpc-service/keycloak"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	user.UnimplementedUserServiceServer
	keycloak.CredentialService
	keycloak.UserService
}

func (u UserController) CreateUser(ctx context.Context, in *user.UserRequest) (*empty.Empty, error) {
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be nil or empty")
	}

	request := domain.UserGRpcRequestToUser(in)

	fmt.Println(request)
	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = u.UserService.CreateUser(ctx, request, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (u UserController) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*empty.Empty, error) {
	if in.Pid == "" {
		return nil, status.Error(codes.InvalidArgument, "user pid cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	user, err := u.UserService.GetUserById(ctx, in.Pid, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	request := user.UpdateUser(in)

	err = u.UserService.UpdateUser(ctx, request, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (u UserController) GetUserById(ctx context.Context, in *wrappers.StringValue) (*user.UserResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "id cannot be nil or empty")
	}

	if in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := u.UserService.GetUserById(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	grpcResp := resp.UserToGRpcResponse()

	return &grpcResp, nil
}

func (u UserController) GetUserByUsername(ctx context.Context, in *wrappers.StringValue) (*user.UserResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "username cannot be nil or empty")
	}

	if in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := u.UserService.GetUserByUsername(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	grpcResp := resp.UserToGRpcResponse()

	return &grpcResp, nil
}

func (u UserController) DeleteUser(ctx context.Context, in *wrappers.StringValue) (*empty.Empty, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "id cannot be nil or empty")
	}

	if in.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = u.UserService.DeleteUser(ctx, in.Value, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (u UserController) AddUserToGroup(ctx context.Context, in *user.UserGroupRequest) (*empty.Empty, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil or empty")
	}

	if in.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "groupId cannot be nil or empty")
	}

	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = u.UserService.AddUserToGroup(ctx, in.UserId, in.GroupId, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (u UserController) RemoveUserFromGroup(ctx context.Context, in *user.UserGroupRequest) (*empty.Empty, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "user id cannot be nil or empty")
	}

	if in.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "groupId cannot be nil or empty")
	}

	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "userId cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = u.UserService.RemoveUserFromGroup(ctx, in.UserId, in.GroupId, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (u UserController) Authenticate(ctx context.Context, in *user.AuthenticateRequest) (*user.AccessTokenResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil or empty")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password cannot be nil or empty")
	}

	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "usernamecannot be nil or empty")
	}

	if in.ClientId == nil || in.ClientId.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "client id cannot be nil or empty")
	}

	if in.ClientSecret != nil && in.ClientSecret.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "client secret cannot be nil or empty")
	}

	clientSecret := ""

	if in.ClientSecret != nil && in.ClientSecret.Value != "" {
		clientSecret = in.ClientSecret.Value
	}

	resp, err := u.UserService.Authenticate(ctx, in.Username, in.Password, in.ClientId.Value, clientSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	gRpc := resp.AccessTokenToGRpcResponse()
	return &gRpc, err
}

func (u UserController) GetAllUsers(ctx context.Context, in *empty.Empty) (*user.UsersResponse, error) {
	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	users, err := u.UserService.GetAllUsers(ctx, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp []*user.UserResponse
	for _, representation := range users {
		gRpcResponse := representation.UserToGRpcResponse()
		resp = append(resp, &gRpcResponse)
	}

	return &user.UsersResponse{
		Users: resp,
	}, nil
}

func (u UserController) GetUsersByIds(ctx context.Context, in *user.StringsRequest) (*user.UsersResponse, error) {
	if in == nil || len(in.Requests) == 0 {
		return nil, status.Error(codes.InvalidArgument, "userIds cannot be nil or empty")
	}
	var ids []string

	for _, request := range in.Requests {
		ids = append(ids, request)
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	users, err := u.UserService.GetUsersByIds(ctx, ids, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp []*user.UserResponse
	for _, representation := range users {
		gRpcResponse := representation.UserToGRpcResponse()
		resp = append(resp, &gRpcResponse)
	}

	return &user.UsersResponse{
		Users: resp,
	}, nil
}

func (u UserController) GetUsersByUsernames(ctx context.Context, in *user.StringsRequest) (*user.UsersResponse, error) {
	if in == nil || len(in.Requests) == 0 {
		return nil, status.Error(codes.InvalidArgument, "usernames cannot be nil or empty")
	}
	var usernames []string

	for _, request := range in.Requests {
		usernames = append(usernames, request)
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	users, err := u.UserService.GetUsersByUsernames(ctx, usernames, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp []*user.UserResponse
	for _, representation := range users {
		gRpcResponse := representation.UserToGRpcResponse()
		resp = append(resp, &gRpcResponse)
	}

	return &user.UsersResponse{
		Users: resp,
	}, nil
}

func (u UserController) SetUserPassword(ctx context.Context, in *user.PasswordRequest) (*wrappers.BoolValue, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil or empty")
	}

	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id cannot be nil or empty")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password cannot be nil or empty")
	}

	token, err := u.CredentialService.ObtainTokenForOps(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := u.UserService.SetUserPassword(ctx, in.Password, in.UserId, in.Temporary, token.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &wrappers.BoolValue{Value: resp}, nil
}
