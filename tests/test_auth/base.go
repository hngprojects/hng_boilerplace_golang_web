package test_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/telex_be/pkg/controller/auth"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/tests"
)

func SetupAuthTestRouter() (*gin.Engine, *auth.Controller) {
	gin.SetMode(gin.TestMode)

	logger := tests.Setup()
	db := storage.Connection()
	validator := validator.New()

	authController := &auth.Controller{
		Db:        db,
		Validator: validator,
		Logger:    logger,
	}

	r := gin.Default()
	SetupAuthRoutes(r, authController)
	return r, authController
}

func SetupAuthRoutes(r *gin.Engine, authController *auth.Controller) {
	r.PUT("/api/v1/auth/change-password",
		middleware.Authorize(authController.Db.Postgresql),
		authController.ChangePassword)
	r.POST("/api/v1/auth/password-reset", authController.ResetPassword)
	r.POST("/api/v1/auth/password-reset/verify", authController.VerifyResetToken)
	r.POST("/api/v1/auth/magick-link", authController.RequestMagicLink)
	r.POST("/api/v1/auth/magick-link/verify", authController.VerifyMagicLink)
}
