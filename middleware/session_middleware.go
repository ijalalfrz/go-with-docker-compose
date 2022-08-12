package middleware

import (
	"net/http"
	"strconv"
	"strings"

	gctx "github.com/gorilla/context"
	"github.com/ijalalfrz/majoo-test-1/entity"

	"github.com/ijalalfrz/majoo-test-1/exception"
	"github.com/ijalalfrz/majoo-test-1/model"
	"github.com/ijalalfrz/majoo-test-1/response"

	"github.com/ijalalfrz/majoo-test-1/jwt"
)

var (
	header = "Authorization"
)

// Session is concrete struct of jwt authorization.
type Session struct {
	jsonWebToken *jwt.JSONWebToken
}

// NewSessionMiddleware is a constructor.
func NewSessionMiddleware(jsonWebToken *jwt.JSONWebToken) RouteMiddleware {
	return &Session{jsonWebToken}
}

// Verify will verify the http incomming request to ensure it comes within the authorized token.
func (a Session) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get(header)
		if authHeader == "" {
			a.respondUnauthorized(w, jwt.ErrInvalidToken.Error())

			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			a.respondUnauthorized(w, jwt.ErrInvalidToken.Error())

			return
		}

		tokenString := bearerToken[1]

		var customerBearerClaims model.UserBeaerClaims

		if err := a.jsonWebToken.Parse(ctx, tokenString, &customerBearerClaims); err != nil {
			a.respondUnauthorized(w, err.Error())
			return
		}

		userId, _ := strconv.ParseInt(customerBearerClaims.Subject, 10, 64)

		user := entity.User{
			ID: userId,
		}

		gctx.Set(r, "user", user)

		next(w, r)
	})
}

func (a Session) respondUnauthorized(w http.ResponseWriter, message string) {
	resp := response.NewErrorResponse(exception.ErrUnauthorized, http.StatusUnauthorized, nil, response.StatUnauthorized, message)
	response.JSON(w, resp)
}
