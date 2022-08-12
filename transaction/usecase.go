package transaction

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/ijalalfrz/majoo-test-1/exception"
	"github.com/ijalalfrz/majoo-test-1/model"
	"github.com/ijalalfrz/majoo-test-1/response"
	"github.com/sirupsen/logrus"
)

const (
	RFC3339Millis string = "2006-01-02T15:04:05.999Z"
	NormalFormat  string = "2006-01-02"
)

type Usecase interface {
	FindManyByUserID(ctx context.Context, userId int64) (resp response.Response)
	FindOmzetPerMerchantByUserID(ctx context.Context, userId int64, startDate string, endDate string, size int, page int) (resp response.Response)
	FindOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate string, endDate string, size int, page int) (resp response.Response)
}

// collection of message
const (
	transactionUnexpectedErrMessage = "An error occured while getting transaction"
	transactionNotFoundErrMessage   = "Transaction not found"

	transactionSuccessMessage = "List of transaction"
)

type usecase struct {
	serviceName string
	logger      *logrus.Logger
	repository  Repository
}

func NewTransactionUsecase(property UsecaseProperty) Usecase {
	return &usecase{
		serviceName: property.ServiceName,
		logger:      property.Logger,
		repository:  property.Repository,
	}
}

func (u usecase) FindManyByUserID(ctx context.Context, userId int64) (resp response.Response) {
	transactions, err := u.repository.FindManyByUserID(ctx, userId, nil, nil, 10)
	if err != nil {
		u.logger.Error(err)
		if err == exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusNotFound, nil, response.StatNotFound, transactionNotFoundErrMessage)
		}
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, transactionUnexpectedErrMessage)
	}

	return response.NewSuccessResponse(transactions, response.StatOK, transactionSuccessMessage)
}

func (u usecase) FindOmzetPerMerchantByUserID(ctx context.Context, userId int64, startDate string, endDate string, size int, page int) (resp response.Response) {

	start, _ := time.Parse(NormalFormat, startDate)
	end, _ := time.Parse(NormalFormat, endDate)
	diff := end.Sub(start).Hours() / 24

	skip := (page - 1) * size
	startTime := start.AddDate(0, 0, skip)
	endTime := startTime.AddDate(0, 0, size-1)
	totalPage := math.Ceil(diff / float64(size))

	if page > int(totalPage) {
		return response.NewErrorResponse(exception.ErrNotFound, http.StatusNotFound, nil, response.StatNotFound, transactionNotFoundErrMessage)
	}

	omzets, err := u.repository.FindOmzetPerMerchantByUserID(ctx, userId, &startTime, &endTime)
	omzetMap := map[string]model.GroupOmzetPerMerchant{}

	for _, o := range omzets {
		omzetMap[o.TransactionFullDate.Format(NormalFormat)] = o
	}

	var result []model.GroupOmzetPerMerchant
	for d := startTime; d.After(endTime) == false; d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(NormalFormat)
		_, ok := omzetMap[dateStr]
		if !ok {
			result = append(result, model.GroupOmzetPerMerchant{
				TransactionFullDate: d,
			})
		} else {
			result = append(result, omzetMap[dateStr])
		}
	}
	if err != nil {
		u.logger.Error(err)
		if err != exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, transactionUnexpectedErrMessage)
		}
	}
	meta := map[string]interface{}{
		"total_data": int(diff),
		"total_page": totalPage,
		"page":       page,
	}

	return response.NewSuccessResponseWithMeta(result, response.StatOK, transactionSuccessMessage, meta)
}

func (u usecase) FindOmzetPerOutletByUserID(ctx context.Context, userId int64, startDate string, endDate string, size int, page int) (resp response.Response) {

	start, _ := time.Parse(NormalFormat, startDate)
	end, _ := time.Parse(NormalFormat, endDate)
	skip := (page - 1) * size

	omzets, err := u.repository.FindOmzetPerOutletByUserID(ctx, userId, &start, &end, skip, size)

	if err != nil {
		u.logger.Error(err)
		if err != exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, transactionUnexpectedErrMessage)
		}
		return response.NewErrorResponse(exception.ErrNotFound, http.StatusNotFound, nil, response.StatNotFound, transactionNotFoundErrMessage)
	}
	meta := map[string]interface{}{}

	return response.NewSuccessResponseWithMeta(omzets, response.StatOK, transactionSuccessMessage, meta)
}
