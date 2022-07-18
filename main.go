package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/nfnt/resize"
)

var watermarkPath string
var inputFilesPath string
var outputFilesPath string
var offsetX int
var offsetY int
var heightPercentage int

func main() {
	flag.StringVar(&watermarkPath, "watermark", "watermark.png", "Path and file of the watermark (must be .png)")
	flag.StringVar(&inputFilesPath, "input", "./in", "Path to the folder containing the input files (must be .jpg)")
	flag.StringVar(&outputFilesPath, "output", "./out", "Path for the images with watermark")
	flag.IntVar(&offsetX, "offset_x", 0, "Distance of the watermark to the left side of the image")
	flag.IntVar(&offsetY, "offset_y", 0, "Distance of the watermark to the bottom side of the image")
	flag.IntVar(&heightPercentage, "height_percentage", 10, "Percentage of the height of the watermark (relative to the image it is placed on)")
	flag.Parse()

	log.Printf("reading file of watermark at %q\n", watermarkPath)
	wmb, err := os.Open(watermarkPath)
	if err != nil {
		log.Fatalf("failed to open file of watermark: %v\n", err)
	}
	watermark, err := png.Decode(wmb)
	if err != nil {
		log.Fatalf("failed decode watermark: %v\n", err)
	}
	defer wmb.Close()

	log.Printf("getting all images at %q\n", inputFilesPath)
	files, err := ioutil.ReadDir(inputFilesPath)
	if err != nil {
		log.Fatalf("failed to get all files of input directory: %v\n", err)
	}

	log.Printf("creating output directory at %q\n", outputFilesPath)
	err = os.MkdirAll(outputFilesPath, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory: %v\n", err)
	}

	log.Println("generating watermarked images")

	var watermarkedImageCount = 0

	for _, file := range files {
		if file.Mode().IsRegular() {
			filePath := path.Join(inputFilesPath, file.Name())

			log.Printf(" >%s\n", filePath)

			imageFile, err := os.Open(filePath)
			if err != nil {
				log.Fatalf("failed to read image: %v\n", err)
			}
			img, err := jpeg.Decode(imageFile)

			imageFile.Close()
			if err != nil {
				log.Fatalf("failed to decode image: %v\n", err)
			}

			bounds := img.Bounds()
			watermarkedImage := image.NewRGBA(bounds)
			draw.Draw(watermarkedImage, bounds, img, image.Point{}, draw.Src)

			waterMarkForThisImage := resize.Resize(0, uint(heightPercentage*bounds.Dy()/100), watermark, resize.Lanczos3)
			draw.Draw(watermarkedImage, waterMarkForThisImage.Bounds().Add(image.Point{
				X: offsetX,
				Y: bounds.Dy() - waterMarkForThisImage.Bounds().Dy() - offsetY,
			}), waterMarkForThisImage, image.Point{}, draw.Over)

			outputImage, err := os.Create(path.Join(outputFilesPath, file.Name()))
			if err != nil {
				log.Fatalf("failed to create output file: %v\n", err)
			}
			err = jpeg.Encode(outputImage, watermarkedImage, &jpeg.Options{Quality: 100})
			if err != nil {
				log.Fatalf("failed to encode output: %v\n", err)
			}
			outputImage.Close()
			watermarkedImageCount++
		}
	}
	log.Printf("watermarked %d images with %q", watermarkedImageCount, watermarkPath)
}
