package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("must be 2 args")
	}
	var names = os.Args[1:]

	var start = time.Now()
	var images = readAll(names)

	y1, y2, length := overlaps(images[0], images[1])
	fmt.Println(y1, y2, length)
	fmt.Println(time.Since(start).String())

	// var a = imageToRows(images[0], "0.png")
	// var b = imageToRows(images[1], "1.png")
	//
	// var originY, otherY, length = a.overlaps(b)
	//
	// fmt.Println(originY, otherY, length)
	//
	// // debug print union
	// var res = image.NewRGBA(image.Rect(0, 0, 1440, length))
	// for x := 0; x < 720; x++ {
	// 	for y := originY; y < originY+length; y++ {
	// 		res.Set(x, y-originY, a.original.At(x, y))
	// 	}
	// }
	//
	// for x := 0; x < 720; x++ {
	// 	for y := otherY; y < otherY+length; y++ {
	// 		res.Set(x+720, y-otherY, b.original.At(x, y))
	// 	}
	// }
	//
	// f, _ := os.Create("union.jpg")
	// jpeg.Encode(f, res, nil)
	// f.Close()
	//
	// print final image
	var sizeOfOrigin = y1 + length
	var sizeOfOther = images[1].Bounds().Dy() - length - y2
	var res2 = image.NewRGBA(image.Rect(0, 0, images[0].Bounds().Dx(), sizeOfOrigin+sizeOfOther))

	for x := 0; x < images[0].Bounds().Dx(); x++ {
		for y := 0; y < y1+length; y++ {
			res2.Set(x, y, images[0].At(x, y))
		}

		for y := y2 + length; y < images[1].Bounds().Dy(); y++ {
			res2.Set(x, y-y2-length+sizeOfOrigin, images[1].At(x, y))
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

func overlaps(i1, i2 image.Image) (y1, y2, length int) {
	for offset := 1; offset < i1.Bounds().Dy()*2; offset++ {
		var (
			from1 = offset - i1.Bounds().Dy()
			to1   = from1 + i1.Bounds().Dy()
			from2 = i2.Bounds().Dy() - offset
			to2   = from2 + i2.Bounds().Dy()
		)
		if from1 < 0 {
			from1 = 0
		}
		if to1 > i1.Bounds().Dy() {
			to1 = i1.Bounds().Dy()
		}
		if from2 < 0 {
			from2 = 0
		}
		if to2 > i2.Bounds().Dy() {
			to2 = i2.Bounds().Dy()
		}

		if to1-from1 != to2-from2 {
			fmt.Println(to1-from1, to2-from2)
			panic("damn")
		}

		var hasTrues = false
		var mask = make([]bool, to1-from1)

		for yOffset := 0; yOffset < to1-from1; yOffset++ {
			var allSame = true
			for x := 0; x < i1.Bounds().Dx(); x++ {
				if !closeEnough(i1.At(x, from1+yOffset), i2.At(x, from2+yOffset)) {
					allSame = false
					break
				}
			}
			if allSame == true {
				hasTrues = true
			}
			mask[yOffset] = allSame
		}

		var maxPos, maxLen = 0, 0
		if hasTrues {
			maxPos, maxLen = getLongestSequenceOfTrue(mask)
		}

		// fmt.Println("i:", offset)
		if maxLen > length {
			length = maxLen
			y1 = maxPos + from1
			y2 = maxPos + from2

			// fmt.Println("Var:", y1, y2, length)
		}
	}

	return y1, y2, length
}
