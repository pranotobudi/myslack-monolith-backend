package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		userPtr, err := mongo.GetUser(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		fmt.Println("inside room_io_handler-getRoom GetUserByEmail!: ", *userPtr)
		response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
		log.Println("RESPONSE TO BROWSER: ", response)
		// Add CORS headers, if no global CORS setting
		// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		c.JSON(http.StatusOK, response)
	}
}

func UserAuth(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// login
		var userAuth mongodb.UserAuth

		err := c.BindJSON(&userAuth)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		log.Println("GetUserByEmail - email: ", userAuth.Email)

		filter := bson.M{"email": userAuth.Email}
		userPtr, err := mongo.GetUser(filter)
		if err != nil {
			// c.JSON(http.StatusInternalServerError, err)
			// return
			log.Println("inside room_io_handler-UserAuth error: ", err)
		}

		if userPtr != nil {
			// two possibility:
			// empty User{}, means: user found but failed to decode
			// non-empty User, means: user found and success to decode
			// for both result we'll return it anyway

			fmt.Println("inside room_io_handler-UserAuth user found!: ", *userPtr)
			response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
			log.Println("RESPONSE TO BROWSER: ", response)
			// Add CORS headers, if no global CORS setting
			// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

			c.JSON(http.StatusOK, response)
			return
		}

		// user == nil, user not found
		// register
		userDoc := bson.D{{"email", userAuth.Email}, {"username", ""}, {"user_image", userAuth.UserImage}, {"rooms", bson.A{}}}
		userID, err := mongo.AddUser(userDoc)

		// return User data as response
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		filter = bson.M{"_id": objID}
		userPtr, err = mongo.GetUser(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Println("inside room_io_handler-UserAuth user registered! ID: ", *userPtr)
		response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
		log.Println("RESPONSE TO BROWSER: ", response)
		c.JSON(http.StatusOK, response)

	}
}

func UpdateUserRooms(mongo *mongodb.MongoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// login
		var userMongo mongodb.User

		err := c.BindJSON(&userMongo)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// log.Println("UpdateUser - userMongo: ", userMongo)
		// log.Println("UpdateUser - userMongo email: ", userMongo.Email)
		log.Println("UpdateUser - userMongo rooms: ", userMongo.Rooms)
		log.Println("UpdateUser - userMongo ID: ", userMongo.ID)

		// filter := bson.D{{"email", userMongo.Email}}
		// update := bson.M{"rooms": userMongo.Rooms}

		filter := bson.M{"email": userMongo.Email}
		opts := options.Update().SetUpsert(true)

		// remove all user rooms first
		update := bson.D{{"$set", bson.M{"rooms": []string{}}}}
		err = mongo.UpdateUser(filter, update, opts)
		if err != nil {
			log.Println("inside user_io_handler-UpdateUser remove all user rooms, error: ", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		// add room one by one
		for _, room := range userMongo.Rooms {
			update := bson.D{{"$push", bson.M{"rooms": room}}}
			err = mongo.UpdateUser(filter, update, opts)
			if err != nil {
				log.Println("inside user_io_handler-UpdateUser add room, error: ", err)
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}

		// get user with updated rooms element
		userPtr, err := mongo.GetUser(filter)
		if err != nil {
			log.Println("failed to GetUser")
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Println("inside room_io_handler-UserAuth user registered! ID: ", *userPtr)
		response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
		log.Println("RESPONSE TO BROWSER: ", response)
		c.JSON(http.StatusOK, response)

	}
}
