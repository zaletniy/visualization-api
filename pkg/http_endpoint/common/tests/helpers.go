package testHelper

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"time"

	"visualization-api/pkg/grafanaclient/mock"
	"visualization-api/pkg/http_endpoint/authentication"
	"visualization-api/pkg/http_endpoint/common"
	"visualization-api/pkg/logging"
	"visualization-api/pkg/openstack/mock"
)

const tokenHeaderName = "Authorization"

type nullWriter struct{}

func (w nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// InitializeLogger initializes logger for tests with empty file writer
func InitializeLogger() {
	log.InitializeLogger(&nullWriter{}, false, "critical")
}

/*MockClientContainer returns struct populated with all mocks required*/
func MockClientContainer(mockCtrl *gomock.Controller) *common.ClientContainer {
	mockedOpenstack := mock_openstack.NewMockClientInterface(mockCtrl)
	mockedGrafana := mock_grafanaclient.NewMockSessionInterface(mockCtrl)
	return &common.ClientContainer{mockedOpenstack, mockedGrafana}
}

// GetAuthToken returns admin token with expiration date in 2037
func GetAuthToken(secret string, projectID string) string {
	parsedTime, _ := time.Parse(time.RFC3339, "2037-06-15T00:48:41Z")
	token, _ := httpAuth.JWTTokenFromParams(secret, true, projectID,
		parsedTime)
	return token
}

// SetRequestAuthHeader sets authorization bearer header for you
func SetRequestAuthHeader(secret string, projectID string, request *http.Request) {
	token := GetAuthToken(secret, projectID)
	request.Header.Set(tokenHeaderName, fmt.Sprintf("Bearer %s", token))
}
