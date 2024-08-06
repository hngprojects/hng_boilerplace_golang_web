package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/token"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TokenGen(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	token := token.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	tokenUrl := r.Group(fmt.Sprintf("%v/token", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		tokenUrl.GET("/connection", token.GetConnToken)
		tokenUrl.POST("/subscription", token.GetSubToken)
	}
	return r
}
