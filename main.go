package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

func main() {
	// gin setup
	router := gin.Default()
	router.Static("/static", "./static")

	// #1 init global message server as goroutine. this server will be an argument for each client

	// #2 handle url to init websocket client connection (will have func to handle incoming url)
	// each client have pointer to global message server
	// this client will notify subscribe event to the global message server through channel.
	// including (if needed: authentication, get client snapshot and subscribe it to message server)

	// #3 init gin main server
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/", serveMainPage)
	// router.GET("/static", serveStaticPage)
	router.GET("/chat", chatServer)
	router.Run(":8080")
}

func helloWorld(c *gin.Context) {
	firstname := c.DefaultQuery("firstname", "Guest")
	lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
}

func serveMainPage(c *gin.Context) {
	fmt.Println("inside serveMainPage!")
	c.File("static/index.html")
}
func serveStaticPage(c *gin.Context) {
	fmt.Println("inside serveStaticPage!")
	filePath := c.Request.URL.Path
	c.File(filePath)
}

func chatServer(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "")
	fmt.Println("inside chatServer! connection success")

	err = subscribe(c.Request.Context(), conn)
	if err != nil {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Printf("%v", err)
		return
	}

}

// func chatServer(w http.ResponseWriter, r *http.Request) {
// 	conn, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		return
// 	}
// 	defer conn.Close(websocket.StatusInternalError, "")

// 	err = subscribe(r.Context(), conn)
// 	if errors.Is(err, context.Canceled) {
// 		return
// 	}
// 	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
// 		websocket.CloseStatus(err) == websocket.StatusGoingAway {
// 		return
// 	}
// 	if err != nil {
// 		log.Printf("%v", err)
// 		return
// 	}
// }

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	s := &subscriber{
		msgs: make(chan []byte, 512),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
