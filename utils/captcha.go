package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type CaptchaGenerator struct {
	width      int
	height     int
	fontFace   font.Face
	background color.Color
	image      *image.RGBA
	drawer     *font.Drawer
}

func NewCaptchaGenerator(fontFile string, fontSize, width, height int, backgroundColor color.Color) (*CaptchaGenerator, error) {
	rand.Seed(time.Now().UnixNano())

	fontFace, err := loadFont(fontFile, float64(fontSize))
	if err != nil {
		return nil, fmt.Errorf("error loading font: %w", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{225, 35, 60, 255}),
		Face: fontFace,
		Dot:  fixed.Point26_6{},
	}

	return &CaptchaGenerator{
		width:      width,
		height:     height,
		fontFace:   fontFace,
		background: backgroundColor,
		image:      img,
		drawer:     drawer,
	}, nil
}

func (cg *CaptchaGenerator) GenerateImage(value string) *image.RGBA {
	draw.Draw(cg.image, cg.image.Bounds(), &image.Uniform{cg.background}, image.Point{}, draw.Src)

	textWidth := cg.drawer.MeasureString(value)

	textX := (cg.width - int(textWidth.Round())) / 2

	ascent := cg.fontFace.Metrics().Ascent.Ceil()
	descent := cg.fontFace.Metrics().Descent.Ceil()
	totalHeight := ascent + descent

	textY := (cg.height-totalHeight)/2 + ascent

	cg.drawer.Dot = fixed.Point26_6{
		X: fixed.I(textX),
		Y: fixed.I(textY),
	}
	cg.drawer.DrawString(value)

	cg.drawRandomDots(0.015)

	return cg.image
}

func (cg *CaptchaGenerator) drawRandomDots(radiusRatio float64) {
	radius := int(float64(cg.width) * radiusRatio)
	numDots := 100

	for i := 0; i < numDots; i++ {
		x := rand.Intn(cg.width)
		y := rand.Intn(cg.height)

		c := color.RGBA{
			225 - uint8(rand.Uint32()%50+80),
			35 - uint8(rand.Uint32()%15+15),
			60 - uint8(rand.Uint32()%15+15),
			255,
		}

		drawCircle(cg.image, x, y, radius, c)
	}
}

func drawCircle(img *image.RGBA, x, y, radius int, c color.Color) {
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x+dx, y+dy, c)
			}
		}
	}
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
