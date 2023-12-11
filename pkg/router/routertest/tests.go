package routertest

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
)

type RequestOptions struct {
	Headers map[string]string
	Params  map[string]string
	Body    interface{}
}

func Request(method string, handler gin.HandlerFunc, opts RequestOptions) *httptest.ResponseRecorder {
	var route, path string
	for k, v := range opts.Params {
		route += fmt.Sprintf("/:%v", k)
		path += fmt.Sprintf("/%v", v)
	}

	// Setup
	r := router.New(nil)
	r.Handle(route, method, []fizz.OperationOption{}, handler)
	w := httptest.NewRecorder()

	// Make the request
	body := strings.NewReader(fmt.Sprintf("%v", opts.Body))
	req := httptest.NewRequest(method, path, body)

	// Add headers
	if opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	r.ServeHTTP(w, req)
	return w
}

func ToJSON(t interface{}) string {
	b, _ := json.Marshal(t)
	return string(b)
}
