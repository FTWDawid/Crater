package main

import (
	"crater/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.GeneratorsRoute(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Page not found",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"status":  405,
			"message": "Method not allowed",
		})
	})

	r.Run(":8080")
}
