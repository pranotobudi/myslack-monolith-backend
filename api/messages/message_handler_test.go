package messages

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/common"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getMessagesServiceFunc func(roomId string) ([]mongodb.Message, error)
)

type mockMessageService struct{}

func (m *mockMessageService) GetMessages(roomId string) ([]mongodb.Message, error) {
	return getMessagesServiceFunc(roomId)
}
func TestGetMessagesHandler(t *testing.T) {

	tt := []struct {
		Name     string
		mockFunc func(roomId string) ([]mongodb.Message, error)
		CodeWant int
	}{
		{
			Name: "GetMessages Success",
			mockFunc: func(roomId string) ([]mongodb.Message, error) {
				return []mongodb.Message{}, nil
			},
			CodeWant: http.StatusOK,
		},
		{
			Name: "GetMessages Failed",
			mockFunc: func(roomId string) ([]mongodb.Message, error) {
				return nil, errors.New("fail to get messages")
			},
			CodeWant: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getMessagesServiceFunc = tc.mockFunc

			messageHandler := NewMessageHandler()
			messageHandler.service = &mockMessageService{}
			rr := httptest.NewRecorder()
			// gin.SetMode(gin.ReleaseMode)
			c, _ := gin.CreateTestContext(rr)
			c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/messages?room_id=61f61d94fc663b6f4c8f3172", nil)
			// c.Request, _ = http.NewRequest(http.MethodGet, "", nil) // c.Params doesn't work for c.GetQuery("room_id")
			// c.Params = gin.Params{
			// 	{Key: "room_id", Value: "61f61d94fc663b6f4c8f3172"},
			// }

			log.Println(c.Params, c.Request.RequestURI)
			messageHandler.GetMessages(c)

			// check header StatusCode
			assert.EqualValues(t, tc.CodeWant, rr.Code)
			// check response (JSON format) StatusCode
			var response common.Response
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				assert.Errorf(t, err, "response format is not valid")
			}
			assert.EqualValues(t, tc.CodeWant, response.Meta.Code)
			// log.Println("test response: ", rc.Body.String())
		})
	}

}
