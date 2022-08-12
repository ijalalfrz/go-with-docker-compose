package auth

import (
	"github.com/ijalalfrz/majoo-test-1/jwt"
	"github.com/sirupsen/logrus"
)

type UsecaseProperty struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
	JWT         *jwt.JSONWebToken
}
