package endpoint

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"
)

const v1ApiPrefix = "/v1"

// Serve is an entry point to our HTTP API
func Serve(secret string, httpPort int) error {
	rootRouter := chi.NewRouter()
	rootRouter.Use(authenticationMiddleware(secret))
	rootRouter.Mount(v1ApiPrefix, v1Router())
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), rootRouter)
}
