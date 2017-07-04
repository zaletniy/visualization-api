package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "visualization-api/pkg/logging"
	//"github.com/shuaiming/mung/middlewares"
)

// Proxy wapper ReverseProxy
type Proxy struct {
	proxy          *httputil.ReverseProxy
	requestLogging bool
	authHeader     string
}

// NewProxy New proxy
func NewProxy(endpoint string, requestLogging bool, authHeader string) (*Proxy, error) {
	backend, err := url.Parse(fmt.Sprintf("%s", endpoint))
	if err != nil {
		//TODO(illia) refactor logging package. Not very convenient to call log.Logger...
		log.Logger.Error("Can't parse Grafana url err: %s", err)
		return nil, err
	}
	log.Logger.Debugf("Grafana reverse proxy configured for endpoint=%s", endpoint)
	proxy := httputil.NewSingleHostReverseProxy(backend)
	return &Proxy{proxy: proxy, requestLogging: requestLogging, authHeader: authHeader}, nil
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//sess := middlewares.GetSession(r)
	//openid := middlewares.GetOpenIDUser(r)

	// overwrite grafana's login
	//if r.URL.Path == "/login" {
	//  http.Redirect(rw, r, "/openid/login", http.StatusFound)
	//  return
	//}

	// overwrite grafana's logout
	if r.URL.Path == "/logout" {
		//TODO(illia) remove cookies and clear session here
		//  delete(sess.Values, middlewares.OpenIDContextKey)
		//  if err := sess.Save(r, rw); err != nil {
		//    fmt.Fprintln(os.Stderr, err.Error())
		log.Logger.Infof("Logging out user")
		fmt.Fprintf(rw, "<h1>Bye bye</h1><div>%s</div>", "user1")
		return
	}

	log.Logger.Debugf("Proxying %s %s", r.Method, r.RequestURI)

	// redirect to login url if openid not login
	//email, ok := openid["sreg.email"]
	//if !ok {
	//  http.Redirect(rw, r, "/openid/login", http.StatusFound)
	//  return
	//}

	// overwirte X-WEBAUTH-USER with openid email name
	r.Header.Set(p.authHeader, "user1")

	if p.requestLogging {
		log.Logger.Debugf("Request dump  %v", r)
	}
	p.proxy.ServeHTTP(rw, r)
}
