package room

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func GetRooms(db *gorm.DB) ([]models.Room, int, error) {
	var room models.Room

	rooms, err := room.GetRooms(db)
	if err != nil {
		return rooms, http.StatusInternalServerError, err
	}
	return rooms, http.StatusOK, nil
}

func CreateRoom(req models.CreateRoomRequest, db *gorm.DB, userId string) (models.Room, int, error) {
	var joinRoomReq models.JoinRoomRequest

	room := models.Room{
		ID:          utility.GenerateUUID(),
		Name:        req.Name,
		OwnerId:     userId,
		Description: req.Description,
	}

	joinRoomReq.RoomID = room.ID
	joinRoomReq.UserID = userId
	joinRoomReq.Username = req.Username

	err := room.CreateRoom(db)
	if err != nil {
		return room, http.StatusBadRequest, err
	}

	err = room.AddUserToRoom(db, joinRoomReq)
	if err != nil {
		return room, http.StatusBadRequest, err
	}
	return room, http.StatusOK, nil
}

func GetRoom(db *gorm.DB, roomID string) ([]models.UserRoom, int, error) {
	var room models.Room

	fetchedUsers, err := room.GetRoomUsersByID(db, roomID)
	if err != nil {
		return fetchedUsers, http.StatusBadRequest, err
	}
	return fetchedUsers, http.StatusOK, nil
}

func GetRoomByName(db *gorm.DB, name string) (models.Room, int, error) {
	var r models.Room

	room, err := r.GetRoomByName(db, name)
	if err != nil {
		return room, http.StatusBadRequest, err
	}
	return room, http.StatusOK, nil
}

func GetRoomMsg(roomId, userID string, db *gorm.DB) ([]models.Message, int, error) {
	var message models.Message

	resp, err := message.GetMessagesByRoomID(db, userID, roomId)

	if err != nil {
		return []models.Message{}, http.StatusBadRequest, err
	}

	return resp, http.StatusOK, nil

}

func JoinRoom(db *gorm.DB, req models.JoinRoomRequest) (int, error) {
	var room models.Room

	err := room.AddUserToRoom(db, req)

	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func LeaveRoom(db *gorm.DB, room_id, user_id string) (int, error) {
	var room models.Room

	_, _, err := GetRoom(db, room_id)
	if err != nil {
		return http.StatusBadRequest, errors.New("room does not exist")
	}

	err = room.RemoveUserFromRoom(db, room_id, user_id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil

}

func AddRoomMsg(req models.CreateMessageRequest, db *gorm.DB) (int, error) {

	message := models.Message{
		Content: req.Content,
		RoomID:  req.RoomId,
		UserID:  req.UserId,
	}

	err := message.CreateMessage(db)

	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusCreated, nil
}

func UpdateUsername(req models.UpdateRoomUserNameReq, db *gorm.DB, roomId, userId string) (int, error) {

	var userroom models.UserRoom

	err := userroom.UpdateUsername(db, req, roomId, userId)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func DeleteRoom(db *gorm.DB, roomId, userId string) (int, error) {
	var room models.Room

	room, err := room.GetRoomByID(db, roomId)

	fmt.Println(room.OwnerId, userId)

	if room.OwnerId != userId {
		return http.StatusUnauthorized, errors.New("user not authorized")
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = room.Delete(db)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func CountRoomUsers(db *gorm.DB, roomId string) (int, int, error) {
	var userRoom models.UserRoom

	count, err := userRoom.CountRoomUsers(db, roomId)
	if err != nil {
		return count, http.StatusBadRequest, err
	}
	return count, http.StatusOK, nil
}

func UpdateRoom(db *gorm.DB, req models.UpdateRoomRequest, roomId string) (models.Room, error) {
	var (
		room models.Room
	)
	updatedRoom, _, err := room.UpdateRoom(db, req, roomId)
	if err != nil {
		return updatedRoom, err
	}
	return updatedRoom, nil
}

func CheckUser(roomId, userID string, db *gorm.DB) (gin.H, int, error) {
	var (
		userroom models.UserRoom
		resp     gin.H
	)

	status, chk := userroom.CheckUser(db, userID, roomId)

	resp = gin.H{
		"exist": status,
		"msg":   chk,
	}

	return resp, http.StatusOK, nil

}
