package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	main "github.com/pranotobudi/myslack-monolith-backend"
	"github.com/stretchr/testify/assert"
)

// TestRouting will test each end point is active
func TestRouting(t *testing.T) {
	tt := []struct {
		Name   string
		Path   string
		Method string
		Status int
	}{
		{"Get home", "/", http.MethodGet, http.StatusOK},
		{"Get rooms", "/rooms", http.MethodGet, http.StatusOK},
		// {"Get room", "/room", http.MethodGet, http.StatusOK},
		// {"Post room", "/room", http.MethodPost, http.StatusOK},
		// {"Get messages", "/messages", http.MethodGet, http.StatusOK},
		// {"Get userByemail", "/userByEmail", http.MethodGet, http.StatusOK},
		// {"Post userAuth", "/userAuth", http.MethodPost, http.StatusOK},
		// {"Put updateUserRooms", "/updateUserRooms", http.MethodPut, http.StatusOK},
		// {"Get websocket", "/websocket", http.MethodGet, http.StatusOK},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.Path, nil)
			w := httptest.NewRecorder()

			router := main.Router()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			// t.Log("body: " + string(w.Body.String()))

			// srv := httptest.NewServer(Router())
			// defer srv.Close()

			// res, err := http.Get(fmt.Sprintf("%s/roomss", srv.URL))
			// if err != nil {
			// 	t.Fatalf("could not send GET request: %v", err)
			// }

			// if res.StatusCode != http.StatusOK {
			// 	t.Errorf("expected status OK, but got %v", res.StatusCode)
			// }
			// b, err := ioutil.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatalf("expected an integer; got %s", err)
			// }
		})
	}
}

func TestCORS(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// test request, must instantiate a request first
	// req := &http.Request{
	// 	URL:    &url.URL{},
	// 	Header: make(http.Header), // if you need to test headers
	// }
	req, err := http.NewRequest("GET", "/rooms", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// example: req.Header.Add("Accept", "application/json")
	c.Request = req
	main.CORS(c)

	// logrus.Info("header:", c.Writer.Header())
}
