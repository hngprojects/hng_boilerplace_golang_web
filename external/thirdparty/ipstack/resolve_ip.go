package ipstack

import (
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/external_models"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
)

func (r *RequestObj) IpstackResolveIp() (external_models.IPStackResolveIPResponse, error) {

	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPStackResolveIPResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipstack resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	path := "/" + ip + "?access_key=" + key

	logger.Info("ipstack resolve ip", ip)
	err := r.getNewSendRequestObject(nil, map[string]string{}, path).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("ipstack resolve ip", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
