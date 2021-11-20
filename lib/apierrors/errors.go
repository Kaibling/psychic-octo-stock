package apierrors

import "net/http"

type ApiError interface {
	Error() string
	HttpStatus() int
}

type NotFoundError struct {
	HttpStatusCode int
	Err            error
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{HttpStatusCode: http.StatusNotFound, Err: err}
}
func (e *NotFoundError) Error() string   { return e.Err.Error() }
func (e *NotFoundError) HttpStatus() int { return e.HttpStatusCode }

type GeneralError struct {
	HttpStatusCode int
	Err            error
	//ErrorType string
}

func NewGeneralError(err error) *GeneralError {
	return &GeneralError{HttpStatusCode: http.StatusBadGateway, Err: err}
}
func (e *GeneralError) Error() string   { return e.Err.Error() }
func (e *GeneralError) HttpStatus() int { return e.HttpStatusCode }

type ClientError struct {
	HttpStatusCode int
	Err            error
	//ErrorType string
}

func NewClientError(err error) *ClientError {
	return &ClientError{HttpStatusCode: http.StatusUnprocessableEntity, Err: err}
}
func (e *ClientError) Error() string   { return e.Err.Error() }
func (e *ClientError) HttpStatus() int { return e.HttpStatusCode }
