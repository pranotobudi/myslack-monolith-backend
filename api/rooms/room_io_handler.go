package rooms

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRooms(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		rooms, err := mongo.GetRooms()
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
}

func GetAnyRoom(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// request: userId
		// response: user snapshot to load main page
		anyRoom := mongo.GetAnyRoom()
		fmt.Println("inside room_io_handler-getRoom anyRoom!: ", anyRoom)
		objID, err := primitive.ObjectIDFromHex(anyRoom.ID)
		if err != nil {
			panic(err)
		}

		filter := bson.M{"_id": objID}
		room := mongo.GetRoom(filter)
		fmt.Println("inside room_io_handler-getRoom!: ", room)
		response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", room)
		log.Println("RESPONSE TO BROWSER: ", response)
		// Add CORS headers
		// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		c.JSON(http.StatusOK, response)
	}
}

func AddRoom(mongo *mongodb.MongoDB) func(c *gin.Context) {

	return func(c *gin.Context) {
		var room mongodb.Room
		// c.Bind(&roomName)
		c.BindJSON(&room)
		log.Println("JSON roomName: ", room.Name)
		roomId, err := mongo.AddRoom(room.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, roomId)
			return
		}
		fmt.Println("room_io_handler-AddRoom: ", roomId)
		response := common.ResponseFormatter(http.StatusOK, "success", "add room successfull", roomId)
		c.JSON(http.StatusOK, response)
	}
}
