package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hub1989/keycloak-grpc-service/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type UserService interface {
	CreateUser(ctx context.Context, request domain.UserRepresentation, token string) error
	UpdateUser(ctx context.Context, request domain.UserRepresentation, token string) error
	GetUserById(ctx context.Context, id string, token string) (domain.UserRepresentation, error)
	GetUserByUsername(ctx context.Context, username string, token string) (domain.UserRepresentation, error)
	DeleteUser(ctx context.Context, id string, token string) error
	AddUserToGroup(ctx context.Context, userId, groupId, token string) error
	RemoveUserFromGroup(ctx context.Context, userId, groupId, token string) error
	Authenticate(ctx context.Context, username, password string, clientId, clientSecret string) (domain.AccessTokenResponse, error)

	GetAllUsers(ctx context.Context, token string) ([]domain.UserRepresentation, error)
	GetUsersByIds(ctx context.Context, ids []string, token string) ([]domain.UserRepresentation, error)
	GetUsersByUsernames(ctx context.Context, usernames []string, token string) ([]domain.UserRepresentation, error)

	SetUserPassword(ctx context.Context, password string, id string, temporary bool, token string) (bool, error)
}

type DefaultUserService struct {
	Configuration
}

func (d DefaultUserService) CreateUser(ctx context.Context, request domain.UserRepresentation, token string) error {
	endpoint := d.GetUserEndpoint()

	body, err := json.Marshal(request)
	bodyReader := bytes.NewReader(body)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not create user")
		return err
	}

	if resp.StatusCode != 201 {
		log.Error("could not create user: ", resp)
		return errors.New(fmt.Sprintf("could not create user. got status: %v", resp.StatusCode))
	}

	return nil
}

func (d DefaultUserService) UpdateUser(ctx context.Context, request domain.UserRepresentation, token string) error {

	endpoint := fmt.Sprintf("%s/%s", d.GetUserEndpoint(), request.Id)

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not update user")
		return err
	}

	if resp.StatusCode != 204 {
		log.Error("could not update user: ", resp)
		return errors.New(fmt.Sprintf("could not update user. got status: %v", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (d DefaultUserService) GetUserById(ctx context.Context, id string, token string) (domain.UserRepresentation, error) {
	endpoint := fmt.Sprintf("%s/%s", d.GetUserEndpoint(), id)
	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return domain.UserRepresentation{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not get user")
		return domain.UserRepresentation{}, err
	}

	if res.StatusCode != 200 {
		return domain.UserRepresentation{}, errors.New(fmt.Sprintf("could not get user. got status: %v", res.StatusCode))
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return domain.UserRepresentation{}, err
	}

	var userInfo domain.UserRepresentation
	err = json.Unmarshal(data, &userInfo)

	if err != nil {
		return domain.UserRepresentation{}, err
	}

	return userInfo, nil
}

func (d DefaultUserService) DeleteUser(ctx context.Context, id string, token string) error {
	endpoint := fmt.Sprintf("%s/%s", d.GetUserEndpoint(), id)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not delete user")
	}

	if resp.StatusCode != 204 {
		log.Error("could not delete user: ", resp)
		return errors.New(fmt.Sprintf("could not delete user. got status: %v", resp.StatusCode))
	}

	return nil
}

func (d DefaultUserService) Authenticate(ctx context.Context, username, password string, clientId, clientSecret string) (domain.AccessTokenResponse, error) {

	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", d.GetBaseUrl(), d.GetRealm())
	method := "POST"

	password = url.QueryEscape(password)

	requestForm := fmt.Sprintf("client_id=%s&username=%s&password=%s&grant_type=password&scope=openid profile", clientId, username, password)

	if clientSecret != "" {
		requestForm = fmt.Sprintf("%s&client_secret=%s", requestForm, clientSecret)
	}
	payload := strings.NewReader(requestForm)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, method, endpoint, payload)

	if err != nil {
		log.WithError(err).Error("could not get access token for operation")
		return domain.AccessTokenResponse{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("could not get access token for operation")
		return domain.AccessTokenResponse{}, err
	}

	if res.StatusCode != 200 {
		return domain.AccessTokenResponse{}, errors.New("could not authenticate user :" + res.Status)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("could not get access token for operation")
		return domain.AccessTokenResponse{}, err
	}

	var accessToken domain.AccessTokenResponse
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return domain.AccessTokenResponse{}, err
	}

	return accessToken, nil
}

func (d DefaultUserService) GetUserByUsername(ctx context.Context, username string, token string) (domain.UserRepresentation, error) {
	endpoint := fmt.Sprintf("%s?username=%s", d.GetUserEndpoint(), username)
	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return domain.UserRepresentation{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not get user")
		return domain.UserRepresentation{}, err
	}

	if res.StatusCode != 200 {
		log.WithError(err).Error("could not get user")
		return domain.UserRepresentation{}, errors.New("could not get user info -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return domain.UserRepresentation{}, err
	}

	var users []domain.UserRepresentation
	err = json.Unmarshal(data, &users)
	if err != nil {
		return domain.UserRepresentation{}, err
	}

	if len(users) > 0 {
		return users[0], nil
	}

	return domain.UserRepresentation{}, errors.New("no user found")
}

func (d DefaultUserService) AddUserToGroup(ctx context.Context, userId, groupId, token string) error {
	endpoint := fmt.Sprintf("%s/%s/groups/%s", d.GetUserEndpoint(), userId, groupId)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not add user to group")
	}

	if resp.StatusCode != 201 {
		log.Error("could not add user to group: ", resp)
		return errors.New(fmt.Sprintf("could not create user. got status: %v", resp.StatusCode))
	}

	return nil
}

func (d DefaultUserService) RemoveUserFromGroup(ctx context.Context, userId, groupId, token string) error {
	endpoint := fmt.Sprintf("%s/%s/groups/%s", d.GetUserEndpoint(), userId, groupId)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not remove user from group")
	}

	if resp.StatusCode != 201 {
		log.Error("could not remove user from group: ", resp)
		return errors.New(fmt.Sprintf("could not create user. got status: %v", resp.StatusCode))
	}

	return nil
}

