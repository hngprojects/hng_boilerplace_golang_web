package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/telex_be/internal/models"
	service "github.com/hngprojects/telex_be/services/auth"
	"github.com/hngprojects/telex_be/utility"
)

func (base *Controller) RequestMagicLink(c *gin.Context) {
	var (
		req = models.MagicLinkRequest{}
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

	url := c.Request.Header.Get("Referer")

	respData, code, err := service.MagicLinkRequest(req.Email, url, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("magic link sent to email")

	rd := utility.BuildSuccessResponse(http.StatusOK, "Magic link sent to email", respData)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) VerifyMagicLink(c *gin.Context) {
	var (
		req = models.VerifyMagicLinkRequest{}
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

	respData, code, err := service.VerifyMagicLinkToken(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("user login successfully")

	rd := utility.BuildSuccessResponse(http.StatusOK, "User login successfully", respData)
	c.JSON(http.StatusOK, rd)

}
