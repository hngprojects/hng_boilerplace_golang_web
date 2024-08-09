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

func TestRequestMagicLink(t *testing.T) {
	router, authController := SetupAuthTestRouter()
	db := authController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe2",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: password,
	}
	db.Create(&adminData)

	t.Run("Successful Magic Link Request", func(t *testing.T) {
		requestMagicLink := models.MagicLinkRequest{
			Email: adminData.Email,
		}
		reqBody, _ := json.Marshal(requestMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Magic link sent to email")
	})

	t.Run("Invalid Email for Magic Link Request", func(t *testing.T) {
		requestMagicLink := models.MagicLinkRequest{
			Email: "nonexistent@qa.team",
		}
		reqBody, _ := json.Marshal(requestMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusNotFound)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "user not found")
	})

	t.Run("Field Validation Error", func(t *testing.T) {
		requestMagicLink := models.MagicLinkRequest{
			Email: "invalid-email",
		}
		reqBody, _ := json.Marshal(requestMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})
}

func TestVerifyMagicLink(t *testing.T) {
	router, authController := SetupAuthTestRouter()
	db := authController.Db.Postgresql
	currUUID := utility.GenerateUUID()
	password, _ := utility.HashPassword("password")

	adminData := models.User{
		ID:       utility.GenerateUUID(),
		Name:     "admin jane doe2",
		Email:    fmt.Sprintf("testadmin%v@qa.team", currUUID),
		Password: password,
	}
	db.Create(&adminData)

	magicToken := utility.GenerateUUID()
	expirationTime := time.Now().Add(10 * time.Minute)
	magicLinkData := models.MagicLink{
		ID:        utility.GenerateUUID(),
		Email:     adminData.Email,
		Token:     magicToken,
		ExpiresAt: expirationTime,
	}
	db.Create(&magicLinkData)

	t.Run("Successful Magic Link Verification", func(t *testing.T) {
		verifyMagicLink := models.VerifyMagicLinkRequest{
			Token: magicToken,
		}
		reqBody, _ := json.Marshal(verifyMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusOK)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "User login successfully")
	})

	t.Run("Invalid or Expired Token", func(t *testing.T) {
		verifyMagicLink := models.VerifyMagicLinkRequest{
			Token: "invalidtoken",
		}
		reqBody, _ := json.Marshal(verifyMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnauthorized)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "invalid or expired token")
	})

	t.Run("Validation Error", func(t *testing.T) {
		verifyMagicLink := models.VerifyMagicLinkRequest{
			Token: "",
		}
		reqBody, _ := json.Marshal(verifyMagicLink)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/magick-link/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		tests.AssertStatusCode(t, resp.Code, http.StatusUnprocessableEntity)
		response := tests.ParseResponse(resp)
		tests.AssertResponseMessage(t, response["message"].(string), "Validation failed")
	})
}
