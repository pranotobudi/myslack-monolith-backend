package rooms

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
)

type IRoomHandler interface {
	GetRooms(c *gin.Context)
	GetAnyRoom(c *gin.Context)
	AddRoom(c *gin.Context)
}

type roomHandler struct {
	roomService IRoomService
}

func NewRoomHandler() *roomHandler {
	roomService := NewRoomService()
	return &roomHandler{roomService: roomService}
}

func (h *roomHandler) GetRooms(c *gin.Context) {
	// rooms, err := mongo.GetRooms()
	rooms, err := h.roomService.GetRooms()
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getRooms!: ", rooms)
	response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", rooms)
	log.Println("RESPONSE TO BROWSER: ", response)
	// Add CORS headers, if no global CORS setting
	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, response)
}

func (h *roomHandler) GetAnyRoom(c *gin.Context) {
	// request: userId
	// response: user snapshot to load main page
	room, err := h.roomService.GetAnyRoom()
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", room)
	log.Println("RESPONSE TO BROWSER: ", response)
	// Add CORS headers
	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, response)
}

func (h *roomHandler) AddRoom(c *gin.Context) {

	var room mongodb.Room
	// c.Bind(&roomName)

	c.BindJSON(&room)
	log.Println("JSON roomName: ", room.Name)
	// roomId, err := mongo.AddRoom(room.Name)
	roomId, err := h.roomService.AddRoom(room.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, roomId)
		return
	}
	fmt.Println("room_io_handler-AddRoom: ", roomId)
	response := common.ResponseFormatter(http.StatusOK, "success", "add room successfull", roomId)
	c.JSON(http.StatusOK, response)
}
