package utils

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// JSONError custom implementation of JSONAPI error
type JSONError struct {
	Code   int    `json:"code"`      // an application-specific error code
	Title  string `json:"codeTitle"` // a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
	Source string `json:"source"`    // an object containing references to the source of the error
	Form   string `json:"form"`      // an object containing references to the source of the error

	Detail string `json:"detail"` // a human-readable explanation specific to this occurrence of the problem.
}

// JSONErrors for multiple return
type JSONErrors struct {
	Errors []*JSONError `json:"errors"`
}

// Satisfy error interface
func (e JSONErrors) Error() string {
	res, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(res)
}

// DEPRICATE NewError creates new error
func NewError(code int, source, detail string) *JSONError {
	e := new(JSONError)
	e.Code = code
	e.Title = StatusTitle(code)
	e.Source = source
	e.Detail = detail
	// send every created error to elastic
	// logger.Error
	Error("json error -> ", e)
	return e
}

// CreateError creates new error
func CreateError(code int, source, form, detail string) *JSONError {
	e := new(JSONError)
	e.Code = code
	e.Title = StatusTitle(code)
	e.Source = source
	e.Form = form
	e.Detail = detail
	// send every created error to elastic
	// logger.Error
	Error("json error -> ", e)
	return e
}

// ResponseErrors with list or errors
func ResponseErrors(c echo.Context, errors ...*JSONError) error {
	return c.JSON(http.StatusBadRequest, JSONErrors{Errors: errors})
}

// ResponseData with success
func ResponseData(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// Response based on args
func Response(c echo.Context, data interface{}, errors ...*JSONError) error {
	if data != nil {
		return ResponseData(c, data)
	}
	return ResponseErrors(c, errors...)
}
