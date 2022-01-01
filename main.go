package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pranotobudi/myslack-monolith-backend/api/messages"
	"github.com/pranotobudi/myslack-monolith-backend/api/rooms"
	"github.com/pranotobudi/myslack-monolith-backend/api/users"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"github.com/pranotobudi/myslack-monolith-backend/msgserver"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		// production env provide that "APP_ENV" variable
		// executed in development only,
		//production set those on production environment settings

		// load local env variables to os
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("failed to load .env file")
		}
	}

	//mongoDB
	mongo := mongodb.NewMongoDB()
	// mongo.DataSeeder()

	// chat server
	// #1 init global message server as goroutine. this server will be an argument for each client
	hub := msgserver.NewHub()
	go hub.Run()

	// #2 init gin main server
	router := gin.Default()
	router.Use(CORS())

	// #3 handle url to init websocket client connection (will have func to handle incoming url)
	// this client will notify subscribe event to the global message server through channel.
	// through GetUserByEmail, we'll have email for user authentication

	// router.GET("/", serveMainPage)
	// router.Static("/static", "./static")
	router.GET("/rooms", rooms.GetRooms(mongo))
	router.GET("/messages", messages.GetMessages(mongo))
	router.POST("/room", rooms.AddRoom(mongo))
	router.GET("/room", rooms.GetAnyRoom(mongo))
	router.GET("/userByEmail", users.GetUserByEmail(mongo))
	router.POST("/userAuth", users.UserAuth(mongo))
	router.PUT("/updateUserRooms", users.UpdateUserRooms(mongo))
	router.GET("/websocket", msgserver.InitWebsocket(hub, mongo))

	// #4 run router server
	log.Println("server run on port:8080...")
	router.Run(":8080")
}

func serveMainPage(c *gin.Context) {
	// request: userId
	// response: user snapshot to load main page
	fmt.Println("inside serveMainPage!")
	c.File("static/index.html")
}
func serveStaticPage(c *gin.Context) {
	fmt.Println("inside serveStaticPage!")
	filePath := c.Request.URL.Path
	c.File(filePath)
}

func chatServer(c *gin.Context) {
	log.Println("inside chatServer! message Send..")
	c.Writer.WriteHeader(http.StatusAccepted)
	c.Writer.Write([]byte("msg send.."))
}

// func helloWorld(c *gin.Context) {
// 	firstname := c.DefaultQuery("firstname", "Guest")
// 	lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

// 	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
// }
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
