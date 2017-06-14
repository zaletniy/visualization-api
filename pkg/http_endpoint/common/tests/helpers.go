package testHelper

import (
	"github.com/golang/mock/gomock"
	"visualization-api/pkg/http_endpoint/common"
	"visualization-api/pkg/logging"
	"visualization-api/pkg/openstack/mock"
)

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
	return &common.ClientContainer{mockedOpenstack}
}
