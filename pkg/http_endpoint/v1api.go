package endpoint

import (
	"github.com/pressly/chi"
	"net/http"
)

const adminAPIPrefix = "/admin"

func adminRouter() http.Handler {
	router := chi.NewRouter()
	return router
}

func v1Router() http.Handler {
	router := chi.NewRouter()
	/* example of router handler usage
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Logger.Error(err)
		}
	})
	*/

	// example of subrouter usage. all datasources, visualizations, templates
	// have to use separate subrouter
	router.Mount(adminAPIPrefix, adminRouter())
	return router
}
