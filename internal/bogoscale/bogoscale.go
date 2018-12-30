package bogoscale

import (
	"image"
	"image/color"
	"math"
)

// bogo scale.  apply an averaging type filter to input image and translate to
// resulting bitmap.
func Scale(img image.Image, width int, height int) image.Image {
	rect := img.Bounds()

	// 64 bits might be overzealous
	type point struct {
		Red   uint64
		Green uint64
		Blue  uint64
		Count uint64
	}

	output := make([]point, width*height, width*height)

	yscale := float64(height) / float64(rect.Max.Y)
	xscale := float64(width) / float64(rect.Max.X)

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			// map x and y

			c := img.At(x, y)

			// we don't care about the alpha channel
			r, g, b, _ := c.RGBA()

			xx, yy := int(math.Floor(float64(x)*xscale)), int(math.Floor(float64(y)*yscale))
			pos := (yy * width) + xx

			output[pos].Red += uint64(r >> 8)
			output[pos].Green += uint64(g >> 8)
			output[pos].Blue += uint64(b >> 8)
			output[pos].Count++
		}
	}

	rc := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := (y * width) + x

			if output[pos].Count == 0 {
				// zerovalue is fine for this 'pixel'
				continue
			}

			avg := func(v uint64) uint8 {
				return uint8(v / output[pos].Count)
			}

			c := color.RGBA{
				R: avg(output[pos].Red),
				G: avg(output[pos].Green),
				B: avg(output[pos].Blue),
				A: 0xff,
			}
			rc.Set(x, y, c)
		}
	}

	return rc
}
