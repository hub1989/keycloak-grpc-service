package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hub1989/keycloak-grpc-service/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type ClientService interface {
	GetClients(ctx context.Context, token string) ([]domain.Client, error)
	GetClientById(ctx context.Context, clientId, token string) (domain.Client, error)
	GetClientByClientId(ctx context.Context, clientName, token string) (domain.Client, error)
	GetClientsByIds(ctx context.Context, ids []string, token string) ([]domain.Client, error)
	GetClientsByClientIds(ctx context.Context, clientIds []string, token string) ([]domain.Client, error)
}

type DefaultClientService struct {
	Configuration
}

func (d DefaultClientService) GetClients(ctx context.Context, token string) ([]domain.Client, error) {
	endpoint := fmt.Sprintf("%s/admin/realms/%s/clients", d.Configuration.GetBaseUrl(), d.Configuration.GetRealm())
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
		log.WithError(err).Error("could not get realm clients")
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("could not get clients -- see reason " + res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var clients []domain.Client
	err = json.Unmarshal(data, &clients)

	return clients, nil
}

func (d DefaultClientService) GetClientById(ctx context.Context, clientId, token string) (domain.Client, error) {
	clients, err := d.GetClients(ctx, token)
	if err != nil {
		return domain.Client{}, err
	}

	for _, c := range clients {
		if c.Id == clientId {
			return c, err
		}
	}

	return domain.Client{}, errors.New(fmt.Sprintf("could not find client: %s", clientId))
}

func (d DefaultClientService) GetClientByClientId(ctx context.Context, clientName, token string) (domain.Client, error) {
	clients, err := d.GetClients(ctx, token)
	if err != nil {
		return domain.Client{}, err
	}

	for _, c := range clients {
		if strings.EqualFold(c.ClientId, clientName) {
			return c, err
		}
	}

	return domain.Client{}, errors.New(fmt.Sprintf("could not find client: %s", clientName))
}

func (d DefaultClientService) GetClientsByIds(ctx context.Context, ids []string, token string) ([]domain.Client, error) {
	clients, err := d.GetClients(ctx, token)
	if err != nil {
		return nil, err
	}

	clientMap := make(map[string]domain.Client)

	for _, c := range clients {
		clientMap[c.Id] = c
	}

	var response []domain.Client
	for _, id := range ids {
		response = append(response, clientMap[id])
	}

	return response, nil
}

func (d DefaultClientService) GetClientsByClientIds(ctx context.Context, clientIds []string, token string) ([]domain.Client, error) {
	clients, err := d.GetClients(ctx, token)
	if err != nil {
		return nil, err
	}

	clientMap := make(map[string]domain.Client)

	for _, c := range clients {
		clientMap[c.Id] = c
	}

	var response []domain.Client
	for _, id := range clientIds {
		response = append(response, clientMap[id])
	}

	return response, nil
}
