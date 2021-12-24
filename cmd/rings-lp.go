/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"io/ioutil"
	"strings"

	color_extractor "github.com/marekm4/color-extractor"
	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// ringsLpCmd represents the rings command
var ringsLpCmd = &cobra.Command{
	Use:   "lp",
	Short: "generate rings lp token icon",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if err := params.generate(ctx, func(ctx context.Context) (*image.RGBA, error) {
			return drawRingsLPToken(ctx, params.firstAsset, params.secondAsset)
		}); err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("rings lp token generated")
	},
}

func drawRingsLPToken(ctx context.Context, asset1, asset2 Asset) (*image.RGBA, error) {
	const ICON_SIZE = 230

	var (
		rgba = image.NewRGBA(image.Rect(0, 0, SIZE, SIZE))

		icon1 image.Image
		icon2 image.Image
	)

	{
		icon, err := asset1.Image(ctx)
		if err != nil {
			return nil, err
		}

		icon1 = resize.Resize(ICON_SIZE, ICON_SIZE, icon, resize.Bilinear)
	}

	{
		icon, err := asset2.Image(ctx)
		if err != nil {
			return nil, err
		}

		icon2 = resize.Resize(ICON_SIZE, ICON_SIZE, icon, resize.Bilinear)
	}

	{
		svgData, err := ioutil.ReadFile(cfg.Template.RingsLP)
		if err != nil {
			return nil, err
		}

		color1 := asset1.ColorHex
		color2 := asset2.ColorHex
		if color1 == "" {
			color1 = hexColor(color_extractor.ExtractColors(icon1)[0])
		}
		if color2 == "" {
			color2 = hexColor(color_extractor.ExtractColors(icon2)[0])
		}

		svgData = []byte(strings.ReplaceAll(
			strings.ReplaceAll(string(svgData), "#F7931B", color1),
			"#2277EF",
			color2,
		))

		tpl, _ := oksvg.ReadIconStream(bytes.NewReader(svgData))
		tpl.SetTarget(0, 0, float64(SIZE), float64(SIZE))
		tpl.Draw(rasterx.NewDasher(SIZE, SIZE, rasterx.NewScannerGV(SIZE, SIZE, rgba, rgba.Bounds())), 1)
	}

	{
		offset := image.Point{59, 145}
		rect := image.Rectangle{offset, offset.Add(icon1.Bounds().Size())}
		draw.Draw(rgba, rect, icon1, image.Point{0, 0}, draw.Over)

		offset = image.Point{SIZE - 59 - icon1.Bounds().Size().X, 145}
		rect = image.Rectangle{offset, offset.Add(icon2.Bounds().Size())}
		draw.Draw(rgba, rect, icon2, image.Point{0, 0}, draw.Over)
	}
	return rgba, nil
}

func init() {
	ringsCmd.AddCommand(ringsLpCmd)

	ringsLpCmd.Flags().StringVar(&params.outputAsset.AssetID, "a", "", "output asset id")
	ringsLpCmd.Flags().StringVar(&params.outputPath, "o", "outputs", "output path")
	ringsLpCmd.Flags().StringVar(&params.firstAsset.AssetID, "a1", "", "the first asset id")
	ringsLpCmd.Flags().StringVar(&params.firstAsset.ColorHex, "c1", "", "the first color (in hex)")
	ringsLpCmd.Flags().StringVar(&params.firstAsset.Icon, "ic1", "", "the first icon")
	ringsLpCmd.Flags().StringVar(&params.secondAsset.AssetID, "a2", "", "the second asset id")
	ringsLpCmd.Flags().StringVar(&params.secondAsset.ColorHex, "c2", "", "the second color (in hex)")
	ringsLpCmd.Flags().StringVar(&params.secondAsset.Icon, "ic2", "", "the second icon")
}
