// Package bindvalidate is package of validators and sanitizers for strings, structs and collections.
// Form Post is removed as we use JSON API only
package bindvalidate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	valid "github.com/reactima/reactima-go/validator"
)

// BindAndValidate struct
type BindAndValidate struct{}

// JSONError custom implementation of JSONAPI error
type JSONError struct {
	Source string `json:"source"` // an object containing references to the source of the error
	Detail string `json:"detail"` // a human-readable explanation specific to this occurrence of the problem.
}

// JSONErrors for multiple return
type JSONErrors struct {
	Errors []*JSONError `json:"errors"`
}

func (e JSONErrors) Error() string {
	//res, err := json.Marshal(e)
	//if err != nil {
	//	return err.Error()
	//}
	//return string(res)

	res := ""
	for _, error := range e.Errors {
		//res += error.Source + ": " + error.Detail
		res += error.Detail + "\n"
	}
	return res
}

// NewError creates new error
func NewError(source, detail string) *JSONError {
	e := new(JSONError)
	e.Source = source
	e.Detail = detail
	return e
}

// Bind override
// TODO Document Da
func (cb *BindAndValidate) Bind(i interface{}, c echo.Context) error {
	jsonErrors := JSONErrors{}
	req := c.Request()

	ctype := req.Header.Get(echo.HeaderContentType)
	if req.ContentLength == 0 {
		jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.ContentLength", "Request body can't be empty"))
		return jsonErrors
	}

	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):

		fmt.Println(req.Body)

		if err := json.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {

				jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.UnmarshalTypeError",
					fmt.Sprintf("UnmarshalTypeError: expected=%v, got=%v, offset=%v Struct=%v Field=%v", ute.Type, ute.Value, ute.Offset, ute.Struct, ute.Field)))
				return jsonErrors

			} else if se, ok := err.(*json.SyntaxError); ok {

				jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.SyntaxError",
					fmt.Sprintf("SyntaxError: offset=%v, error=%v", se.Offset, se.Error())))
				return jsonErrors

			} else {

				jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.DecodeError",
					"Possible *int64, null related error. Remove unused values. DecodeError: "+err.Error()))
				return jsonErrors
			}
		}

	default:

		jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.ErrUnsupportedMediaType",
			"Unsupported Media Type or missed JSON encording: "+echo.ErrUnsupportedMediaType.Error()))
		return jsonErrors

	}

	// assuming POST method used to create entity check ro required fields is needed
	if c.Request().Method == echo.POST {
		v := reflect.ValueOf(i).Elem()
		t := reflect.TypeOf(i).Elem()
		missing := make([]string, 0)

		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Tag.Get("create") == "required" && v.Field(i).IsNil() == true {
				missing = append(missing, t.Field(i).Name)
			}
		}

		if len(missing) != 0 {
			for _, m := range missing {

				jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.MissingRequiredField",
					"ValidateStructErrorMissingRequiredField -> "+m))
			}
			return jsonErrors
		}
	}

	_, err := valid.ValidateStruct(i)
	if err != nil {
		jsonErrors.Errors = append(jsonErrors.Errors, NewError("BindAndValidate.ValidateStructError",
			"ValidateStructError -> "+err.Error()))
		return jsonErrors
	}

	return nil
}
