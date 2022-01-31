package messages

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pranotobudi/myslack-monolith-backend/mongodb"
	"github.com/stretchr/testify/assert"
)

var (
	getMessagesFunc func(filter interface{}) ([]mongodb.Message, error)
)

type mockService struct{}

func (m *mockService) GetMessages(filter interface{}) ([]mongodb.Message, error) {
	return getMessagesFunc(filter)
}
func TestGetMessages(t *testing.T) {

	tt := []struct {
		Name     string
		mockFunc func(filter interface{}) ([]mongodb.Message, error)
		CodeWant int
	}{
		{
			Name: "GetMessages Success",
			mockFunc: func(filter interface{}) ([]mongodb.Message, error) {
				return []mongodb.Message{}, nil
			},
			CodeWant: http.StatusOK,
		},
		{
			Name: "GetMessages Failed",
			mockFunc: func(filter interface{}) ([]mongodb.Message, error) {
				return nil, errors.New("fail to get messages")
			},
			CodeWant: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			getMessagesFunc = tc.mockFunc

			// messageHandler := NewMessageHandler(&mockService{})
			messageHandler := NewMessageHandler()
			messageHandler.service = &mockService{}
			rc := httptest.NewRecorder()
			// gin.SetMode(gin.ReleaseMode)
			c, _ := gin.CreateTestContext(rc)
			// c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/messages?room_id=61f61d94fc663b6f4c8f3172", nil)
			c.Request, _ = http.NewRequest(http.MethodGet, "", nil) // c.Params doesn't work for c.GetQuery("room_id")
			c.Params = gin.Params{
				{Key: "room_id", Value: "61f61d94fc663b6f4c8f3172"},
			}
			log.Println(c.Params, c.Request.RequestURI)
			messageHandler.GetMessages(c)

			assert.EqualValues(t, tc.CodeWant, rc.Code)
			log.Println("test response: ", rc.Body.String())
			// var response common.Response
			// err := json.Unmarshal(rc.Body.Bytes(), &response)
			// assert.Nil(t, err)
			// assert.
		})
	}

}
