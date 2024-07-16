package main

import (
	"crater/routes"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ratelimitStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 90,
	})

	ratelimitMiddleware := ratelimit.RateLimiter(ratelimitStore, &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
			c.JSON(429, gin.H{
				"status":  429,
				"message": "Too many requests",
				"retryIn": info.ResetTime.Sub(time.Now()).Round(time.Second).String(),
			})
		},
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
	})

	routes.GeneratorsRoute(r, ratelimitMiddleware)

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
