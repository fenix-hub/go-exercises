// Ref: https://go.dev/tour/methods/20

package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	var err error 
	if x < 0 {
		err = ErrNegativeSqrt(x)
	}
	return math.Sqrt(x), err
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
