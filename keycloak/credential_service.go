package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hub1989/keycloak-grpc-service/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type CredentialService interface {
	ObtainTokenForOps(ctx context.Context) (domain.AccessTokenResponse, error)
}

type DefaultCredentialService struct {
	Configuration
}

func (d DefaultCredentialService) ObtainTokenForOps(ctx context.Context) (domain.AccessTokenResponse, error) {
	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", d.GetBaseUrl(), d.GetRealm())

	credentials := d.Configuration.GetClientCredentials()

	requestForm := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", credentials.ClientId, credentials.ClientSecret)
	payload := strings.NewReader(requestForm)

	client := d.GetClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, payload)

	if err != nil {
		log.WithError(err).Error("could not get access token for operation")
		return domain.AccessTokenResponse{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("could not get access token for operation")
		return domain.AccessTokenResponse{}, err
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

	return accessToken, err
}
