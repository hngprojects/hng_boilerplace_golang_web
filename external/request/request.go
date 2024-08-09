package request

import (
	"fmt"

	"github.com/hngprojects/telex_be/external/mocks"
	"github.com/hngprojects/telex_be/external/thirdparty/ipstack"
	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (
	// serializers
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"

	// requests
	IpstackResolveIp string = "ipstack_resolve_ip"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	var (
		config = config.GetConfig()
	)
	if !er.Test {
		switch name {
		case IpstackResolveIp:
			obj := ipstack.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v", config.IPStack.BaseUrl),
				Method:       "GET",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.IpstackResolveIp()
		default:
			return nil, fmt.Errorf("request not found")
		}

	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
