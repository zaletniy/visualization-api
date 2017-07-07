package v1Api

import (
	"encoding/json"
	"strconv"
	"time"
	"visualization-api/pkg/http_endpoint/authentication"
	"visualization-api/pkg/http_endpoint/common"
	"visualization-api/pkg/logging"
)

// TokenIssueHours defines on how much hours our token would be issued
const TokenIssueHours = 3

// V1Handler is implementation of Handler interface
type V1Handler struct{}

// AuthOpenstack uses provided keystone token to create jwt token
func (h *V1Handler) AuthOpenstack(clients *common.ClientContainer,
	clock common.ClockInterface, openstackToken string,
	secret string) ([]byte, error) {

	tokenValid, err := clients.Openstack.ValidateToken(openstackToken)
	if err != nil {
		log.Logger.Errorf("Error validating openstack Token: %s", err)
		return nil, err
	}
	if !tokenValid {
		return nil, common.InvalidOpenstackToken{}
	}

	tokenInfo, err := clients.Openstack.GetTokenInfo(openstackToken)
	if err != nil {
		log.Logger.Errorf("Error retrieving openstack Token: %s", err)
		return nil, err
	}

	expirationTime := clock.Now().Add(TokenIssueHours * time.Hour)

	grafanaOrg, err := clients.Grafana.GetOrCreateOrgByName(
		tokenInfo.ProjectName + "-" + tokenInfo.ProjectID)
	if err != nil {
		return nil, err
	}
	grafanaOrgID := strconv.Itoa(grafanaOrg.ID)

	token, err := httpAuth.JWTTokenFromParams(secret, tokenInfo.IsAdmin(),
		grafanaOrgID, expirationTime)
	if err != nil {
		return nil, err
	}

	var payload struct {
		JWT   string `json:"jwt"`
		Token struct {
			OrganizationID string    `json:"organizationId"`
			ExpiresAt      time.Time `json:"expiresAt"`
			IsAdmin        bool      `json:"isAdmin"`
		} `json:"token"`
	}

	payload.JWT = token
	payload.Token.OrganizationID = grafanaOrgID
	payload.Token.ExpiresAt = expirationTime
	payload.Token.IsAdmin = tokenInfo.IsAdmin()

	return json.Marshal(payload)
}
