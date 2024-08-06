package ipstack_mocks

import (
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/external_models"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/config"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func IpstackResolveIp(logger *utility.Logger, idata interface{}) (external_models.IPStackResolveIPResponse, error) {

	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPStackResolveIPResponse
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipstack resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}
	outBoundResponse.Ip = ip
	outBoundResponse.City = "city"
	outBoundResponse.CountryName = "name"

	path := "/" + ip + "?access_key=" + key

	logger.Info("ipstack resolve ip", ip, path)

	return outBoundResponse, nil
}
