package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func (base *Controller) GoogleLogin(c *gin.Context) {

	var (
		req = models.GoogleRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed",
			utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := auth.CreateGoogleUser(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, respData)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user sign in successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user sign in successfully", respData)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ProviderLogout(c *gin.Context) {

	err := gothic.Logout(c.Writer, c.Request)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user logged out successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "user logged out successfully", nil)
	c.JSON(http.StatusOK, rd)

}
