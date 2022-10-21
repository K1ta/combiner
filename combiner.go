package combiner

import (
	"errors"
	"image"
	"image/color"
	"math"
)

func Combine(first, second image.Image) (combined image.Image, err error) {
	if first.Bounds().Dx() != second.Bounds().Dy() {
		return nil, errors.New("width must be same")
	}

	y1, y2, length := getOverlap(first, second)

	var res image.Image
	if y1 < y2 {
		res = combineImages(second, first, y2, y1, length)
	} else {
		res = combineImages(first, second, y1, y2, length)
	}

	return res, nil
}

func combineImages(first, second image.Image, firstFromY, secondFromY, length int) image.Image {
	var (
		firstLength   = firstFromY + length
		secondLength  = second.Bounds().Dy() - (secondFromY + length)
		combinedSizeY = firstLength + secondLength
		res           = image.NewRGBA(image.Rect(0, 0, first.Bounds().Dx(), combinedSizeY))
	)

	for x := 0; x < first.Bounds().Dx(); x++ {
		for y := 0; y < firstLength; y++ {
			res.Set(x, y, first.At(x, y))
		}

		for y := 0; y < secondLength; y++ {
			res.Set(x, firstLength+y, second.At(x, secondFromY+length+y))
		}
	}

	return res
}

func getOverlap(i1, i2 image.Image) (y1, y2, length int) {
	var long, short = i1, i2
	if long.Bounds().Dy() < short.Bounds().Dy() {
		long, short = short, long
		defer func() {
			y1, y2 = y2, y1
		}()
	}

	// Для проверки пересечения short изображение последовательно смещается относительно long сначала на высоту short,
	// а затем на высоту long. При сдвиге на каждый пиксель идет сравнение
	var maxOffset = short.Bounds().Dy() + long.Bounds().Dy()

	for offset := 1; offset < maxOffset; offset++ {
		var (
			fromLong  = offset - short.Bounds().Dy()
			toLong    = fromLong + short.Bounds().Dy()
			fromShort = short.Bounds().Dy() - offset
		)
		if fromLong < 0 {
			fromLong = 0
		}
		if toLong > long.Bounds().Dy() {
			toLong = long.Bounds().Dy()
		}
		if fromShort < 0 {
			fromShort = 0
		}

		var mask = make([]bool, toLong-fromLong)

		var noNeedToCheckMask = false
		for x := 0; x < long.Bounds().Dx(); x++ {
			for yOffset := 0; yOffset < toLong-fromLong; yOffset++ {

				var rowValue = true
				if !closeEnough(long.At(x, fromLong+yOffset), short.At(x, fromShort+yOffset)) {
					rowValue = false
				}
				if x > 0 {
					if rowValue != mask[yOffset] {
						mask[yOffset] = false
					}
				} else {
					mask[yOffset] = rowValue
				}
			}
			if x > 0 {
				var _, maxLen = getLongestSequenceOfTrue(mask)
				if maxLen < 100 {
					noNeedToCheckMask = true
					break
				}
			}
		}

		var maxPos, maxLen = 0, 0
		if !noNeedToCheckMask {
			maxPos, maxLen = getLongestSequenceOfTrue(mask)
		}

		if maxLen > length {
			length = maxLen
			y1 = maxPos + fromLong
			y2 = maxPos + fromShort
		}
	}

	return y1, y2, length
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

		diffR = float64(r1>>8) - float64(r2>>8)
		diffG = float64(g1>>8) - float64(g2>>8)
		diffB = float64(b1>>8) - float64(b2>>8)
	)
	res := math.Sqrt(diffR*diffR + diffG*diffG + diffB*diffB)
	return res < 150
}
