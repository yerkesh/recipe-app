package writer

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"recipe-app/pkg/domain/constant"
	"recipe-app/pkg/util/fault"
)

type ServiceResponse struct {
	Code        string            `json:"code,omitempty"`
	Status      string            `json:"status,omitempty"`
	Message     string            `json:"message,omitempty"`
	Debug       string            `json:"debug,omitempty"`
	MessageCode string            `json:"message_code,omitempty"`
	Validation  map[string]string `json:"validation,omitempty"`
}

type EmptyJSON struct{}

func ServiceResponseOk(msg string) ServiceResponse {
	return ServiceResponse{ //nolint:exhaustivestruct // partial response
		Code:    strconv.Itoa(http.StatusOK),
		Status:  http.StatusText(http.StatusOK),
		Message: msg,
	}
}

func ServiceResponseCreated(msg string) ServiceResponse {
	return ServiceResponse{ //nolint:exhaustivestruct // partial response
		Code:    strconv.Itoa(http.StatusCreated),
		Status:  http.StatusText(http.StatusCreated),
		Message: msg,
	}
}

func UnhandledServiceResponse(e error) *ServiceResponse {
	return &ServiceResponse{ //nolint:exhaustivestruct // partial response
		Code:    strconv.Itoa(http.StatusInternalServerError),
		Status:  http.StatusText(http.StatusInternalServerError),
		Debug:   e.Error(),
		Message: constant.MsgUnhandledErr,
	}
}

func ErrServiceResponse(e *fault.WhsError) *ServiceResponse {
	return &ServiceResponse{ //nolint:exhaustivestruct // partial response
		Code:       strconv.Itoa(e.HTTPStatus),
		Status:     http.StatusText(e.HTTPStatus),
		Debug:      e.Debug,
		Message:    e.Message,
		Validation: e.Validation,
	}
}

func checkSetSsoError(err error) (status int, value *ServiceResponse) {
	if e, ok := err.(*fault.WhsError); ok { //nolint:errorlint // don't know how to use
		status = e.HTTPStatus
		value = ErrServiceResponse(e)
	} else {
		status = http.StatusInternalServerError
		value = UnhandledServiceResponse(err)
	}

	return status, value
}

func okServiceResponse(body interface{}) (status int, value interface{}) {
	status = http.StatusOK

	if body != nil {
		value = body
	} else {
		value = &ServiceResponse{ //nolint:exhaustivestruct // partial response
			Code:    strconv.Itoa(http.StatusOK),
			Status:  http.StatusText(http.StatusOK),
			Message: constant.MsgCreated,
		}
	}

	return status, value
}

func HTTPResponseWriter(resp http.ResponseWriter, err error, body interface{}) {
	var (
		status int
		value  interface{}
	)

	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		status, value = checkSetSsoError(err)
	} else {
		status, value = okServiceResponse(body)
	}

	resp.WriteHeader(status)

	err = json.NewEncoder(resp).Encode(value)

	if err != nil {
		log.Printf("couldn't encode the value[%v] into json: %v", value, err)
	}
}
