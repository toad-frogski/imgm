package image

import (
	"image"
	"image/color"
	"math"
)

type NNProcessor struct {
	colorSchema []color.Color
}

func (nn *NNProcessor) convert(pixel color.Color) color.Color {
	if len(nn.colorSchema) == 0 {
		return pixel
	}

	r0, g0, b0, _ := pixel.RGBA()
	minDist := float64(math.MaxFloat64)
	var nearest color.Color

	for _, c := range nn.colorSchema {
		r, g, b, _ := c.RGBA()
		dr := float64(r0) - float64(r)
		dg := float64(g0) - float64(g)
		db := float64(b0) - float64(b)
		dist2 := dr*dr + dg*dg + db*db
		if dist2 < minDist {
			minDist = dist2
			nearest = c
		}
	}

	return nearest
}

func (nn *NNProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			updated := nn.convert(pixel)
			newImg.Set(x, y, updated)
		}
	}

	return newImg
}
