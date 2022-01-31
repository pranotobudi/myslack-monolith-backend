package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
)

type IUserHandler interface {
	GetUserByEmail(c *gin.Context)
	UserAuth(c *gin.Context)
	UpdateUserRooms(c *gin.Context)
	HelloWorld(c *gin.Context)
}
type userHandler struct {
	userService IUserService
}

func NewUserHandler() *userHandler {
	userService := NewUserService()
	return &userHandler{userService: userService}
}

func (h *userHandler) GetUserByEmail(c *gin.Context) {
	email, ok := c.GetQuery("email")
	log.Println("GetUserByEmail - email: ", email)
	if !ok {
		c.JSON(http.StatusBadRequest, email)
		return
	}
	// filter := bson.M{"email": email}
	// userPtr, err := mongo.GetUser(filter)
	userPtr, err := h.userService.GetUser(email)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
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

func (h *userHandler) UserAuth(c *gin.Context) {
	// return func(c *gin.Context) {
	// login
	var userAuth mongodb.UserAuth

	err := c.BindJSON(&userAuth)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	log.Println("GetUserByEmail - email: ", userAuth.Email)
	userPtr, err := h.userService.UserAuth(userAuth)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// filter := bson.M{"email": userAuth.Email}
	// userPtr, err := mongo.GetUser(filter)
	// if err != nil {
	// 	// c.JSON(http.StatusInternalServerError, err)
	// 	// return
	// 	log.Println("inside room_io_handler-UserAuth error: ", err)
	// }

	// if userPtr != nil {
	// 	// two possibility:
	// 	// empty User{}, means: user found but failed to decode
	// 	// non-empty User, means: user found and success to decode
	// 	// for both result we'll return it anyway

	// 	fmt.Println("inside room_io_handler-UserAuth user found!: ", *userPtr)
	// 	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	// 	log.Println("RESPONSE TO BROWSER: ", response)
	// 	// Add CORS headers, if no global CORS setting
	// 	// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// 	// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	// 	c.JSON(http.StatusOK, response)
	// 	return
	// }

	// // user == nil, user not found
	// // register
	// userDoc := bson.D{{"email", userAuth.Email}, {"username", ""}, {"user_image", userAuth.UserImage}, {"rooms", bson.A{}}}
	// userID, err := mongo.AddUser(userDoc)

	// // return User data as response
	// objID, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	// filter = bson.M{"_id": objID}
	// userPtr, err = mongo.GetUser(filter)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	fmt.Println("inside room_io_handler-UserAuth user registered! ID: ", *userPtr)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUserRooms(c *gin.Context) {
	// login
	var userMongo mongodb.User

	err := c.BindJSON(&userMongo)
	log.Println("UpdateUserRooms userMongo: ", userMongo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userPtr, err := h.userService.UpdateUserRooms(userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", *userPtr)
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)

}

func HelloWorld(c *gin.Context) {
	// firstname := c.DefaultQuery("firstname", "Guest")
	// lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	// c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	response := common.ResponseFormatter(http.StatusOK, "success", "get user successfull", "Hello from MySlack App. A Persistence Chat App... The server is running at the background..")
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)
}
