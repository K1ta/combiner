package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

var names = []string{"photo_5341437942142451795_y.jpg", "photo_5341437942142451796_y.jpg"}

func main() {
	var images = readAll(names)

	var a = imageToRows(images[0], "0.png")
	var b = imageToRows(images[1], "1.png")

	var originY, otherY, length = a.overlaps(b)

	fmt.Println(originY, otherY, length)

	// debug print union
	var res = image.NewRGBA(image.Rect(0, 0, 1440, length))
	for x := 0; x < 720; x++ {
		for y := originY; y < originY+length; y++ {
			res.Set(x, y-originY, a.original.At(x, y))
		}
	}

	for x := 0; x < 720; x++ {
		for y := otherY; y < otherY+length; y++ {
			res.Set(x+720, y-otherY, b.original.At(x, y))
		}
	}

	f, _ := os.Create("union.jpg")
	jpeg.Encode(f, res, nil)
	f.Close()

	// print final image
	var sizeOfOrigin = originY + length
	var sizeOfOther = images[1].Bounds().Dy() - length - otherY
	var res2 = image.NewRGBA(image.Rect(0, 0, images[0].Bounds().Dx(), sizeOfOrigin+sizeOfOther))

	for x := 0; x < images[0].Bounds().Dx(); x++ {
		for y := 0; y < originY+length; y++ {
			res2.Set(x, y, images[0].At(x, y))
		}

		for y := otherY + length; y < images[1].Bounds().Dy(); y++ {
			res2.Set(x, y-otherY-length+sizeOfOrigin, images[1].At(x, y))
		}
	}

	f1, _ := os.Create("woohoo.jpg")
	jpeg.Encode(f1, res2, nil)
	f1.Close()
}

func readAll(names []string) []image.Image {
	var res = make([]image.Image, len(names))

	for i := range names {
		var f, err = os.Open(names[i])
		if err != nil {
			panic(err)
		}

		if res[i], err = jpeg.Decode(f); err != nil {
			fmt.Println("failed to decode " + names[i])
			panic(err)
		}

		f.Close()
	}

	return res
}

type ImageRows struct {
	original image.Image
	rows     []string
}

func imageToRows(img image.Image, name string) *ImageRows {
	var rows = make([]string, img.Bounds().Dy())

	return &ImageRows{
		original: img,
		rows:     rows,
	}
}

// mask is never empty
func getLongestSequenceOfTrue(mask []bool) (pos, length int) {
	var (
		curPos = 0
		count  = 0
	)
	for j := range mask {
		if mask[j] == true {
			if count == 0 {
				curPos = j
			}
			count++
		} else {
			if count > length {
				length = count
				pos = curPos
			}
			count = 0
		}
	}

	if count > 0 && count > length {
		length = count
		pos = curPos
	}

	return pos, length
}

func (origin *ImageRows) overlaps(other *ImageRows) (originY, otherY, length int) {
	var (
		size        = origin.original.Bounds().Dy()
		originSlice = getTwoDimensionSlice(origin.original, color.White, size)
	)

	for i := 0; i < size*2; i++ {
		var otherSlice = getTwoDimensionSlice(other.original, color.Black, i)

		var mask = getRowBoolMask(originSlice, otherSlice)

		// now get the longest 'true' sequence
		var maxPos, maxLen = getLongestSequenceOfTrue(mask)

		fmt.Println("i:", i)
		if maxLen > length {
			length = maxLen
			originY = maxPos - size
			otherY = maxPos - i

			fmt.Println("Var:", originY, otherY, length)
		}
	}

	return originY, otherY, length
}

func getTwoDimensionSlice(img image.Image, base color.Color, offset int) [][]color.Color {
	var res = make([][]color.Color, img.Bounds().Dy()*3)

	for i := 0; i < len(res); i++ {
		res[i] = make([]color.Color, img.Bounds().Dx())
		for j := range res[i] {
			if i < offset || i >= offset+img.Bounds().Dy() {
				res[i][j] = base
			} else {
				res[i][j] = img.At(j, i-offset)
			}
		}
	}

	return res
}

func getRowBoolMask(sl1, sl2 [][]color.Color) []bool {
	var res = make([]bool, len(sl1))

	for i := range sl1 {
		var allSame = true
		for j := range sl1[i] {
			if !closeEnough(sl1[i][j], sl2[i][j]) {
				allSame = false
				break
			}
		}
		res[i] = allSame
	}

	return res
}

func closeEnough(c1, c2 color.Color) bool {
	var (
		r1, g1, b1, _ = c1.RGBA()
		r2, g2, b2, _ = c2.RGBA()

		diffR = float64(r1) - float64(r2)
		diffG = float64(g1) - float64(g2)
		diffB = float64(b1) - float64(b2)
	)
	res := math.Sqrt(diffR*diffR + diffG*diffG + diffB*diffB)
	return res < 60000
}
