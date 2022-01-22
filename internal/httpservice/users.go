package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetUsers(c *gin.Context) {
	users, err := h.app.marvel.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}
