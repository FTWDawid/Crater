package routes

import (
	"image/png"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func CaptchaGenerator(c *gin.Context) {
	height := captcha.StdHeight
	width := captcha.StdWidth

	if queryHeight := c.Query("height"); queryHeight != "" {
		if h, err := strconv.Atoi(queryHeight); err == nil && h > 0 {
			height = h
		}
	}

	if queryWidth := c.Query("width"); queryWidth != "" {
		if w, err := strconv.Atoi(queryWidth); err == nil && w > 0 {
			width = w
		}
	}

	id := captcha.New()

	captchaDigits := captcha.RandomDigits(6)
	captchaImage := captcha.NewImage(id, captchaDigits, width, height)

	c.Header("Content-Type", "image/png")
	_ = png.Encode(c.Writer, captchaImage)
}

func GeneratorsRoute(ginInstance *gin.Engine, limiter gin.HandlerFunc) {
	generator := ginInstance.Group("/v1/generators", limiter)

	generator.GET("/captcha", CaptchaGenerator)
}
