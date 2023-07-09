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

type GroupService interface {
	CreateGroup(ctx context.Context, request domain.GroupOverview, token string) error
	GetGroupsInRealm(ctx context.Context, token string) ([]domain.GroupOverview, error)
	GetGroupById(ctx context.Context, groupId, token string) (domain.Group, error)
	DeleteGroup(ctx context.Context, groupId, token string) error
	GetGroupMembers(ctx context.Context, groupId, token string) ([]domain.UserRepresentation, error)
	AddRoleToGroup(ctx context.Context, role domain.Role, groupId string, token string) error
}

type DefaultGroupService struct {
	Configuration
}

func (d DefaultGroupService) CreateGroup(ctx context.Context, request domain.GroupOverview, token string) error {
	endpoint := d.Configuration.GetGroupEndpoint()

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

	if resp.StatusCode != 201 {
		return errors.New(fmt.Sprintf("could not create group, see reason: %v", resp.Status))
	}

	if err != nil {
		log.WithError(err).Error("could not create group")
	}

	return nil
}

func (d DefaultGroupService) GetGroupsInRealm(ctx context.Context, token string) ([]domain.GroupOverview, error) {
	endpoint := d.Configuration.GetGroupEndpoint()
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
		log.WithError(err).Error("could not get realm groups")
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

	var groups []domain.GroupOverview
	err = json.Unmarshal(data, &groups)

	return groups, nil
}

func (d DefaultGroupService) GetGroupById(ctx context.Context, groupId, token string) (domain.Group, error) {
	endpoint := fmt.Sprintf("%s/%s", d.GetGroupEndpoint(), groupId)
	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return domain.Group{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not get realm group")
		return domain.Group{}, err
	}

	if res.StatusCode != 200 {
		return domain.Group{}, errors.New("could not get user info -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return domain.Group{}, err
	}

	var group domain.Group
	err = json.Unmarshal(data, &group)

	return group, nil
}

func (d DefaultGroupService) DeleteGroup(ctx context.Context, groupId, token string) error {
	endpoint := fmt.Sprintf("%s/%s", d.GetGroupEndpoint(), groupId)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not delete group")
	}

	if res.StatusCode != 204 {
		return errors.New("could not delete group -- see reason " + res.Status)
	}

	return nil
}

func (d DefaultGroupService) GetGroupMembers(ctx context.Context, groupId, token string) ([]domain.UserRepresentation, error) {
	endpoint := fmt.Sprintf("%s/%s/members", d.GetGroupEndpoint(), groupId)
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
		log.WithError(err).Error("could not get group members")
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("could not get group members -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var users []domain.UserRepresentation
	err = json.Unmarshal(data, &users)

	return users, nil
}

func (d DefaultGroupService) AddRoleToGroup(ctx context.Context, role domain.Role, groupId string, token string) error {
	endpoint := fmt.Sprintf("%s/%s/role-mappings/realm", d.GetGroupEndpoint(), groupId)

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

	res, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("could not add role to group")
	}

	if res.StatusCode != 201 {
		return errors.New("could not add role to group -- see reason " + res.Status)
	}

	return nil
}
