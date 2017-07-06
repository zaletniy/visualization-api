package proxy

import (
	"github.com/shuaiming/mung/middlewares"
	"net/http"

	log "visualization-api/pkg/logging"
)

type OpenStackHandler struct {
	loginPage []byte
}

func NewOpenStackHandler(loginPage []byte) (*OpenStackHandler, error) {
	return &OpenStackHandler{loginPage: loginPage}, nil
}

func (oh *OpenStackHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sess := middlewares.GetSession(r)
	log.Logger.Debugf("Session values for request are %v", sess.Values)

	if r.RequestURI == "/auth/openstack" {
		if r.Method == http.MethodGet {
			r.Header.Set("Content-Type", "text/html")
			rw.Write(oh.loginPage)
			return
		} else if r.Method == http.MethodPost {
			// processing credentials and setting OS_USER if needed
			//proces username/password form here
			log.Logger.Debugf("OpenStack user is setting as 'john'")
			sess.Values[SessionUsername] = "john"
			//saving cookies
			sess.Save(r, rw)
			//We are good now 302 redirect
			http.Redirect(rw, r, "/", http.StatusFound)
			return

		} else {
			http.Error(rw, "Unexpected method", http.StatusBadRequest)
			return
		}
	}

	if data, ok := sess.Values[SessionUsername]; ok {
		log.Logger.Debugf("OpenStack user is defined in session as %s", data)
		next(rw, r)
		return
	} else {
		//let user enter credentials
		http.Redirect(rw, r, "/auth/openstack", http.StatusFound)
		return
	}

	//TODO(illia):
	// * check if user authenticated
}
