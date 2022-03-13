package emails

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
)

type IEmailHandler interface {
	MailChat(w http.ResponseWriter, r *http.Request)
}
type emailHandler struct {
	emailService IEmailService
}

// NewEmailHandler will initialize emailHandler object
func NewEmailHandler() *emailHandler {
	emailService := NewEmailService()
	return &emailHandler{emailService: emailService}
}

// UserMailChat will send email of chat with each room and its messages to the current user
func (h *emailHandler) MailChat(c *gin.Context) {
	// login
	var userMongo mongodb.User
	err := c.BindJSON(&userMongo)
	log.Println("UserMailChat userMongo: ", userMongo)

	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	msg, err := h.emailService.MailChat(userMongo)
	if err != nil {
		response := common.ResponseErrorFormatter(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := common.ResponseFormatter(http.StatusOK, "success", "email sent successfully", msg)
	log.Println("RESPONSE TO BROWSER: ", response)
	c.JSON(http.StatusOK, response)
}
