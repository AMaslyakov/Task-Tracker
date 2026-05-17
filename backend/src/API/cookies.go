package API

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseCookie(c *gin.Context) string {
	value, err := c.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Printf("Cookies not found")
			return ""
		}
		fmt.Printf("No valid cookies")
		return ""
	}
	return value
}
