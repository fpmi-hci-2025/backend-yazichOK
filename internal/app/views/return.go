package views

import (
	"encoding/json"
	"errors"
	"net/http"

	"speech-processing-service/internal/errs"

	"go.uber.org/zap"
)

const (
	//internalServerError group of errors
	codeMinio          = 10
	codeExecutionQuery = 11

	//bad request group of errors
	codeTypeMustBeNumeric = 20

	//not found group of errors
	codeNotFound = 30

	//unknowError
	codeUnknown = 999
)

func Return(
	logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
	err error,
) {
	if err != nil {
		errResp := ErrorResponse{
			Error: Error{
				ErrorCode: defineErrorCode(err),
				Msg:       err.Error(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(defineStatusCode(err))

		if err := json.NewEncoder(w).Encode(errResp); err != nil {
			logger.Error("views.Return", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	resp := SuccessResponse{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("views.Return", zap.Error(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	}

}

func defineStatusCode(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, errs.ErrExecutionQuery) || errors.Is(err, errs.ErrMinio):
		return http.StatusInternalServerError
	case errors.Is(err, errs.ErrTypeMustBeNumeric) || errors.Is(err, errs.ErrDecodingJSON) ||
		errors.Is(err, errs.ErrTypeMustBeUUID):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func defineErrorCode(err error) int {
	switch {
	case errors.Is(err, errs.ErrExecutionQuery):
		return codeExecutionQuery
	case errors.Is(err, errs.ErrMinio):
		return codeMinio
	case errors.Is(err, errs.ErrTypeMustBeNumeric) || errors.Is(err, errs.ErrDecodingJSON) ||
		errors.Is(err, errs.ErrTypeMustBeUUID):
		return codeTypeMustBeNumeric
	case errors.Is(err, errs.ErrNotFound):
		return codeNotFound
	default:
		return codeUnknown
	}
}
