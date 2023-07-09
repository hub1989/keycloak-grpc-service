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
)

type RoleService interface {
	AssignRoleToUser(ctx context.Context, userId string, role []domain.Role, token string) error
	GetUserRoles(ctx context.Context, userId string, token string) ([]domain.Role, error)
	GetAvailableRoles(ctx context.Context, userId string, token string) ([]domain.Role, error)
	RemoveRoleFromUser(ctx context.Context, userId string, role []domain.Role, token string) error
	CreateRole(ctx context.Context, role domain.Role, token string) error
}

type DefaultRoleService struct {
	Configuration
	ClientService
}

func (d DefaultRoleService) AssignRoleToUser(ctx context.Context, userId string, role []domain.Role, token string) error {
	endpoint := fmt.Sprintf("%s/%s/role-mappings/realm", d.GetUserEndpoint(), userId)

	body, err := json.Marshal(role)
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
		log.WithError(err).Error("could not assign role user")
	}

	if resp.StatusCode != 204 {
		return errors.New(fmt.Sprintf("could not assign role to user, see reason: %v", resp.Status))
	}

	return nil
}

func (d DefaultRoleService) GetUserRoles(ctx context.Context, userId string, token string) ([]domain.Role, error) {
	endpoint := fmt.Sprintf("%s/%s/role-mappings/realm", d.GetUserEndpoint(), userId)
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
		log.WithError(err).Error("could not available user roles")
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("could not get user roles -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var roles []domain.Role
	err = json.Unmarshal(data, &roles)

	return roles, nil
}

func (d DefaultRoleService) GetAvailableRoles(ctx context.Context, userId string, token string) ([]domain.Role, error) {
	endpoint := fmt.Sprintf("%s/%s/role-mappings/realm/available", d.GetUserEndpoint(), userId)
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
		log.WithError(err).Error("could not available user roles")
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("could not get user info -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var roles []domain.Role
	err = json.Unmarshal(data, &roles)

	return roles, nil
}

func (d DefaultRoleService) RemoveRoleFromUser(ctx context.Context, userId string, role []domain.Role, token string) error {
	endpoint := fmt.Sprintf("%s/%s/role-mappings/realm", d.GetUserEndpoint(), userId)

	body, err := json.Marshal(role)
	bodyReader := bytes.NewReader(body)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not delete user")
	}

	if res.StatusCode != 204 {
		return errors.New("could not get user info -- see reason " + res.Status)
	}

	return nil
}

func (d DefaultRoleService) CreateRole(ctx context.Context, role domain.Role, token string) error {
	clientName := d.Configuration.GetClientCredentials().ClientId
	client, err := d.ClientService.GetClientByClientId(ctx, clientName, token)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/admin/realms/%s/clients/%s/roles", d.Configuration.GetBaseUrl(), d.Configuration.GetRealm(), client.Id)
	body, err := json.Marshal(role)
	bodyReader := bytes.NewReader(body)

	httpClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		log.WithError(err).Error("could not assign role user")
	}

	if resp.StatusCode != 201 {
		return errors.New(fmt.Sprintf("could not create role, see reason: %v", resp.Status))
	}

	return nil
}
