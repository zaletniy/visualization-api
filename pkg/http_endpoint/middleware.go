package endpoint

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// ContextJWTProperty defines the key validated JWT token stored in http context
const ContextJWTProperty = "AuthToken"

func authenticationMiddleware(secret string) func(http.Handler) http.Handler {

	/*
	   go-jwt-middleware does not provide support for chi framework, that's
	   why this wrapper was implemented.

	   at the time of middleware initialization - jwtmiddleware.JWTMiddleware
	   struct is initialized in wrapped context. authenticationMiddleware
	   returns function compatible with chi framework interface. this function
	   uses JWTMiddleware functions to validate JWT token.
	*/

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  ContextJWTProperty,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if jwtMiddleware.CheckJWT(w, r) != nil {
				// 401 code and error message is set in http response by
				// jwtMiddleware, we have just to return response
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
