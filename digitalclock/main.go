// +build !solution

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"time"
	"unicode"
)

const (
	widthNumber  int = 8
	heightNumber int = 12
	widthDots    int = 4
)

var numToImg = [...]string{Zero, One, Two, Three, Four, Five, Six, Seven, Eight, Nine}

func main() {
	serverPort := flag.Int("port", 8000, "a port")
	flag.Parse()

	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(fmt.Sprintf("localhost:%d", *serverPort), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Query().Get("k")
	zoom, _ := strconv.Atoi(k)
	if zoom == 0 {
		zoom = 1
	}

	timeString := r.URL.Query().Get("time")
	if timeString == "" {
		timeString = time.Now().Format("15:04:05")
	}

	resultImage, _ := stringToPng(timeString)
	resultImage = scaleImage(resultImage, zoom)

	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, resultImage)

	// TODO return 400
	//w.WriteHeader(http.StatusBadRequest)
}

func stringToPng(s string) (*image.RGBA, error) {
	imageWidth := 2*widthDots + 6*widthNumber
	imageHeight := heightNumber
	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	currentOffset := 0

	for _, char := range s {
		switch {
		case unicode.IsDigit(char):
			curDigit := char - '0'
			sss(rgba, currentOffset, numToImg[curDigit%10], widthNumber)
			currentOffset += widthNumber
		case char == ':':
			sss(rgba, currentOffset, Colon, widthDots)
			currentOffset += widthDots
		default:
			return nil, fmt.Errorf("f")
		}
	}

	return rgba, nil
}

func scaleImage(rgba *image.RGBA, zoom int) *image.RGBA {
	width := rgba.Bounds().Size().X
	height := rgba.Bounds().Size().Y

	newImage := image.NewRGBA(image.Rect(0, 0, width*zoom, height*zoom))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for ii := 0; ii < zoom; ii++ {
				for jj := 0; jj < zoom; jj++ {
					newImage.Set(i*zoom+ii, j*zoom+jj, rgba.At(i, j))
				}
			}
		}
	}

	return newImage
}

func sss(matrix *image.RGBA, offset int, char string, charWidth int) {
	for i := 0; i < heightNumber; i++ {
		for j := 0; j < charWidth; j++ {
			z := char[i*(charWidth+1)+j]
			if z == '1' {
				matrix.Set(offset+j, i, Cyan)
			} else {
				matrix.Set(offset+j, i, color.White)
			}
		}
	}
}
