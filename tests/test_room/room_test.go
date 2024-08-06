package test_room

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

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/room"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func TestRoomEndpoints(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	currUUID := utility.GenerateUUID()
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	roomController := room.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tst.SignupUser(t, r, auth, userSignUpData, false)

	token := tst.GetLoginToken(t, r, auth, loginData)

	createRoomReq := models.CreateRoomRequest{
		Name:        "Test Room",
		Description: "This is a test room",
		Username:    userSignUpData.UserName,
	}

	room_id := tst.CreateRoom(t, r, roomController, db, createRoomReq, token)

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Message      string
		Method       string
		Headers      map[string]string
		RequestURI   url.URL
	}{
		{
			Name: "Create Room Action",
			RequestBody: models.CreateRoomRequest{
				Name:        "Test Room",
				Description: "This is a test room",
				Username:    userSignUpData.UserName,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "room created successfully",
			Method:       http.MethodPost,
			RequestURI:   url.URL{Path: "/api/v1/rooms/"},
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Get Rooms Action",
			ExpectedCode: http.StatusOK,
			Message:      "rooms retrieved successfully",
			Method:       http.MethodGet,
			RequestURI:   url.URL{Path: "/api/v1/rooms/"},
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Get Room Action",
			ExpectedCode: http.StatusOK,
			Message:      "room retreived successfully",
			Method:       http.MethodGet,
			RequestURI:   url.URL{Path: fmt.Sprintf("/api/v1/rooms/%s", room_id)},
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Leave Room Action",
			ExpectedCode: http.StatusOK,
			Message:      "user left room successfully",
			Method:       http.MethodPost,
			RequestURI:   url.URL{Path: fmt.Sprintf("/api/v1/rooms/%s/leave", room_id)},
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		},
		{
			Name:         "Join Room Action",
			ExpectedCode: http.StatusOK,
			Message:      "room joined successfully",
			RequestBody:  models.JoinRoomRequest{
				Username: userSignUpData.UserName,
			},
			Method:     http.MethodPost,
			RequestURI: url.URL{Path: fmt.Sprintf("/api/v1/rooms/%s/join", room_id)},
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token,
			},
		}, 
	}

	room := room.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		roomUrl := r.Group(fmt.Sprintf("%v", "/api/v1/rooms"), middleware.Authorize(db.Postgresql))
		{
			roomUrl.GET("/", room.GetRooms)
			roomUrl.POST("/", room.CreateRoom)
			roomUrl.GET("/:roomId", room.GetRoom)
			roomUrl.POST("/:roomId/join", room.JoinRoom)
			roomUrl.POST("/:roomId/leave", room.LeaveRoom)
		}

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(test.Method, test.RequestURI.String(), &b)
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
