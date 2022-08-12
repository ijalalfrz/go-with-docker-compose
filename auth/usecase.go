package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/ijalalfrz/majoo-test-1/exception"
	"github.com/ijalalfrz/majoo-test-1/jwt"
	"github.com/ijalalfrz/majoo-test-1/model"
	"github.com/ijalalfrz/majoo-test-1/response"
	"github.com/sirupsen/logrus"
)

type Usecase interface {
	Login(ctx context.Context, username string, password string) (resp response.Response)
}

const (
	loginUnexpectedErrMessage = "Unexpected error while trying to login"
	loginNotFoundErrMessage   = "User not found"
	loginInvalidErrMessage    = "Invalid username or password"
	loginSuccessMessage       = "User credential"
)

type usecase struct {
	serviceName string
	logger      *logrus.Logger
	repository  Repository
	jwt         *jwt.JSONWebToken
}

func NewAuthUsecase(property UsecaseProperty) Usecase {
	return &usecase{
		serviceName: property.ServiceName,
		logger:      property.Logger,
		repository:  property.Repository,
		jwt:         property.JWT,
	}
}

func (u usecase) Login(ctx context.Context, username string, password string) (resp response.Response) {

	user, err := u.repository.FindByUsername(ctx, username)
	if err != nil {
		u.logger.Error(err)
		if err != exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, loginUnexpectedErrMessage)
		}
		return response.NewErrorResponse(err, http.StatusNotFound, nil, response.StatNotFound, loginNotFoundErrMessage)

	}

	passwordHash := md5.Sum([]byte(password))
	passwordHashHex := hex.EncodeToString(passwordHash[:])

	if user.Password != passwordHashHex {
		return response.NewErrorResponse(err, http.StatusUnauthorized, nil, response.StatUnauthorized, loginInvalidErrMessage)
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7).Unix()
	claims := model.UserBeaerClaims{}
	claims.UserName = user.UserName
	claims.Subject = fmt.Sprintf("%d", user.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = expiresAt
	claims.Issuer = "majoo"
	tokenString, err := u.jwt.Sign(ctx, claims)
	if err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, err.Error())
	}

	var userSignedIn model.UserSuccessSignIn
	userSignedIn.Name = user.Name
	userSignedIn.Token = &tokenString

	return response.NewSuccessResponse(userSignedIn, response.StatOK, loginSuccessMessage)

}
