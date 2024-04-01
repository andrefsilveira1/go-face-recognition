package main

import (
	"fmt"
	"image"
	"log"
	"math"

	"github.com/gonum/matrix/mat64"
)

func compare(image1, image2 *mat64.Dense, threshold float64) bool {
	fmt.Println("Comparing images...")
	rows1, cols1 := image1.Dims()
	rows2, cols2 := image2.Dims()
	if rows1 != rows2 || cols1 != cols2 {
		log.Fatal("Error: Image dimensions do not match")
		return false
	}

	diff := mat64.NewDense(rows1, cols1, nil)
	diff.Sub(image1, image2)

	var sumSquared float64
	diff.Apply(func(i, j int, v float64) float64 {
		sumSquared += v * v
		return v
	}, diff)

	rmsDiff := sumSquared / float64(rows1*cols1)
	rmsDiff = math.Sqrt(rmsDiff)

	return rmsDiff <= threshold

}

func Identify(images []image.Image) {
	fmt.Println("Images received!")
}
