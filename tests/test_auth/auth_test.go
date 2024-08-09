package test_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/telex_be/internal/models"
	"github.com/hngprojects/telex_be/pkg/controller/auth"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	tst "github.com/hngprojects/telex_be/tests"
	"github.com/hngprojects/telex_be/utility"
)

// test user signup
func TestUserSignup(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/auth/users/signup"}
	currUUID := utility.GenerateUUID()

	tests := []struct {
		Name         string
		RequestBody  models.CreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful user register",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "verification email sent successfully",
		}, {
			Name: "details already exist",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "user already exists with the given email",
		}, {
			Name: "invalid email",
			RequestBody: models.CreateUserRequestModel{
				Email:       "emailtest",
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "email address is invalid",
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateUserRequestModel{
				PhoneNumber: "090909",
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/api/v1/auth/users/signup", auth.RegisterUser)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}

// test admin signup
func TestAdminSignup(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: "/api/v1/auth/admin/signup"}
	currUUID := utility.GenerateUUID()

	tests := []struct {
		Name         string
		RequestBody  models.CreateUserRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful user register",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testadmin%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "admin",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusCreated,
			Message:      "verification email sent successfully",
		}, {
			Name: "details already exist",
			RequestBody: models.CreateUserRequestModel{
				Email:       fmt.Sprintf("testadmin%v@qa.team", currUUID),
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "admin",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "user already exists with the given email",
		}, {
			Name: "invalid email",
			RequestBody: models.CreateUserRequestModel{
				Email:       "emailtest",
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
				FirstName:   "test",
				LastName:    "admin",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "email address is invalid",
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateUserRequestModel{
				PhoneNumber: "090909",
				FirstName:   "test",
				LastName:    "admin",
				Password:    "password",
				UserName:    fmt.Sprintf("test_username%v", currUUID),
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/api/v1/auth/admin/signup", auth.CreateAdmin)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}
		})

	}

}

// test login
func TestLogin(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	var (
		loginPath      = "/api/v1/auth/login"
		loginURI       = url.URL{Path: loginPath}
		currUUID       = utility.GenerateUUID()
		userSignUpData = models.CreateUserRequestModel{
			Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
			PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			FirstName:   "test",
			LastName:    "user",
			Password:    "password",
			UserName:    fmt.Sprintf("test_username%v", currUUID),
		}
	)

	tests := []struct {
		Name         string
		RequestBody  models.LoginRequestModel
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK email login successful",
			RequestBody: models.LoginRequestModel{
				Email:    userSignUpData.Email,
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "user login successfully",
		}, {
			Name:         "password not provided",
			RequestBody:  models.LoginRequestModel{},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "username or phone or email not provided",
			RequestBody: models.LoginRequestModel{
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "email does not exist",
			RequestBody: models.LoginRequestModel{
				Email:    utility.GenerateUUID(),
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "incorrect password",
			RequestBody: models.LoginRequestModel{
				Email:    fmt.Sprintf("testuser%v@qa.team", currUUID),
				Password: "incorrect",
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	r.POST(loginPath, auth.LoginUser)

	tst.SignupUser(t, r, auth, userSignUpData, false)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}

// test user logout
func TestLogout(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	var (
		logoutPath = "/api/v1/auth/logout"
		logoutURI  = url.URL{Path: logoutPath}
	)
	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "Hashira@password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	authen := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, authen, userSignUpData, false)

	token := tst.GetLoginToken(t, r, authen, loginData)

	tests := []struct {
		Name         string
		ExpectedCode int
		RequestBody  interface{}
		Message      string
		Headers      map[string]string
	}{
		{
			Name:         "user logout successfully",
			ExpectedCode: http.StatusOK,
			RequestBody:  nil,
			Message:      "user logout successfully",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, {
			Name:         "Session is invalid!",
			ExpectedCode: http.StatusUnauthorized,
			RequestBody:  nil,
			Message:      "Session is invalid!",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
	}

	authRoute := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	for _, test := range tests {
		r = gin.Default()

		authUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql))
		{
			authUrl.POST("/auth/logout", authRoute.LogoutUser)

		}
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, logoutURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
