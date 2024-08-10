package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/telex_be/internal/config"
	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/internal/models/migrations"
	"github.com/hngprojects/telex_be/pkg/controller/auth"
	"github.com/hngprojects/telex_be/pkg/controller/room"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/pkg/repository/storage/postgresql"
	"github.com/hngprojects/telex_be/pkg/repository/storage/redis"
	"github.com/hngprojects/telex_be/utility"
)

func Setup() *utility.Logger {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../../app")

	postgresql.ConnectToDatabase(logger, config.TestDatabase)
	redis.ConnectToRedis(logger, config.Redis)
	db := storage.Connection()
	if config.TestDatabase.Migrate {
		migrations.RunAllMigrations(db)
	}
	return logger
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
func AssertBool(t *testing.T, got, expected bool) {
	if got != expected {
		t.Errorf("handler returned wrong boolean: got %v expected %v", got, expected)
	}
}

func AssertValidationError(t *testing.T, response map[string]interface{}, field string, expectedMessage string) {
	errors, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'error' field in response")
	}

	errorMessage, exists := errors[field]
	if !exists {
		t.Fatalf("expected validation error message for field '%s'", field)
	}

	if errorMessage != expectedMessage {
		t.Errorf("unexpected error message for field '%s': got %v, want %v", field, errorMessage, expectedMessage)
	}
}

func SignupUser(t *testing.T, r *gin.Engine, auth auth.Controller, userSignUpData models.CreateUserRequestModel, admin bool) {
	var (
		signupPath = "/api/v1/auth/register"
		signupURI  = url.URL{Path: signupPath}
	)

	r.POST(signupPath, auth.RegisterUser)

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(userSignUpData)
	req, err := http.NewRequest(http.MethodPost, signupURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
}

func GetLoginToken(t *testing.T, r *gin.Engine, auth auth.Controller, loginData models.LoginRequestModel) string {
	var (
		loginPath = "/api/v1/auth/login"
		loginURI  = url.URL{Path: loginPath}
	)
	r.POST(loginPath, auth.LoginUser)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(loginData)
	req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		return ""
	}

	data := ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	token := dataM["access_token"].(string)

	return token
}

func CreateRoom(t *testing.T, r *gin.Engine, room room.Controller, db *storage.Database, CreateData models.CreateRoomRequest, token string) (string, string) {
	var (
		createPath = "/api/v1/rooms/"
		createURI  = url.URL{Path: createPath}
	)

	roomUrl := r.Group(fmt.Sprintf("%v", "/api/v1/rooms"), middleware.Authorize(db.Postgresql))
	{
		roomUrl.POST("/", room.CreateRoom)
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(CreateData)
	req, err := http.NewRequest(http.MethodPost, createURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		return "", ""
	}

	data := ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	roomID := dataM["room_id"].(string)
	roomName := dataM["name"].(string)

	return roomID, roomName
}
