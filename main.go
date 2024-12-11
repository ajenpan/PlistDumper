package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

var Input string = "./"
var Output string = "./"

var Ext string = "json,plist,fnt,atlas"
var TrimDir bool = false

func main() {
	cmd := &cli.Command{
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "input",
				Value:       "./",
				Max:         1,
				Destination: &Input,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ext",
				Aliases:     []string{"e"},
				Value:       "json,plist,fnt,atlas",
				Destination: &Ext,
				Config:      cli.StringConfig{TrimSpace: true},
			}, &cli.BoolFlag{
				Name:        "trimdir",
				Destination: &TrimDir,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			var ext = []string{}
			arr := strings.Split(Ext, ",")
			for _, v := range arr {
				ext = append(ext, "."+v)
			}

			if Input == "" {
				Input = "./"
			}

			filenames := []string{}
			if IsDir(Input) {
				files := GetAllFiles(Input, ext)
				filenames = append(filenames, files...)
			} else {
				filenames = append(filenames, Input)
			}

			for i, filename := range filenames {
				p := fmt.Sprintf("[%d/%d]", i+1, len(filenames))
				log.Println("开始导出", p, filename)
				err := dumpByFileName(filename)
				if err != nil {
					if !errors.Is(err, ErrSkip) {
						log.Println("错误", p, filename, err)
					}
				}
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
