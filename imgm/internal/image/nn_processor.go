package image

import (
	"cmp"
	"image"
	"image/color"
	"math"
	"slices"
)

type NNProcessor struct {
	colorSchema []color.Color
}

func (nn *NNProcessor) convert(pixel color.Color) color.Color {
	if len(nn.colorSchema) == 0 {
		return pixel
	}

	nearest := slices.MinFunc(nn.colorSchema, func(a, b color.Color) int {
		return cmp.Compare(nn.distance(pixel, a), nn.distance(pixel, b))
	})

	return nearest
}

func (nn *NNProcessor) distance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	dr := float64(r1) - float64(r2)
	dg := float64(g1) - float64(g2)
	db := float64(b1) - float64(b2)

	return math.Sqrt(dr*dr + dg*dg + db*db)
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
