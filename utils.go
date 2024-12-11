package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

// condi ? a : b
func IfElse[T any](condi bool, a, b T) T {
	if condi {
		return a
	}
	return b
}

func LoadImage(path string) (img image.Image, err error) {
	return imaging.Open(path)
}

func SaveImage(path string, img image.Image) (err error) {

	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func GetAllFiles(dir string, allow []string) []string {
	allowMap := map[string]bool{}
	for _, v := range allow {
		allowMap[v] = true
	}

	ret := []string{}
	filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		ext := path.Ext(fpath)
		if allowMap[ext] {
			ret = append(ret, filepath.ToSlash(fpath))
		}

		return nil
	})

	return ret
}

func subImage(src image.Image, x, y, w, h int) image.Image {
	r := image.Rect(0, 0, x+w, y+h)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dst, r, src, image.Point{x, y}, draw.Src)
	return dst
}

func rotateImage(src image.Image) image.Image {
	w := src.Bounds().Max.X
	h := src.Bounds().Max.Y
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dst.Set(y, w-x, src.At(x, y))
		}
	}
	return dst
}

func intArr(str string) []int {
	s := strings.Replace(str, "{", "", -1)
	s = strings.Replace(s, "}", "", -1)
	s = strings.Replace(s, " ", "", -1)

	sA := strings.Split(s, ",")

	ret := make([]int, len(sA))

	for i, v := range sA {
		if len(v) == 0 {
			ret[i] = 0
			continue
		}
		value, err := strconv.ParseFloat(v, 32)
		if err != nil {
			value, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				panic(err)
			}
			ret[i] = int(value)
		} else {
			ret[i] = int(value)
		}
	}

	return ret
}
