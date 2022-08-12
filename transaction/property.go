package transaction

import "github.com/sirupsen/logrus"

type UsecaseProperty struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
}
