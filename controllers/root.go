package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, os.Getenv("FE_URL"))
}
