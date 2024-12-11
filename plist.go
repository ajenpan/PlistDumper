package main

import (
	"fmt"
	"image"

	"howett.net/plist"
)

type FrameV0 struct {
	Height         int     `plist:"height"`
	Width          int     `plist:"width"`
	X              int     `plist:"x"`
	Y              int     `plist:"y"`
	OriginalWidth  int     `plist:"originalWidth"`
	OriginalHeight int     `plist:"originalHeight"`
	OffsetX        float32 `plist:"offsetX"`
	OffsetY        float32 `plist:"offsetY"`
}

type PlistV0 struct {
	Frames map[string]*FrameV0 `plist:"frames"`
}

type FrameV1 struct {
	Frame      string `plist:"frame"`
	Offset     string `plist:"offset"`
	SourceSize string `plist:"sourceSize"`
}

type PlistV1 struct {
	Frames map[string]*FrameV1 `plist:"frames"`
}

type FrameV2 struct {
	Rotated         bool   `plist:"rotated"`
	Frame           string `plist:"frame"`
	Offset          string `plist:"offset"`
	SourceColorRect string `plist:"sourceColorRect"`
	SourceSize      string `plist:"sourceSize"`
}

type PlistV2 struct {
	Frames map[string]*FrameV2 `plist:"frames"`
}

type FrameV3 struct {
	SpriteOffset     string `plist:"spriteOffset"`
	SpriteSize       string `plist:"spriteSize"`
	SpriteSourceSize string `plist:"spriteSourceSize"`
	TextureRect      string `plist:"textureRect"`
	TextureRotated   bool   `plist:"textureRotated"`
}
type PlistV3 struct {
	Frames map[string]*FrameV3 `plist:"frames"`
}

type MetaData struct {
	Format      int    `plist:"format"`
	RealTexture string `plist:"realTextureFileName"`
	Size        string `plist:"size"`
	SmartUpdate string `plist:"smartupdate"`
	Texture     string `plist:"textureFileName"`
}

type Version struct {
	MetaData *MetaData `plist:"metadata"`
}

func dumpPlist(ctx *DumpContext) error {
	version := Version{}
	_, err := plist.Unmarshal(ctx.FileContent, &version)
	if err != nil {
		return err
	}

	if version.MetaData == nil {
		return fmt.Errorf("unmarshal context error, got MetaData nil, filename:%s", ctx.FileName)
	}

	part := ctx.AppendPart()
	part.ImageSoureFile = version.MetaData.Texture

	switch version.MetaData.Format {
	case 0:
		plistData := PlistV0{}
		_, err = plist.Unmarshal(ctx.FileContent, &plistData)
		if err != nil {
			return err
		}

		for k, v := range plistData.Frames {
			part.Frames[k] = &Frame{
				Rect:         image.Rect(v.X, v.Y, v.X+v.Width, v.Y+v.Height),
				OriginalSize: image.Point{v.OriginalWidth, v.OriginalHeight},
				Offset:       image.Point{int(v.OffsetX), int(v.OffsetY)},
				Rotated:      0,
			}
		}
	case 1:
		plistData := PlistV1{}
		_, err = plist.Unmarshal(ctx.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			part.Frames[k] = &Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      0,
			}
		}
	case 2:
		plistData := PlistV2{}
		_, err = plist.Unmarshal(ctx.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.Frame)
			o := intArr(v.Offset)
			s := intArr(v.SourceSize)
			part.Frames[k] = &Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      IfElse(v.Rotated, 90, 0),
			}
		}
	case 3:
		plistData := PlistV3{}
		_, err = plist.Unmarshal(ctx.FileContent, &plistData)
		if err != nil {
			return err
		}
		for k, v := range plistData.Frames {
			f := intArr(v.TextureRect)
			o := intArr(v.SpriteOffset)
			s := intArr(v.SpriteSourceSize)
			part.Frames[k] = &Frame{
				Rect:         image.Rect(f[0], f[1], f[2]+f[0], f[3]+f[1]),
				OriginalSize: image.Point{s[0], s[1]},
				Offset:       image.Point{o[0], o[1]},
				Rotated:      IfElse(v.TextureRotated, 90, 0),
			}
		}
	default:
		return fmt.Errorf("unknow version.MetaData.Format:%d", version.MetaData.Format)
	}

	return nil
}
