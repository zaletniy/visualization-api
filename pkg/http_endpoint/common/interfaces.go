package common

import (
	"time"
	"visualization-api/pkg/openstack"
)

/*ClientContainer represents container for storing different clients
It was created to have mockable architecture*/
type ClientContainer struct {
	Openstack openstack.ClientInterface
}

/*HandlerInterface represents set of handlers for api
It was created to have mockable architecture*/
type HandlerInterface interface {
	AuthOpenstack(*ClientContainer, ClockInterface, string, string) ([]byte, error)
}

// ClockInterface serves for testing purposes of functions, that require time
type ClockInterface interface {
	Now() time.Time
}
