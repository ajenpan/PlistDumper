package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
	"strings"
)

/*

btn_ksks.png
size: 304,304
format: RGBA8888
filter: Linear,Linear
repeat: none
00
  rotate: false
  xy: 153, 12
  size: 25, 39
  orig: 140, 40
  offset: 7, 1
  index: -1
*/

// the spine atlas file format
// the frist line is empty
// the second line is sourse image file
// meta data
// frame data
// frame data ...

type SpineAtlasReadStep = int

const (
	statNone       SpineAtlasReadStep = 0
	statAtlasImage SpineAtlasReadStep = 1
	statAtlasMeta  SpineAtlasReadStep = 2
	statFramePart  SpineAtlasReadStep = 3
)

type SpineAtlasPart struct {
	ImageFileName string
	Meta          map[string]string
	Parts         map[string]*FramePart
}

type FramePart struct {
	Name   string
	Rotate string
	XY     string
	Size   string
	Orig   string
	Offset string
	Index  string
}

type SpineAtlas struct {
	Parts map[string]*SpineAtlasPart
}

func dumpSpine(c *DumpContext) error {
	var step SpineAtlasReadStep = 0

	atlas := SpineAtlas{}
	atlas.Parts = map[string]*SpineAtlasPart{}

	atlasPart := &SpineAtlasPart{
		Meta:  make(map[string]string),
		Parts: make(map[string]*FramePart),
	}
	var currFrame *FramePart

	fileScanner := bufio.NewScanner(bytes.NewReader(c.FileContent))
	for fileScanner.Scan() {
		// fmt.Println(fileScanner.Text())
		line := fileScanner.Text()

		switch step {
		case statNone:
			step = statAtlasImage
			if line != "" {
				log.Printf("parse spine file got err, the first line is not empty,line:%v\n", line)
			}
		case statAtlasImage:
			atlasPart.ImageFileName = line
			atlas.Parts[atlasPart.ImageFileName] = atlasPart

			step = statAtlasMeta
		case statAtlasMeta:
			isTitle := !strings.Contains(line, ":")
			if isTitle {
				step = statFramePart

				currFrame = &FramePart{
					Name: strings.TrimSpace(line),
				}
				atlasPart.Parts[currFrame.Name] = currFrame

			} else {
				kv := strings.Split(line, ":")
				if len(kv) == 2 {
					k := strings.TrimSpace(kv[0])
					v := strings.TrimSpace(kv[1])
					atlasPart.Meta[k] = v
				}
			}
		case statFramePart:
			isTitle := !strings.Contains(line, ":")
			if isTitle {
				currFrame = &FramePart{
					Name: strings.TrimSpace(line),
				}
				atlasPart.Parts[currFrame.Name] = currFrame
			} else {
				if currFrame == nil {
					return fmt.Errorf("parse got bug, currFrame is nil")
				}
				line = strings.TrimSpace(line)
				kv := strings.Split(line, ":")
				v := strings.TrimSpace(kv[1])
				switch kv[0] {
				case "rotate":
					currFrame.Rotate = v
				case "xy":
					currFrame.XY = v
				case "size":
					currFrame.Size = v
				case "orig":
					currFrame.Orig = v
				case "offset":
					currFrame.Offset = v
				case "index":
					currFrame.Index = v
				default:
					return errors.New("error parse line " + line)
				}
			}
		default:
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	for _, sp := range atlas.Parts {
		part := c.AppendPart()
		for _, spf := range sp.Parts {

			xy := intArr(spf.XY)
			orig := intArr(spf.Orig)
			offset := intArr(spf.Offset)
			size := intArr(spf.Size)

			if len(xy) != 2 || len(orig) == 0 || len(offset) != 2 || len(size) != 2 {
				fmt.Println(spf.Orig, spf.Offset)
				log.Println("parse to error, param len not 2")
			}

			part.ImageSoureFile = sp.ImageFileName
			part.Frames[spf.Name] = &Frame{
				Rotated:      IfElse(strings.ToLower(spf.Rotate) == "true", 270, 0),
				OriginalSize: image.Pt(orig[0], orig[1]),
				Offset:       image.Pt(offset[0], offset[1]),
				Rect:         image.Rect(xy[0], xy[1], xy[0]+size[0], xy[1]+size[1]),
			}
		}
	}
	return nil
}
