package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/check", func(c *gin.Context) {
		domain := c.Query("domain")
		if domain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Domain parameter is required"})
			return
		}

		status := Check(domain, "80")

		c.JSON(http.StatusOK, status)
	})

	r.Run(":8080")
}

