package routertest

import (
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
)

type RequestOptions struct {
	Headers map[string]string
}

func Request(method string, handler gin.HandlerFunc, opts RequestOptions) *httptest.ResponseRecorder {
	// Setup
	r := router.New(nil)
	r.Handle(method, "/", []fizz.OperationOption{}, handler)
	w := httptest.NewRecorder()

	// Make the request
	req := httptest.NewRequest(method, "/", nil)

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
