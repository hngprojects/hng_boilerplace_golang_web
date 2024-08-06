package utility

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Response struct {
	Status     string      `json:"status,omitempty"`
	Code       int         `json:"status_code,omitempty"`
	Name       string      `json:"name,omitempty"` //name of the error
	Message    string      `json:"message,omitempty"`
	Error      interface{} `json:"error,omitempty"` //for errors that occur even if request is successful
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Extra      interface{} `json:"extra,omitempty"`
}

// BuildResponse method is to inject data value to dynamic success response
func BuildSuccessResponse(code int, message string, data interface{}, pagination ...interface{}) Response {
	res := ResponseMessage(code, "success", "", message, nil, data, pagination, nil)
	return res
}

// BuildErrorResponse method is to inject data value to dynamic failed response
func BuildErrorResponse(code int, status string, message string, err interface{}, data interface{}, logger ...bool) Response {
	res := ResponseMessage(code, status, "", message, err, data, nil, nil)
	return res
}

// ResponseMessage method for the central response holder
func ResponseMessage(code int, status string, name string, message string, err interface{}, data interface{}, pagination interface{}, extra interface{}) Response {
	if pagination != nil && reflect.ValueOf(pagination).IsNil() {
		pagination = nil
	}

	if code == http.StatusInternalServerError {
		fmt.Println("internal server error", message, err, data)
		message = "internal server error"
		err = message
	}

	res := Response{
		Code:       code,
		Name:       name,
		Status:     status,
		Message:    message,
		Error:      err,
		Data:       data,
		Pagination: pagination,
		Extra:      extra,
	}
	return res
}

func UnauthorisedResponse(code int, status string, name string, message string) Response {
	res := ResponseMessage(http.StatusUnauthorized, status, name, message, nil, nil, nil, nil)
	return res
}

func ValidationResponse(err error, validate *validator.Validate) validator.ValidationErrorsTranslations {
	errs := err.(validator.ValidationErrors)
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	return errs.Translate(trans)
}
