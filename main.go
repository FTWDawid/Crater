package main

import (
	"crater/routes"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(ratelimit.RateLimiter(ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 90,
	}), &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
			c.JSON(429, gin.H{
				"status":  429,
				"message": "Too many requests",
				"retryIn": time.Until(info.ResetTime).Round(time.Second),
			})
		},
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}))

	routes.RegisterGeneratorsRoutes(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Page not found",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"status":  405,
			"message": "Method not allowed",
		})
	})

	router.Run(":4496")
}
