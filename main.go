package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"path/filepath"

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

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func main() {

	imagesDir := "go-opencv/images"
	image1 := "image1.jpg"
	image2 := "image2.jpg"

	image1Path := filepath.Join(imagesDir, image1)
	image2Path := filepath.Join(imagesDir, image2)

	if compare(image1Path, image2Path, 0.1) {
		fmt.Println("EQUAL!")
	} else {
		fmt.Println("FALSE")
	}
	fmt.Println("Images received!")
}
