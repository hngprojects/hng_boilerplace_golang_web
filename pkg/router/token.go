package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/pkg/controller/token"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/utility"
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
