package proxy

import (
	"github.com/shuaiming/mung/middlewares"
	"net/http"
	log "visualization-api/pkg/logging"
)

// VisualizationAPIMiddleware creates users and organizations in Grafana
// via visualization-api according to model in session
type VisualizationAPIMiddleware struct {
}

// NewVisualizationAPIMiddleware returns middleware for managing users and
// organizations in Grafana via Visualization API
func NewVisualizationAPIMiddleware() (*VisualizationAPIMiddleware, error) {
	return &VisualizationAPIMiddleware{}, nil
}

func (vm *VisualizationAPIMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Logger.Debugf("We are creating all users and organizations in this middleware")
	sess := middlewares.GetSession(r)

	if cmd, ok := sess.Values[GrafanaUpdateCommandSessionKey]; ok {
		log.Logger.Debugf("Visualization API command is %+v", cmd)

		//TODO(illia) code should be here :)
		//creating everything here based on GrafanaUpdateCommand

		// * check if user already exists in Grafana in order to not create it automatically
		// * check if organizaton exists
		// * give user roles
		// * remove user's roles if needed

		delete(sess.Values, GrafanaUpdateCommandSessionKey)
		//saving cookies/session
		err := sess.Save(r, rw)
		if err != nil {
			log.Logger.Errorf("Can't save session change err: %s", err)
			http.Error(rw, "Internal error", http.StatusInternalServerError)
			return
		}
		next(rw, r)
		return
	}

	log.Logger.Debugf("Visualization API has been configured already for user. Skipping this step.")
	next(rw, r)
	return
}
