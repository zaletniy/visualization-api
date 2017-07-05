package proxy

import (
  "net/http"
  "fmt"
  "github.com/shuaiming/mung/middlewares"

  log "visualization-api/pkg/logging"
)

type OpenStackHandler struct {
}

const OpenStackUser = "OS_USER"

func NewOpenStackHandler() (*OpenStackHandler, error) {
  return &OpenStackHandler{}, nil
}

func (oh *OpenStackHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  log.Logger.Debugf("Checking the session")
  sess := middlewares.GetSession(r)
  log.Logger.Debugf("Session values for request is %v", sess.Values)

  if r.RequestURI == "/auth/openstack" {
    if r.Method == http.MethodGet {
      fmt.Fprint(rw, "<form action=\"/auth/openstack\">"+
        "username:<br><input type=\"text\" name=\"username\">"+
        "password:<br><input type=\"password\" name=\"password\">"+
        "<input type=\"submit\" value=\"Login\"></form>")
      return
    } else if r.Method == http.MethodPost {
      // processing credentials and setting OS_USER if needed
      //proces username/password form here
      log.Logger.Debugf("OpenStack user is setting as 'john'")
      sess.Values[OpenStackUser] = "john"
      //saving cookies
      sess.Save(r, rw)
      next(rw, r)
      return

    } else {
      http.Error(rw, "Unexpected method", http.StatusBadRequest)
      return
    }
  }

  if data, ok := sess.Values[OpenStackUser]; ok {
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
  // * redirect to auth form if not
  // * pass if context is initialized
}
