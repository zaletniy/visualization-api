package proxy

import (
	"fmt"
	"github.com/shuaiming/mung/middlewares"
	"net/http"
	"time"

	log "visualization-api/pkg/logging"
)

var (
	// DefaultOpenStackGrafanaRolesMapping default roles mapping
	DefaultOpenStackGrafanaRolesMapping = map[string]string{"admin": GrafanaRoleEditor,
		"Member": GrafanaRoleReadOnlyEditor}
)

// OpenStackAuthHandler middleware for handling authentication via Keystone
type OpenStackAuthHandler struct {
	loginPage       []byte
	grafanaStateTTL int
	rolesMapping    map[string]string
}

// NewOpenStackAuthHandler returns OpenStackAuthHandler
func NewOpenStackAuthHandler(loginPage []byte, grafanaStateTTL int, rolesMapping map[string]string) (*OpenStackAuthHandler, error) {
	if rolesMapping == nil {
		return &OpenStackAuthHandler{loginPage: loginPage,
			grafanaStateTTL: grafanaStateTTL, rolesMapping: DefaultOpenStackGrafanaRolesMapping}, nil
	}
	return &OpenStackAuthHandler{loginPage: loginPage,
		grafanaStateTTL: grafanaStateTTL, rolesMapping: rolesMapping}, nil
}

func (oh *OpenStackAuthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sess := middlewares.GetSession(r)
	log.Logger.Debugf("Session values for request are %v", sess.Values)

	if r.RequestURI == "/auth/openstack" {
		if r.Method == http.MethodGet {
			r.Header.Set("Content-Type", "text/html")
			_, err := rw.Write(oh.loginPage)
			if err != nil {
				log.Logger.Errorf("Error during writing response err: %s", err)
			}
			return

		} else if r.Method == http.MethodPost {

			username := r.FormValue("username")
			password := r.FormValue("password")

			ok, err := oh.authenticate(username, password)
			if err != nil {
				log.Logger.Infof("User %s is not authenticated. Due to err: %s", username, err)
				http.Redirect(rw, r, "/auth/openstack", http.StatusInternalServerError)
			}

			if !ok {
				log.Logger.Infof("User %s is not authenticated", username)
				http.Redirect(rw, r, "/auth/openstack", http.StatusForbidden)
			}

			log.Logger.Debugf("%s is setting as '%s'", SessionUsername, username)

			sess.Values[SessionUsername] = username
			sess.Values[GrafanaUpdateCommandSessionKey], err = oh.getGrafanaUpdateCommand(username)
			sess.Values[OrgAndUserStateExpiresAt] = time.Now().
				Add(time.Duration(oh.grafanaStateTTL) * time.Second).
				Format(TimeFormat)

			if err != nil {
				log.Logger.Errorf("Can't create CGafanaUpdateCommand for user %s err: %s", username, err)
			}
			//saving cookies/session
			err = sess.Save(r, rw)
			if err != nil {
				log.Logger.Errorf("Error during updating session err: %s", err)
				http.Error(rw, "Internal error", http.StatusInternalServerError)
				return
			}
			//We are good now 302 redirect
			http.Redirect(rw, r, "/", http.StatusFound)
			return

		} else {
			http.Error(rw, "Unexpected method", http.StatusBadRequest)
			return
		}
	}

	if data, ok := sess.Values[SessionUsername]; ok {
		// looks like user already authenticated, let us check how fresh
		// is Grafana state
		username := fmt.Sprintf("%s", data)

		if _, ok := sess.Values[OrgAndUserStateExpiresAt]; !ok {
			log.Logger.Errorf("%s key expected in session", OrgAndUserStateExpiresAt)
			http.Error(rw, "Internal error", http.StatusInternalServerError)
			return
		}

		gTTL := sess.Values[OrgAndUserStateExpiresAt]

		t, err := time.Parse(TimeFormat, fmt.Sprintf("%s", gTTL))
		if err != nil {
			log.Logger.Errorf("Can't parse time %s", gTTL)
			http.Error(rw, "Internal error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(t) {
			// grafana data is outdated let us refresh if
			//TODO(illia) looks like code duplication. Should be refactored
			sess.Values[GrafanaUpdateCommandSessionKey], err = oh.getGrafanaUpdateCommand(username)
			sess.Values[OrgAndUserStateExpiresAt] = time.Now().
				Add(time.Duration(oh.grafanaStateTTL) * time.Second).
				Format(TimeFormat)

			if err != nil {
				log.Logger.Errorf("Can't create GafanaUpdateCommand for user %s. err: %s", username, err)
			}
			//saving cookies/session
			err := sess.Save(r, rw)
			if err != nil {
				log.Logger.Errorf("Error during updating session err: %s", err)
				http.Error(rw, "Internal error", http.StatusInternalServerError)
				return
			}
		}

		log.Logger.Debugf("OpenStack user is defined in session as %s", data)
		next(rw, r)
		return
	}
	//let user enter credentials
	http.Redirect(rw, r, "/auth/openstack", http.StatusFound)
	return
}

func (oh *OpenStackAuthHandler) getGrafanaUpdateCommand(user string) (GrafanaUpdateCommand, error) {
	//TODO implement it
	return GrafanaUpdateCommand{
		User{
			Login: user,
		},
		[]Organization{Organization{"ww", GrafanaRoleReadOnlyEditor}},
	}, nil
}

func (oh *OpenStackAuthHandler) authenticate(user string, password string) (bool, error) {
	//TODO implement it
	return true, nil
}
