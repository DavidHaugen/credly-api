package httpservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DavidHaugen/golang-boilerplate/internal/httpservice/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

/*
 * Test routes are properly configured
 * and call the correct handler
 */

var anyGinContext = gomock.AssignableToTypeOf(&gin.Context{})

func serveHTTP(method, url string, handler RouteHandler) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(method, url, nil)

	router := getRouter(handler)
	router.ServeHTTP(response, request)
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		method string
		url    string
		expect func(h *mock.MockRouteHandler)
	}{
		{
			http.MethodGet,
			"/ping",
			func(h *mock.MockRouteHandler) {
				h.EXPECT().Ping(anyGinContext).Times(1)
			},
		},
	}

	for _, test := range tests {
		controller := gomock.NewController(t)
		handler := mock.NewMockRouteHandler(controller)
		test.expect(handler)
		serveHTTP(test.method, test.url, handler)
		controller.Finish()
	}
}
