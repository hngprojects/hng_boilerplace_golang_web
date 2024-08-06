package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/elliotchance/phpserialize"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type SendRequestObject struct {
	Name         string
	Logger       *utility.Logger
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	Data         interface{}
	DecodeMethod string
	UrlPrefix    string
}

func GetNewSendRequestObject(logger *utility.Logger, name, path, method, urlPrefix, decodeMethod string, headers map[string]string, successCode int, data interface{}) *SendRequestObject {
	return &SendRequestObject{
		Logger:       logger,
		Name:         name,
		Path:         path,
		Method:       method,
		UrlPrefix:    urlPrefix,
		DecodeMethod: decodeMethod,
		Headers:      headers,
		SuccessCode:  successCode,
		Data:         data,
	}
}

var (
	ResponseCode int
	ResponseBody string
)

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (r *SendRequestObject) SendRequest(response interface{}) error {
	var (
		data   = r.Data
		logger = r.Logger
		name   = r.Name
		err    error
	)

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		logger.Error("encoding error", name, err.Error())
	}

	logger.Info("before prefix", name, r.Path, data, buf)
	if r.UrlPrefix != "" {
		r.Path += r.UrlPrefix
	}
	logger.Info("after prefix", name, r.Path, data, buf)

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.Path, buf)
	if err != nil {
		logger.Error("request creation error", name, err.Error())
		return err
	}

	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	logger.Info("request", name, r.Path, r.Method, r.Headers)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("client do", name, err.Error())
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("readin body error", name, err.Error())
		return err
	}

	if r.DecodeMethod != PhpSerializerMethod {
		err = json.Unmarshal(body, response)
		if err != nil {
			logger.Error("json decoding error", name, err.Error())
			return err
		}
	}

	logger.Info("response body", name, r.Path, string(body))

	if r.DecodeMethod == PhpSerializerMethod {
		err := phpserialize.Unmarshal(body, response)
		if err != nil {
			logger.Error("php serializer decoding error", name, err.Error())
			return err
		}
	}

	defer res.Body.Close()
	ResponseCode = res.StatusCode

	if res.StatusCode == r.SuccessCode {
		return nil
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("external requests error for request %v, code %v", name, strconv.Itoa(res.StatusCode))
	}

	return nil
}
