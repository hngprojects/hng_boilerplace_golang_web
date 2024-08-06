package room

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/room"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest

	claims, exists := c.Get("userClaims")
	if !exists {
		base.Logger.Info("error getting claims")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error getting claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)
	userId := userClaims["user_id"].(string)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		base.Logger.Info("error parsing request body")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		base.Logger.Info("validation failed")
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := room.CreateRoom(req, base.Db.Postgresql, userId)
	if err != nil {
		base.Logger.Info("error creating room")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("room created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "room created successfully", respData)
	c.JSON(http.StatusCreated, rd)
}

func (base *Controller) GetRooms(c *gin.Context) {

	respData, code, err := room.GetRooms(base.Db.Postgresql)
	if err != nil {
		base.Logger.Info("error getting rooms")
		rd := utility.BuildErrorResponse(code, "error",
			err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("rooms retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "rooms retrieved successfully", respData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetRoom(c *gin.Context) {
	room_id := c.Param("roomId")

	if _, err := uuid.Parse(room_id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := room.GetRoom(base.Db.Postgresql, room_id)
	if err != nil {
		base.Logger.Info("error getting room")
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("room retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "room retreived successfully", respData)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetRoomMsg(c *gin.Context) {

	RoomId := c.Param("roomId")

	if _, err := uuid.Parse(RoomId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	claims, exists := c.Get("userClaims")

	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", errors.New("user not authorized"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)

	UserId := userClaims["user_id"].(string)

	respData, code, err := room.GetRoomMsg(RoomId, UserId, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("room messages fetched successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "room messages fetched successfully", respData)
	c.JSON(code, rd)
}

func (base *Controller) AddRoomMsg(c *gin.Context) {
	var (
		req models.CreateMessageRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	req.RoomId = c.Param("roomId")

	if _, err := uuid.Parse(req.RoomId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")

	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", errors.New("user not authorized"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)

	req.UserId = userClaims["user_id"].(string)

	code, err := room.AddRoomMsg(req, base.Db.Postgresql)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("message added successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "message added successfully", gin.H{})
	c.JSON(code, rd)
}

func (base *Controller) JoinRoom(c *gin.Context) {
	var (
		req models.JoinRoomRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		base.Logger.Info("validation failed")
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	room_id := c.Param("roomId")

	claims, exists := c.Get("userClaims")
	if !exists {
		base.Logger.Info("error getting claims")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error getting claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)

	user_id := userClaims["user_id"].(string)

	req.RoomID = room_id
	req.UserID = user_id

	code, err := room.JoinRoom(base.Db.Postgresql, req)
	if err != nil {
		base.Logger.Info("error joining room")
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("room joined successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "room joined successfully", nil)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) LeaveRoom(c *gin.Context) {

	roomId := c.Param("roomId")

	if _, err := uuid.Parse(roomId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	claims, exists := c.Get("userClaims")

	if !exists {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", errors.New("user not authorized"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	userClaims := claims.(jwt.MapClaims)

	user_id := userClaims["user_id"].(string)

	code, err := room.LeaveRoom(base.Db.Postgresql, roomId, user_id)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user left room successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "user left room successfully", gin.H{})
	c.JSON(code, rd)
}
