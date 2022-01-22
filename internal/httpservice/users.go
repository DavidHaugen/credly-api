package httpservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetUsers(c *gin.Context) {
	users, err := h.app.marvel.GetUsers()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}
