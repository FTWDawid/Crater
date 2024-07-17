package routes

import (
	"bytes"
	"crater/utils"
	"encoding/base64"
	"image/color"
	"image/png"

	"github.com/gin-gonic/gin"
)

func HandleCaptchaRequest(ctx *gin.Context, captchaGen *utils.CaptchaGenerator) {
	ctx.Header("Content-Type", "application/json")

	randomCode := utils.GenerateRandomCode()

	img := captchaGen.GenerateImage(randomCode)
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to encode image"})
		return
	}

	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	ctx.JSON(200, gin.H{
		"code":  randomCode,
		"image": imageBase64,
	})
}

func RegisterGeneratorsRoutes(router *gin.Engine) {
	generatorGroup := router.Group("/v1/generators")

	captchaGen, error := utils.NewCaptchaGenerator("assets/poppins_bold.ttf", 64, 360, 180, color.Transparent)
	println(error)
	generatorGroup.GET("/captcha", func(ctx *gin.Context) {
		HandleCaptchaRequest(ctx, captchaGen)
	})
}
