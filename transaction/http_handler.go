package transaction

import (
	"net/http"
	"strconv"

	gctx "github.com/gorilla/context"

	"github.com/gorilla/mux"
	"github.com/ijalalfrz/majoo-test-1/entity"
	"github.com/ijalalfrz/majoo-test-1/middleware"
	"github.com/ijalalfrz/majoo-test-1/response"
	"github.com/sirupsen/logrus"
)

const (
	basePath = "/majoo-service"
)

// HTTPHandler is an http handler concrete struct.
type HTTPHandler struct {
	Logger  *logrus.Logger
	Usecase Usecase
}

// NewTransactionHTTPHandler is a constructor
func NewTransactionHTTPHandler(logger *logrus.Logger, sessionMiddleware middleware.RouteMiddleware, router *mux.Router, usecase Usecase) {

	handler := &HTTPHandler{
		Logger:  logger,
		Usecase: usecase,
	}

	router.HandleFunc(basePath+"/v1/transactions", sessionMiddleware.Verify(handler.FindMany)).
		Methods(http.MethodGet)

	router.HandleFunc(basePath+"/v1/transactions/omzet/per-merchant", sessionMiddleware.Verify(handler.FindManyOmzetPerMerchant)).
		Methods(http.MethodGet).Queries("startDate", "{startDate}", "endDate", "{endDate}", "size", "{size:[0-9]+}", "page", "{page:[0-9]+}")

	router.HandleFunc(basePath+"/v1/transactions/omzet/per-outlet", sessionMiddleware.Verify(handler.FindManyOmzetPerOutlet)).
		Methods(http.MethodGet).Queries("startDate", "{startDate}", "endDate", "{endDate}", "size", "{size:[0-9]+}", "page", "{page:[0-9]+}")

}

// FindMany is a function that handle incomming request
func (handler HTTPHandler) FindMany(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	ctx := r.Context()
	user := gctx.Get(r, "user").(entity.User)

	resp = handler.Usecase.FindManyByUserID(ctx, user.ID)
	response.JSON(w, resp)
	return
}

// FindManyOmzetPerMerchant is a function that handle incomming request
func (handler HTTPHandler) FindManyOmzetPerMerchant(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	ctx := r.Context()

	user := gctx.Get(r, "user").(entity.User)

	queryString := r.URL.Query()
	startDate := queryString.Get("startDate")
	endDate := queryString.Get("endDate")
	size := 10
	page := 1
	if queryString.Get("size") != "" {
		size, _ = strconv.Atoi(queryString.Get("size"))
	}
	if queryString.Get("page") != "" {
		page, _ = strconv.Atoi(queryString.Get("page"))
	}

	resp = handler.Usecase.FindOmzetPerMerchantByUserID(ctx, user.ID, startDate, endDate, size, page)
	response.JSON(w, resp)
	return
}

// FindManyOmzetPerOutlet is a function that handle incomming request
func (handler HTTPHandler) FindManyOmzetPerOutlet(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	ctx := r.Context()

	user := gctx.Get(r, "user").(entity.User)

	queryString := r.URL.Query()
	startDate := queryString.Get("startDate")
	endDate := queryString.Get("endDate")
	size := 10
	page := 1
	if queryString.Get("size") != "" {
		size, _ = strconv.Atoi(queryString.Get("size"))
	}
	if queryString.Get("page") != "" {
		page, _ = strconv.Atoi(queryString.Get("page"))
	}

	resp = handler.Usecase.FindOmzetPerOutletByUserID(ctx, user.ID, startDate, endDate, size, page)
	response.JSON(w, resp)
	return
}
