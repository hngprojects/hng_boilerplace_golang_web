package room

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"

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

func (base *Controller) UpdateUsername(c *gin.Context) {
	var req models.UpdateRoomUserNameReq

	roomId := c.Param("roomId")

	if _, err := uuid.Parse(roomId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

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

	code, err := room.UpdateUsername(req, base.Db.Postgresql, roomId, userId)
	if err != nil {
		base.Logger.Info("error creating room")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("username updated successfully")
	rd := utility.BuildSuccessResponse(code, "username updated successfully", nil)
	c.JSON(code, rd)
}

func (base *Controller) DeleteRoom(c *gin.Context) {

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

	code, err := room.DeleteRoom(base.Db.Postgresql, RoomId, UserId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("room deleted successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "room deleted successfully", nil)
	c.JSON(code, rd)
}

func (base *Controller) GetRoomByName(c *gin.Context) {
	name := c.Params.ByName("roomName")

	respData, code, err := room.GetRoomByName(base.Db.Postgresql,name)
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

func (base *Controller) CountRoomUsers(c *gin.Context) {
	roomId := c.Param("roomId")

	if _, err := uuid.Parse(roomId); err != nil {
		base.Logger.Info("failed to get roomId")
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid room id format", errors.New("failed to parse room id"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	totalCount, code, err := room.CountRoomUsers(base.Db.Postgresql, roomId)
	if err != nil {
		base.Logger.Info("error getting total room users")
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("room users count retrieved successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "room users count retrieved successfully", totalCount)
	c.JSON(code, rd)
}

func (base *Controller) UpdateRoom(c *gin.Context) {
	id := c.Param("room_id")
	var req models.UpdateRoomRequest

	if _, err := uuid.Parse(id); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid ID format", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Invalid request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validator.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	result, err := room.UpdateRoom(base.Db.Postgresql, req ,id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rd := utility.BuildErrorResponse(http.StatusNotFound, "error", "Room not found", err, nil)
			c.JSON(http.StatusNotFound, rd)
		} else {
			rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to update room", err, nil)
			c.JSON(http.StatusInternalServerError, rd)
		}
		return
	}

	base.Logger.Info("Room updated successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Room updated successfully", result)
	c.JSON(http.StatusOK, rd)
}
