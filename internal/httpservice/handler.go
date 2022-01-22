package httpservice

import (
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -package=mock -destination=mock/handler.go

// RouteHandler :
type RouteHandler interface {
	Ping(c *gin.Context)
}

// Handler : handles all http requests
type Handler struct {
	app *applicationConfig
}

// newHandler : returns a configured route handler.
func newHandler(config *applicationConfig) RouteHandler {
	return Handler{config}
}
