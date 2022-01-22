package httpservice

import (
	"github.com/DavidHaugen/golang-boilerplate/internal/config"
	"github.com/gin-gonic/gin"
)

func ListenAndServe(config *config.Config) {
	applicationConfig := newApplicationConfig(config)
	handler := newHandler(applicationConfig)
	router := getRouter(handler)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getRouter(handler RouteHandler) *gin.Engine {
	router := gin.New()
	addRoutes(router, handler)
	return router
}

func addRoutes(router *gin.Engine, handler RouteHandler) {
	addBaseRoutes(router.Group(""), handler)
}

func addBaseRoutes(rg *gin.RouterGroup, handler RouteHandler) {
	rg.GET("/ping", handler.Ping)
	rg.GET("/users", handler.GetUsers)
}
