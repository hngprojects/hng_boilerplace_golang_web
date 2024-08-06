package mocks

import (
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/mocks/ipstack_mocks"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type ExternalRequest struct {
	Logger     *utility.Logger
	Test       bool
	RequestObj RequestObj
}

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod string
	Logger       *utility.Logger
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	switch name {
	case "ipstack_resolve_ip":
		return ipstack_mocks.IpstackResolveIp(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
}