func (d DefaultUserService) GetAllUsers(ctx context.Context, token string) ([]domain.UserRepresentation, error) {
	endpoint := d.GetUserEndpoint()
	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not get user")
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("could not get users. got status: %v", res.StatusCode))
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var userInfo []domain.UserRepresentation
	err = json.Unmarshal(data, &userInfo)

	return userInfo, nil
}

func (d DefaultUserService) GetUsersByIds(ctx context.Context, ids []string, token string) ([]domain.UserRepresentation, error) {
	users, err := d.GetAllUsers(ctx, token)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]domain.UserRepresentation)

	for _, c := range users {
		userMap[c.Id] = c
	}

	var response []domain.UserRepresentation
	for _, id := range ids {
		user, ok := userMap[id]
		if ok {
			response = append(response, user)
		}
	}

	return response, nil
}

func (d DefaultUserService) GetUsersByUsernames(ctx context.Context, usernames []string, token string) ([]domain.UserRepresentation, error) {
	users, err := d.GetAllUsers(ctx, token)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]domain.UserRepresentation)

	for _, c := range users {
		userMap[c.Username] = c
	}

	var response []domain.UserRepresentation
	for _, username := range usernames {
		response = append(response, userMap[username])
	}

	return response, nil
}

func (d DefaultUserService) SetUserPassword(ctx context.Context, password string, id string, temporary bool, token string) (bool, error) {
	endpoint := fmt.Sprintf("%s/%s/reset-password", d.GetUserEndpoint(), id)

	credentials := domain.Credential{
		Value:     password,
		Type:      "password",
		Temporary: false,
	}

	body, err := json.Marshal(credentials)
	bodyReader := bytes.NewReader(body)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bodyReader)

	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not update user password")
	}

	if resp.StatusCode != 204 {
		log.Error("could not update user password: ", resp)
		return false, errors.New(fmt.Sprintf("could not setuser password user. got status: %v", resp.StatusCode))
	}

	return true, nil
}
