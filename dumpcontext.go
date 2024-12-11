package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type Frame struct {
	Rect         image.Rectangle
	Offset       image.Point
	OriginalSize image.Point
	Rotated      int
}

type AtlasPart struct {
	ImageSoureFile string
	Frames         map[string]*Frame
}

type DumpContext struct {
	FileName    string
	FileContent []byte
	Atlases     []*AtlasPart
}

func (dc *DumpContext) AppendPart() *AtlasPart {
	part := &AtlasPart{}
	part.Frames = map[string]*Frame{}
	dc.Atlases = append(dc.Atlases, part)
	return part
}

func (dc *DumpContext) Dump() error {
	var err error

	subdir := filepath.Base(dc.FileName)
	subdir = strings.TrimSuffix(subdir, filepath.Ext(subdir))
	fmt.Println(subdir)

	outdir := dc.FileName + ".out"

	if !IsDir(outdir) {
		err = os.MkdirAll(outdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, part := range dc.Atlases {
		textureFileName := filepath.Join(filepath.Dir(dc.FileName), part.ImageSoureFile)
		textureImage, err := LoadImage(textureFileName)
		if err != nil {
			return fmt.Errorf("open image error:" + textureFileName)
		}
		err = part.dumpFrames(outdir, textureImage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (part *AtlasPart) dumpFrames(outdir string, textureImage image.Image) error {
	for filename, frame := range part.Frames {
		var subImage image.Image

		w, h := frame.Rect.Size().X, frame.Rect.Size().Y
		ox, oy := frame.Offset.X, frame.Offset.Y
		ow, oh := frame.OriginalSize.X, frame.OriginalSize.Y
		x, y := frame.Rect.Min.X, frame.Rect.Min.Y

		if frame.Rotated == 90 {
			subImage = imaging.Crop(textureImage, image.Rect(x, y, x+h, y+w))
			subImage = imaging.Rotate90(subImage)
		} else if frame.Rotated == 270 {
			subImage = imaging.Crop(textureImage, image.Rect(x, y, x+h, y+w))
			subImage = imaging.Rotate270(subImage)
		} else {
			subImage = imaging.Crop(textureImage, image.Rect(x, y, x+w, y+h))
		}

		destImage := image.NewRGBA(image.Rect(0, 0, ow, oh))
		newImage := imaging.Paste(destImage, subImage, image.Point{(ow-w)/2 + ox, (oh-h)/2 - oy})

		if TrimDir {
			filename = filepath.Base(filename)
			filename = strings.ReplaceAll(filename, " ", "")
		}

		savepath := filepath.Join(outdir, filename)

		if filepath.Ext(savepath) == "" {
			savepath += ".png"
		}

		SaveImage(savepath, newImage)
	}
	return nil
}

func dumpByFileName(filename string) error {
	c := DumpContext{
		FileName: filename,
		Atlases:  []*AtlasPart{},
	}
	var err error

	data, err := os.ReadFile(c.FileName)
	if err != nil {
		return fmt.Errorf("read file err %v", err)
	}

	c.FileContent = data

	ext := path.Ext(filename)
	switch ext {
	case ".plist":
		err = dumpPlist(&c)
	case ".json":
		err = dumpJson(&c)
	case ".fnt":
		err = dumpFnt(&c)
	case ".atlas":
		err = dumpSpine(&c)
	default:
		err = errors.New("not support file type")
	}

	if err != nil {
		return err
	}
	return c.Dump()
}
