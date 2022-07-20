package main

import (
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/nfnt/resize"
)

var watermarkPath string
var inputFilesPath string
var outputFilesPath string
var offsetX int
var offsetY int
var heightPercentage int
var outputSizes string
var jpegQuality int

func main() {
	flag.StringVar(&watermarkPath, "watermark", "watermark.png", "Path and file of the watermark (must be .png)")
	flag.StringVar(&inputFilesPath, "input", "./in", "Path to the folder containing the input files (must be .jpg)")
	flag.StringVar(&outputFilesPath, "output", "./out", "Path for the images with watermark")
	flag.IntVar(&offsetX, "offset_x", 0, "Distance of the watermark to the left side of the image")
	flag.IntVar(&offsetY, "offset_y", 0, "Distance of the watermark to the bottom side of the image")
	flag.IntVar(&heightPercentage, "height_percentage", 10, "Percentage of the height of the watermark (relative to the image it is placed on)")
	flag.StringVar(&outputSizes, "output_sizes", "source", "List of sizes in which the output images are stored (width in pixels) separated by comma. Special value: source")
	flag.IntVar(&jpegQuality, "jpeg_quality", 85, "Quality of the output image (ranges from 1 to 100 inclusive, higher is better)")
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

	var sizes = map[string]uint{}
	log.Println("parsing output size values (-output_sizes)")
	for _, size := range strings.Split(strings.ReplaceAll(outputSizes, " ", ""), ",") {
		if size == "source" {
			sizes["source"] = 0
		}
		if size == "" || size == "source" {
			continue
		}
		s, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			log.Fatalf("failed to parse value of flag -output_sizes (must be number):%v\n", err)
		}
		sizes[size] = uint(s)
	}

	if len(sizes) == 0 {
		log.Fatalf("Need output file sizes (flag: -output_sizes)")
	}

	log.Printf("creating output directories at %q\n", outputFilesPath)
	for sizeName := range sizes {
		log.Printf(" >%s", path.Join(outputFilesPath, sizeName))
		err = os.MkdirAll(path.Join(outputFilesPath, sizeName), os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create output directory: %v\n", err)
		}
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

			outputFile(file.Name(), watermarkedImage, sizes)
			watermarkedImageCount++
		}
	}
	log.Printf("watermarked %d images with %q", watermarkedImageCount, watermarkPath)
}

func outputFile(outputFile string, image *image.RGBA, sizes map[string]uint) {
	for sizeName, size := range sizes {
		log.Printf("  >width: %s\n", sizeName)
		outputImage, err := os.Create(path.Join(outputFilesPath, sizeName, outputFile))
		if err != nil {
			log.Warnf("failed to create output file: %v\n", err)
			continue
		}
		imageInNewSize := resize.Resize(size, 0, image, resize.Lanczos3)
		err = jpeg.Encode(outputImage, imageInNewSize, &jpeg.Options{Quality: jpegQuality})

		if err != nil {
			log.Warnf("failed to encode output: %v\n", err)
		}
		outputImage.Close()
	}
}
