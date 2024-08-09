package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/tests"
	"github.com/hngprojects/telex_be/utility"
)

func TestResetPassword(t *testing.T) {
	router, authController := SetupAuthTestRouter()
	db := authController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	someData, _ := utility.HashPassword(currUUID)

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe2",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: someData,
	}
	db.Create(&adminData)
	t.Run("Successful reset Password Request", func(t *testing.T) {
		forgotPasswordRequest := models.ForgotPasswordRequestModel{
			Email: adminData.Email,
		}
		reqBody, _ := json.Marshal(forgotPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Password reset email sent")
	})

	t.Run("Invalid Email reset Password Request", func(t *testing.T) {
		forgotPasswordRequest := models.ForgotPasswordRequestModel{
			Email: "nonexistent@qa.team",
		}
		reqBody, _ := json.Marshal(forgotPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})

	t.Run("Field Validation Error", func(t *testing.T) {
		forgotPasswordRequest := models.ForgotPasswordRequestModel{
			Email: "nonexistent@qa",
		}
		reqBody, _ := json.Marshal(forgotPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})
}

func TestVerifyResetPassword(t *testing.T) {
	router, authController := SetupAuthTestRouter()
	db := authController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	someData, _ := utility.HashPassword(currUUID)

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe2",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: someData,
	}
	db.Create(&adminData)

	resetToken := utility.GenerateUUID()
	expirationTime := time.Now().Add(30 * time.Minute)
	passwordResetData := models.PasswordReset{
		ID:        utility.GenerateUUID(),
		Email:     adminData.Email,
		Token:     resetToken,
		ExpiresAt: expirationTime,
	}
	db.Create(&passwordResetData)

	t.Run("Successful Password Reset", func(t *testing.T) {
		resetPasswordRequest := models.ResetPasswordRequestModel{
			Token:       resetToken,
			NewPassword: "newpassword",
		}
		reqBody, _ := json.Marshal(resetPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Password has been reset successfully")
	})

	t.Run("Invalid or Expired Token", func(t *testing.T) {
		resetPasswordRequest := models.ResetPasswordRequestModel{
			Token:       "invalidtoken",
			NewPassword: "newpassword",
		}
		reqBody, _ := json.Marshal(resetPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "invalid or expired token")
	})

	t.Run("Validation Error", func(t *testing.T) {
		resetPasswordRequest := models.ResetPasswordRequestModel{
			Token: resetToken,
		}
		reqBody, _ := json.Marshal(resetPasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})

	t.Run("New password length less than 7", func(t *testing.T) {
		changePasswordRequest := models.ChangePasswordRequestModel{
			OldPassword: "newpassword",
			NewPassword: "newpas",
		}
		reqBody, _ := json.Marshal(changePasswordRequest)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
		tests.AssertValidationError(t, response, "ResetPasswordRequestModel.NewPassword", "NewPassword must be at least 7 characters in length")
	})

}
