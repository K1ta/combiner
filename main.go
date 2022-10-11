package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"strconv"
)

var names = []string{"photo_5341437942142451795_y.jpg", "photo_5341437942142451796_y.jpg"}

func main() {
	var images = readAll(names)

	var a = imageToRows(images[0], "0.png")
	var b = imageToRows(images[1], "1.png")

	var originY, otherY, length = a.overlaps2(b)

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

	// for i := 0; i < len(a.rows); i++ {
	// 	if a.rows[i] != b.rows[i] {
	// 		fmt.Println(i)
	// 	}
	// }

	// var res = image.NewRGBA(image.Rect(0, 0, 720, 2560))
	//
	// for x := 0; x < 720; x++ {
	// 	for y := 0; y < originY; y++ {
	// 		res.Set(x, y, images[0].At(x, y))
	// 	}
	// }
	//
	// for x := 0; x < 720; x++ {
	// 	for y := otherY; y < otherY+1280; y++ {
	// 		res.Set(x, y, images[1].At(x, y-otherY))
	// 	}
	// }
	//
	// var fres, err = os.Create("res2.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fres.Close()
	//
	// if err = jpeg.Encode(fres, res, nil); err != nil {
	// 	panic(err)
	// }

	// ---------

	// fmt.Println(findMaxIntersect(images[0], images[1]))

	// var width, height = 0, 0
	//
	// for i := range images {
	// 	width += images[i].Bounds().Size().X
	// 	if images[i].Bounds().Size().Y > height {
	// 		height = images[i].Bounds().Size().Y
	// 	}
	// }
	//
	// var res = image.NewRGBA(image.Rect(0, 0, width, height))
	//
	// var offset = 0
	// for _, img := range images {
	// 	for x := 0; x < img.Bounds().Size().X; x++ {
	// 		for y := 0; y < img.Bounds().Size().Y; y++ {
	// 			res.Set(x+offset, y, img.At(x, y))
	// 		}
	// 	}
	// 	offset += img.Bounds().Size().X
	// }
	//
	// var fres, err = os.Create("res.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fres.Close()
	//
	// if err = jpeg.Encode(fres, res, nil); err != nil {
	// 	panic(err)
	// }
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

	// var newImg = image.NewRGBA(img.Bounds())
	//
	// const d = 65535 / 255
	//
	// for x := 0; x < img.Bounds().Dx(); x++ {
	// 	for y := 0; y < img.Bounds().Dy(); y++ {
	// 		var r, g, b, _ = img.At(x, y).RGBA()
	// 		var ri, gi, bi = r / d, g / d, b / d
	// 		// rows[y] += strconv.FormatUint(ri, 10) + strconv.FormatUint(gi, 10) +
	// 		//				strconv.FormatUint(bi, 10) + strconv.FormatUint(ai, 10)
	// 		var closest = getClosestColor(ri, gi, bi)
	// 		rows[y] += closest.getString()
	// 		newImg.Set(x, y, color.RGBA{R: uint8(closest.r), G: uint8(closest.g), B: uint8(closest.b), A: uint8(255)})
	// 	}
	// }
	//
	// var newFile, err = os.Create(name)
	// if err != nil {
	// 	panic(err)
	// }
	// defer newFile.Close()
	//
	// if err = png.Encode(newFile, newImg); err != nil {
	// 	panic(err)
	// }

	return &ImageRows{
		original: img,
		rows:     rows,
	}
}

func (origin *ImageRows) overlaps(other *ImageRows) (originY, otherY, length int) {
	var (
		size        = origin.original.Bounds().Dy()
		originSlice = getBigSlice(origin.rows, "x", size)
	)

	for i := 0; i < size*2; i++ {
		var otherSlice = getBigSlice(other.rows, "y", i)

		var mask = getBoolMask(originSlice, otherSlice)

		// now get the longest 'true' sequence
		var maxPos, maxLen = getLongestSequenceOfTrue(mask)

		if maxLen > length {
			length = maxLen
			originY = maxPos - size
			otherY = maxPos - i

			// fmt.Println("Var:", originY, otherY, length)
		}
	}

	return originY, otherY, length
}

func getBoolMask(sl1, sl2 []string) []bool {
	var res = make([]bool, len(sl1))
	for i := range sl1 {
		if sl1[i] == sl2[i] {
			res[i] = true
		}
	}
	return res
}

func getBigSlice(rows []string, base string, offset int) []string {
	var res = make([]string, len(rows)*3)

	for i := 0; i < len(res); i++ {
		if i < offset || i >= offset+len(rows) {
			res[i] = base
		} else {
			res[i] = rows[i-offset]
		}
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

type basicColor struct {
	r, g, b uint32
}

func (bc *basicColor) getString() string {
	return strconv.FormatUint(uint64(bc.r), 10) + strconv.FormatUint(uint64(bc.g), 10) +
		strconv.FormatUint(uint64(bc.b), 10)
}

var basicColors = []basicColor{
	{0, 0, 0},
	{255, 255, 255},
	{255, 0, 0},
	{0, 255, 0},
	{0, 0, 255},
	{255, 255, 0},
	{0, 255, 255},
	{255, 0, 255},
	{192, 192, 192},
	{128, 128, 128},
	{128, 0, 0},
	{128, 128, 0},
	{0, 128, 0},
	{128, 0, 128},
	// {0, 128, 128},
	// {0, 0, 128},
}

func getClosestColor(r, g, b uint32) basicColor {
	var closest = basicColors[0]
	var d = getColorDistance(r, g, b, closest)

	for _, col := range basicColors[1:] {
		if td := getColorDistance(r, g, b, col); td < d {
			closest = col
			d = td
		}
	}

	return closest
}

func getColorDistance(r, g, b uint32, c basicColor) float64 {
	var (
		diffR = float64(r) - float64(c.r)
		diffG = float64(g) - float64(c.g)
		diffB = float64(b) - float64(c.b)
	)
	return math.Sqrt(diffR*diffR + diffG*diffG + diffB*diffB)
}

func (origin *ImageRows) overlaps2(other *ImageRows) (originY, otherY, length int) {
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
