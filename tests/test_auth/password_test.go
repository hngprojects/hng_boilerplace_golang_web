package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/controller/auth"
	"github.com/hngprojects/telex_be/tests"
	"github.com/hngprojects/telex_be/utility"
)

func TestUpdateUserPassword(t *testing.T) {
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

	loginData := models.LoginRequestModel{
		Email:    adminData.Email,
		Password: currUUID,
	}

	auth := auth.Controller{Db: authController.Db, Validator: authController.Validator, Logger: authController.Logger}
	token := tests.GetLoginToken(t, router, auth, loginData)

	t.Run("Successful Password Change", func(t *testing.T) {
		changePasswordRequest := models.ChangePasswordRequestModel{
			OldPassword: currUUID,
			NewPassword: currUUID + "nextest",
		}
		reqBody, _ := json.Marshal(changePasswordRequest)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/change-password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Password updated successfully")
	})

	t.Run("Incorrect Old Password", func(t *testing.T) {
		changePasswordRequest := models.ChangePasswordRequestModel{
			OldPassword: currUUID,
			NewPassword: currUUID + "nextest",
		}
		reqBody, _ := json.Marshal(changePasswordRequest)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/change-password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusBadRequest)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "old password is incorrect")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		changePasswordRequest := models.ChangePasswordRequestModel{
			OldPassword: currUUID,
			NewPassword: currUUID + "nextest",
		}
		reqBody, _ := json.Marshal(changePasswordRequest)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/change-password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Token could not be found!")
		tests.AssertResponseMessage(t, response["error"].(string), "Unauthorized")
	})

}
