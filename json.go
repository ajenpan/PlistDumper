package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
)

// var (
// 	ErrNotSupportJsonType = errors.New("not support json type")
// 	ErrNotSupportFileType = errors.New("not support file type")
// )

var ErrSkip = errors.New("skip")

type JsonSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type JsonRect struct {
	W int `json:"w"`
	H int `json:"h"`
	X int `json:"x"`
	Y int `json:"y"`
}

type JsonMetaData struct {
	Image   string `json:"image"`
	Version string `json:"version"`
}

type JsonVersion struct {
	Meta   *JsonMetaData `json:"meta"`
	Frames interface{}   `json:"frames"`
}

type JsonFrameHashV1 struct {
	Frames map[string]*JsonFrameV1 `json:"frames"`
}

type JsonFrameArrayV1 struct {
	Frames []*JsonFrameV1 `json:"frames"`
}

type JsonFrameV1 struct {
	Frame            *JsonRect `json:"frame"`
	Rotated          bool      `json:"rotated"`
	Trimmed          bool      `json:"trimmed"`
	SpriteSourceSize *JsonRect `json:"spriteSourceSize"`
	SourceSize       *JsonSize `json:"sourceSize"`
	Filename         string    `json:"filename"`
}

func dumpJson(ctx *DumpContext) error {

	version := JsonVersion{}
	err := json.Unmarshal(ctx.FileContent, &version)
	if err != nil {
		return ErrSkip
	}

	if version.Meta == nil {
		return ErrSkip
	}

	if version.Meta.Version != "1.0" {
		return fmt.Errorf("unknow version:%s", version.Meta.Version)
	}

	part := ctx.AppendPart()
	part.ImageSoureFile = version.Meta.Image

	frames := map[string]*JsonFrameV1{}

	switch version.Frames.(type) {
	case map[string]interface{}:
		jsonData := JsonFrameHashV1{}
		err = json.Unmarshal(ctx.FileContent, &jsonData)
		if err != nil {
			return err
		}
		frames = jsonData.Frames
	case []interface{}:
		jsonData := JsonFrameArrayV1{}
		err = json.Unmarshal(ctx.FileContent, &jsonData)
		if err != nil {
			return err
		}
		for _, v := range jsonData.Frames {
			frames[v.Filename] = v
		}
	default:
		return errors.New("unknow version:[" + version.Meta.Version + "]")
	}

	for k, v := range frames {
		f := v.Frame
		s := v.SourceSize
		part.Frames[k] = &Frame{
			Rect:         image.Rect(f.X, f.Y, f.X+f.W, f.Y+f.H),
			OriginalSize: image.Point{s.W, s.H},
			Rotated:      IfElse(v.Rotated, 90, 0),
			Offset:       image.Point{-v.SpriteSourceSize.X / 2, -v.SpriteSourceSize.Y / 2}, //plist offset in center, json in left-top
		}
	}

	return nil
}
