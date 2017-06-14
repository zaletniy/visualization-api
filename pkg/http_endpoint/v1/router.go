package v1Api

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"
	"visualization-api/pkg/http_endpoint/authentication"
	"visualization-api/pkg/http_endpoint/common"
	"visualization-api/pkg/logging"
)

const adminAPIPrefix = "/admin"
const authPrefix = "/auth"

func authRouter(clients *common.ClientContainer,
	handler common.HandlerInterface, secret string) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/openstack", func(w http.ResponseWriter, r *http.Request) {
		// expected HEADER token name
		const openstackHeaderName = "X-OpenStack-Auth-Token"
		// response to be returned with 401 code to user
		const errorMsg = "Unauthorized. Token is invalid or expired."

		// check if openstack token is in request HEADERs
		openstackToken := r.Header.Get(openstackHeaderName)
		if openstackToken == "" {
			common.WriteErrorToResponse(w, http.StatusUnauthorized, errorMsg,
				fmt.Sprintf("header %s is not specified", openstackHeaderName))
			return
		}

		// try to authenticate with provided token
		token, err := handler.AuthOpenstack(clients, &common.RealClock{},
			openstackToken, secret)
		if err != nil {
			switch err.(type) {
			// common.InvalidOpenstackToken means, that user provided invalid
			// token. We have to return 401 error then
			case common.InvalidOpenstackToken:
				common.WriteErrorToResponse(w, http.StatusUnauthorized,
					errorMsg, err.Error())
				return
			// If any other error happened -> return 500 error
			default:
				log.Logger.Error(err)
				common.WriteErrorToResponse(w, http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
					"Internal server error occured")
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write(token)
	})
	return router
}

func adminRouter(authMiddleware func(http.Handler) http.Handler, secret string) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMiddleware)
	router.Use(httpAuth.AdminAuthenticationMiddleware(secret))
	/*
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("You are admin\n"))
			if err != nil {
			}
		})
	*/
	return router
}

// InitializeRouter initializes /v1 routers
func InitializeRouter(clients *common.ClientContainer,
	handler common.HandlerInterface, secret string) *chi.Mux {
	router := chi.NewRouter()
	authMiddleware := httpAuth.AuthenticationMiddleware(secret)
	router.Mount(adminAPIPrefix, adminRouter(authMiddleware, secret))
	router.Mount(authPrefix, authRouter(clients, handler, secret))
	return router
}
