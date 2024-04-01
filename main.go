package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
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
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func imageToGrayscaleMatrix(img *image.Gray) *mat64.Dense {
	bounds := img.Bounds()
	rows, cols := bounds.Dy(), bounds.Dx()
	data := make([]float64, rows*cols)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			grayValue := img.GrayAt(x, y).Y
			data[y*cols+x] = float64(grayValue)
		}
	}
	return mat64.NewDense(rows, cols, data)
}

func compareResize(image1, image2 *mat64.Dense, threshold float64) bool {
	fmt.Println("Comparing images...")
	rows1, cols1 := image1.Dims()
	rows2, cols2 := image2.Dims()

	// Resize images to the same dimensions if they are different
	if rows1 != rows2 || cols1 != cols2 {
		// Resize image2 to match dimensions of image1
		resizedImage2 := imaging.Resize(image2, cols1, rows1, imaging.Linear)

		// Convert resized image2 to *mat64.Dense
		rows2, cols2 := resizedImage2.Bounds().Dy(), resizedImage2.Bounds().Dx()
		data2 := make([]float64, rows2*cols2)
		for y := 0; y < rows2; y++ {
			for x := 0; x < cols2; x++ {
				r, _, _, _ := resizedImage2.At(x, y).RGBA()
				data2[y*cols2+x] = float64(r)
			}
		}
		image2 = mat64.NewDense(rows2, cols2, data2)
	}

	// Perform comparison
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

func matToImage(mat *mat64.Dense, width, height int) image.Image {
	rows, cols := mat.Dims()

	img := image.NewGray(image.Rect(0, 0, width, height))

	data := mat.RawMatrix().Data
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			val := uint8(data[y*cols+x])
			img.SetGray(x, y, color.Gray{val})
		}
	}

	return img
}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	imagesDir := filepath.Join(dir, "images")
	image1 := "image1.jpg"
	image2 := "image6.jpg"

	image1Path := filepath.Join(imagesDir, image1)
	image2Path := filepath.Join(imagesDir, image2)

	img1, err := loadImage(image1Path)

	if err != nil {
		fmt.Println("Error loading image 1:", err)
		return
	}
	img2, err := loadImage(image2Path)

	if err != nil {
		fmt.Println("Error loading image 2:", err)
		return
	}

	gray1 := image.NewGray(img1.Bounds())
	gray2 := image.NewGray(img2.Bounds())

	newImage1 := imageToGrayscaleMatrix(gray1)
	newImage2 := imageToGrayscaleMatrix(gray2)

	if compare(newImage1, newImage2, 0.1) {
		fmt.Println("EQUAL!")
	} else {
		fmt.Println("FALSE")
	}

	fmt.Println("Images received!")
}
