// Ref: https://go.dev/tour/methods/25

package main

import (
	"golang.org/x/tour/pic"
	"image"
	"image/color"
)
	
type Image struct{
	W,H uint8
	Color uint8
}

func (myImg Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(myImg.W), int(myImg.H))
}

func (image Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (image Image) At(x,y int) color.Color {
	return color.RGBA{ uint8(x * y), uint8(x + y), 255, 255 }
}

func main() {
	m := Image{50, 50, 127}
	pic.ShowImage(m)
}
