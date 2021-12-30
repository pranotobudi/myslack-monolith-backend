package messages

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMessages(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		roomId, ok := c.GetQuery("room_id")
		log.Println("GetMessages - roomId: ", roomId)
		if !ok {
			c.JSON(http.StatusBadRequest, roomId)
			return
		}
		filter := bson.M{"room_id": roomId}
		messages, err := mongo.GetMessages(filter)
		if err != nil {
			response := common.ResponseErrorFormatter(err)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		fmt.Println("inside room_io_handler-getMessages!: ", messages)
		response := common.ResponseFormatter(http.StatusOK, "success", "get messages successfull", messages)
		log.Println("RESPONSE TO BROWSER: ", response)
		// Add CORS headers, if no global CORS setting
		// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		c.JSON(http.StatusOK, response)
	}
}
