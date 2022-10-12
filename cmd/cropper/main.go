package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

var (
	dir       string
	outputDir string
	cropSize  int
)

func init() {
	flag.StringVar(&dir, "d", "screens", "directory with screens")
	flag.StringVar(&outputDir, "o", "screens_cropped", "output directory")
	flag.IntVar(&cropSize, "s", 80, "amount of pixels to crop")
}

func main() {
	flag.Parse()

	var entries, err = os.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read input dir:", err)
		return
	}

	if err = os.Mkdir(outputDir, 0777); err != nil && os.IsNotExist(err) {
		fmt.Println("failed to create output dir:", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		processFile(dir, entry.Name(), outputDir, cropSize)
	}
}

func processFile(dir, fileName, outputDir string, cropSize int) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(filepath.Join(dir, fileName)); err != nil {
		fmt.Printf("failed to read file '%s': %v\n", fileName, err)
		return
	}
	defer f.Close()

	var img image.Image
	if img, _, err = image.Decode(f); err != nil {
		fmt.Printf("failed to decode file '%s': %v\n", fileName, err)
		return
	}

	var newImg = image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()-cropSize))

	for x := 0; x < newImg.Bounds().Dx(); x++ {
		for y := 0; y < newImg.Bounds().Dy(); y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	var newF *os.File
	if newF, err = os.Create(filepath.Join(outputDir, fileName)); err != nil {
		fmt.Printf("failed to create output file '%s': %v\n", fileName, err)
		return
	}
	defer newF.Close()

	if err = jpeg.Encode(newF, newImg, &jpeg.Options{Quality: 100}); err != nil {
		fmt.Printf("failed to encode file '%s': %v\n", fileName, err)
	}
}
