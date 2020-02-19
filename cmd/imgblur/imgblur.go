package main

import (
	"image"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/fogleman/gg"
)

func main() {
	// draw the source of light
	// it's the buffer on which we will apply the bloom effect
	dc := gg.NewContext(200, 200)
	dc.SetLineWidth(3.0)
	dc.SetRGB255(147, 112, 219)
	dc.DrawRectangle(10, 10, 180, 180)
	dc.DrawCircle(100, 100, 50)
	dc.Stroke()

	// store the original source of light to draw it back later
	original := dc.Image()

	// bloom this source of light
	bloomed := Bloom(dc.Image())

	// now, let's do our final rendering, let's starts by rendering
	// a gray rectangle in a new buffer
	dc = gg.NewContext(220, 220)
	dc.SetRGB255(40, 40, 40)
	dc.DrawRectangle(0, 0, 220, 220)
	dc.Fill()

	// draw our bloomed light
	dc.DrawImage(bloomed, 0, 0)

	// re-apply the original source of light
	dc.DrawImage(original, 10, 10)

	// save the result in a PNG
	dc.SavePNG("output.png")
}

// Bloom applies a bloom effect on the given image.
// Because of the nature of the effect, a larger image is returned.
// 10px padding is added to each side of the image, growing it by
// 20px on X and 20px on Y.
func Bloom(img image.Image) image.Image {
	// create a larger image
	size := img.Bounds().Size()
	newSize := image.Rect(0, 0, size.X+20, size.Y+20)

	// copy the original in this larger image, slightly translated to the center
	var extended image.Image
	extended = translateImage(img, newSize, 10, 10)

	// dilate the image to have a bigger source of light
	dilated := effect.Dilate(extended, 3)

	// blur the image
	bloomed := blur.Gaussian(dilated, 10.0)

	return bloomed
}

// translateImage copies the src image applying the given offset on a new Image
// bounds is the size of the resulting image.
func translateImage(src image.Image, bounds image.Rectangle, xOffset, yOffset int) image.Image {
	rv := image.NewRGBA(bounds)
	size := src.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			rv.Set(xOffset+x, yOffset+y, src.At(x, y))
		}
	}
	return rv
}
