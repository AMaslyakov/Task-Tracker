package API

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertTask(c *gin.Context) {
	ctx := c.Request.Context()

	SessionId := ParseCookie(c)
	if _, ok := CorrectSession(ctx, SessionId); !ok {
		c.JSON(http.StatusFound, `Location: "Login Page"`)
	}
	var t Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("Parse JSON Error: %w", err))
		return
	}

}
