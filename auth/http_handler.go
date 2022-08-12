package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ijalalfrz/majoo-test-1/model"
	"github.com/ijalalfrz/majoo-test-1/response"
	"github.com/sirupsen/logrus"
)

const (
	basePath = "/majoo-service"
)

// HTTPHandler is an http handler concrete struct.
type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

// NewAuthHTTPHandler is a constructor
func NewAuthHTTPHandler(logger *logrus.Logger, vld *validator.Validate, router *mux.Router, usecase Usecase) {

	handler := &HTTPHandler{
		Logger:   logger,
		Usecase:  usecase,
		Validate: vld,
	}

	router.HandleFunc(basePath+"/v1/login", handler.Login).
		Methods(http.MethodPost)

}

// Login is a function that handle incomming request
func (handler HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var payload model.UserLoginParams

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	if err := handler.validateRequestBody(payload); err != nil {
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = handler.Usecase.Login(ctx, payload.UserName, payload.Password)
	response.JSON(w, resp)
	return
}

func (handler HTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = handler.Validate.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("Invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
