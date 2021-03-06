package messages

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
)

type IMessageHandler interface {
	GetMessages(c *gin.Context)
}
type messageHandler struct {
	service IMessageService
}

// NewMessageHandler initialize messageHandler object
func NewMessageHandler() *messageHandler {

	// func NewMessageHandler(messageService IMessageService) *messageHandler {
	messageService := NewMessageService()
	return &messageHandler{service: messageService}
}

// GetMessages will return list of messages for a room_id
func (h *messageHandler) GetMessages(c *gin.Context) {
	// TODO: we actually should use params instead of query, need changes in the frontend also
	// because room_id is unique resource (params), not filtering (query)
	// roomId := c.Param("room_id")
	// log.Println("GetMessages params: ", c.Params, "roomId: ", roomId)
	roomId, ok := c.GetQuery("room_id")
	log.Println("GetMessages - roomId: ", roomId, "Ok: ", ok)
	if !ok {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, errors.New("failed to get request Query"))
		c.JSON(http.StatusBadRequest, response)
		return
	}
	messages, err := h.service.GetMessages(roomId)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("inside room_io_handler-getMessages!: ", messages)
	response := common.ResponseFormatter(http.StatusOK, "success", "get messages successfull", messages)
	log.Println("RESPONSE TO BROWSER: ", response)

	c.JSON(http.StatusOK, response)
}
