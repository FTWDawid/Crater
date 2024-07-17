package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type CaptchaGenerator struct {
	width   int
	height  int
	font    font.Face
	bgColor color.Color
	img     *image.RGBA
	d       *font.Drawer
}

func NewCaptchaGenerator(fontFile string, fontSize, width, height int, bgColor color.Color) (*CaptchaGenerator, error) {
	fontFace, err := loadFont(fontFile, float64(fontSize))
	if err != nil {
		return nil, fmt.Errorf("error loading font: %w", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{225, 35, 60, 255}), // Customize this as needed
		Face: fontFace,
		Dot:  fixed.Point26_6{},
	}

	return &CaptchaGenerator{
		width:   width,
		height:  height,
		font:    fontFace,
		bgColor: bgColor,
		img:     img,
		d:       d,
	}, nil
}

func (cg *CaptchaGenerator) GenerateImage(value string) *image.RGBA {
	draw.Draw(cg.img, cg.img.Bounds(), &image.Uniform{cg.bgColor}, image.Point{}, draw.Src)

	textWidth := cg.d.MeasureString(value)

	textX := (cg.width - int(textWidth.Round())) / 2

	ascent := cg.font.Metrics().Ascent.Ceil()
	descent := cg.font.Metrics().Descent.Ceil()
	totalHeight := ascent + descent

	textY := (cg.height-totalHeight)/2 + ascent

	cg.d.Dot = fixed.Point26_6{
		X: fixed.I(textX),
		Y: fixed.I(textY),
	}
	cg.d.DrawString(value)

	return cg.img
}

func loadFont(fontFile string, size float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, fmt.Errorf("error reading font file: %w", err)
	}

	parsedFont, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}

	return truetype.NewFace(parsedFont, &truetype.Options{
		Size: size,
		DPI:  72,
	}), nil
}
