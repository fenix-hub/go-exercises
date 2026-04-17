// Ref: https://go.dev/tour/moretypes/18

package main

import (
	"golang.org/x/tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	var img [][]uint8 = make([][]uint8, dy)
	for i := range img {
		img[i] = make([]uint8, dx)
		for j := range len(img) {
			img[i][j] = uint8( (i+j)/2 )
		}
	}

	return img
}

func main() {
	pic.Show(Pic)
}
