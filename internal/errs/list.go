package errs

import "errors"

var (
	ErrInitialization = errors.New("object initialization error")

	//Internal server errors
	ErrExecutionQuery       = errors.New("query execution error")
	ErrUseCaseExecution     = errors.New("use case execution error")
	ErrMinio                = errors.New("minio error")
	ErrMarshalingJSON       = errors.New("marshaling json error")
	ErrExecutionRequest     = errors.New("request execution error")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")

	//Bad request errors
	ErrTypeMustBeNumeric = errors.New("type must be numeric")
	ErrTypeMustBeUUID    = errors.New("type must be uuid")
	ErrDecodingJSON      = errors.New("decoding json error")

	//Not found errors
	ErrNotFound = errors.New("not found")
)
