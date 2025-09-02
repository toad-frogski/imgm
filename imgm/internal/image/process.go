package image

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"strconv"
)

type (
	ImageProcessor interface {
		Process(image.Image) image.Image
	}

	Options struct {
		ImageP    string
		SchemaP   string
		Processor string
		OutputP   string
	}
)

func Process(options Options) error {
	img, err := loadImage(options.ImageP)
	if err != nil {
		return err
	}

	colors, err := loadColors(options.SchemaP)
	if err != nil {
		return err
	}

	processor, err := getProcessor(options.Processor, colors)
	if err != nil {
		return err
	}

	updated := processor.Process(img)

	if err = saveImage(updated, options.OutputP); err != nil {
		return err
	}

	return nil
}

func loadColors(schemaP string) ([]color.Color, error) {
	f, err := os.Open(schemaP)
	if err != nil {
		return nil, err
	}

	parseHexNibble := func(c byte) uint8 {
		v, _ := strconv.ParseUint(string([]byte{c, c}), 16, 8)
		return uint8(v)
	}

	var hexColors []string
	if err = json.NewDecoder(f).Decode(&hexColors); err != nil {
		return nil, err
	}

	colors := make([]color.Color, 0, len(hexColors))

	for _, hex := range hexColors {
		if len(hex) > 0 && hex[0] == '#' {
			hex = hex[1:]
		}

		var r, g, b uint8

		switch len(hex) {
		case 6:
			val, err := strconv.ParseUint(hex, 16, 32)
			if err != nil {
				return nil, err
			}
			r = uint8(val >> 16)
			g = uint8((val >> 8) & 0xFF)
			b = uint8(val & 0xFF)
		case 3:
			r = parseHexNibble(hex[0])
			g = parseHexNibble(hex[1])
			b = parseHexNibble(hex[2])
		default:
			return nil, fmt.Errorf("unsupported hex color format: %s", hex)
		}

		colors = append(colors, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	return colors, nil
}

func loadImage(imgP string) (image.Image, error) {
	f, err := os.Open(imgP)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s", imgP)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("image format not allowed. please, use one of these: [jpeg, png]")
	}

	return img, nil
}

func saveImage(img image.Image, outputP string) error {
	out, err := os.Create(outputP)
	if err != nil {
		return fmt.Errorf("unable to save file: %v", err)
	}
	defer out.Close()

	if err = png.Encode(out, img); err != nil {
		return fmt.Errorf("unable to save file: %v", err)
	}

	return nil
}

func getProcessor(processor string, colors []color.Color) (ImageProcessor, error) {
	switch processor {
	case "nn":
		return &NNProcessor{colorSchema: colors}, nil

	default:
		return nil, fmt.Errorf("unknown processor: %s", processor)
	}
}
