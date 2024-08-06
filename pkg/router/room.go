package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/room"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Room(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	room := room.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	roomUrl := r.Group(fmt.Sprintf("%v/rooms", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		roomUrl.GET("/", room.GetRooms)
		roomUrl.POST("/", room.CreateRoom)
		roomUrl.GET("/:roomId", room.GetRoom)
		roomUrl.GET("/:roomId/messages", room.GetRoomMsg)
		roomUrl.POST("/:roomId/messages", room.AddRoomMsg)
		roomUrl.POST("/:roomId/join", room.JoinRoom)
		roomUrl.POST("/:roomId/leave", room.LeaveRoom)
	}
	return r
}
