package httpservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetBadges(c *gin.Context) {
	err := h.app.credly.GetBadges()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
