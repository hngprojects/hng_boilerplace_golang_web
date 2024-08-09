package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/pkg/controller/room"
	"github.com/hngprojects/telex_be/pkg/middleware"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/utility"
)

func Room(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	room := room.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	roomUrl := r.Group(fmt.Sprintf("%v/rooms", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		roomUrl.POST("/", room.CreateRoom)
		roomUrl.POST("/:roomId/messages", room.AddRoomMsg)
		roomUrl.POST("/:roomId/join", room.JoinRoom)
		roomUrl.POST("/:roomId/leave", room.LeaveRoom)
		roomUrl.DELETE("/:roomId", room.DeleteRoom)
		roomUrl.PATCH("/:roomId/username", room.UpdateUsername)
		roomUrl.GET("/", room.GetRooms)
		roomUrl.GET("/:roomId", room.GetRoom)
		roomUrl.GET("/:roomId/messages", room.GetRoomMsg)
		roomUrl.GET("/:roomId/user-exist", room.CheckUser)
		roomUrl.GET("/name/:roomName", room.GetRoomByName)
		roomUrl.GET("/:roomId/num-users", room.CountRoomUsers)
		roomUrl.PATCH("/:roomId", room.UpdateRoom)
	}
	return r
}
