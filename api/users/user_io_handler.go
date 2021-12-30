package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByEmail(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		email, ok := c.GetQuery("email")
		log.Println("GetUserByEmail - email: ", email)
		if !ok {
			c.JSON(http.StatusBadRequest, email)
			return
		}
		filter := bson.M{"email": email}
		user, err := mongo.GetUser(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, user)
			return
		}
		fmt.Println("inside room_io_handler-getRoom GetUserByEmail!: ", user)
		response := common.ResponseFormatter(http.StatusOK, "success", "get rooms successfull", user)
		log.Println("RESPONSE TO BROWSER: ", response)
		// Add CORS headers, if no global CORS setting
		// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		c.JSON(http.StatusOK, response)
	}
}
