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
	"strings"

	color_extractor "github.com/marekm4/color-extractor"
	"github.com/nfnt/resize"
	"github.com/pandodao/icon-gen/templates"
	"github.com/spf13/cobra"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// fswapCmd represents the fswap command
var fswapCmd = &cobra.Command{
	Use:   "fswap",
	Short: "generate 4swap lp token icon",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := params.generate(ctx, func(ctx context.Context) (*image.RGBA, error) {
			return drawFswapLPToken(ctx, params.firstAsset, params.secondAsset)
		}); err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("4swap lp token generated")
	},
}

func drawFswapLPToken(ctx context.Context, asset1, asset2 Asset) (*image.RGBA, error) {
	const (
		ICON_SIZE = 240
		GAP       = 40
	)

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
		color1 := asset1.ColorHex
		color2 := asset2.ColorHex
		if color1 == "" {
			color1 = hexColor(color_extractor.ExtractColors(icon1)[0])
		}
		if color2 == "" {
			color2 = hexColor(color_extractor.ExtractColors(icon2)[0])
		}

		svgData := []byte(strings.ReplaceAll(
			strings.ReplaceAll(templates.FswapLpTemplate, "#DE400B", color1),
			"#19C0FC",
			color2,
		))

		tpl, _ := oksvg.ReadIconStream(bytes.NewReader(svgData))
		tpl.SetTarget(0, 0, float64(SIZE), float64(SIZE))
		tpl.Draw(rasterx.NewDasher(SIZE, SIZE, rasterx.NewScannerGV(SIZE, SIZE, rgba, rgba.Bounds())), 1)
	}

	{
		offset := (SIZE - ICON_SIZE) / 2
		iconSize := icon1.Bounds().Size()
		rect := image.Rectangle{image.Point{offset, 0}, image.Point{offset + iconSize.X, iconSize.Y}}
		draw.Draw(rgba, rect, icon1, image.Point{0, 0}, draw.Over)

		iconSize = icon2.Bounds().Size()
		rect = image.Rectangle{image.Point{offset, SIZE/2 + GAP/2}, image.Point{offset + iconSize.X, SIZE/2 + GAP/2 + iconSize.Y}}
		draw.Draw(rgba, rect, icon2, image.Point{0, 0}, draw.Over)
	}
	return rgba, nil
}

func init() {
	rootCmd.AddCommand(fswapCmd)

	fswapCmd.Flags().StringVar(&params.outputAsset.AssetID, "a", "", "output asset id")
	fswapCmd.Flags().StringVar(&params.outputPath, "o", "outputs", "output path")
	fswapCmd.Flags().StringVar(&params.firstAsset.AssetID, "a1", "", "the first asset id")
	fswapCmd.Flags().StringVar(&params.firstAsset.ColorHex, "c1", "", "the first color (in hex)")
	fswapCmd.Flags().StringVar(&params.firstAsset.Icon, "ic1", "", "the first icon")
	fswapCmd.Flags().StringVar(&params.secondAsset.AssetID, "a2", "", "the second asset id")
	fswapCmd.Flags().StringVar(&params.secondAsset.ColorHex, "c2", "", "the second color (in hex)")
	fswapCmd.Flags().StringVar(&params.secondAsset.Icon, "ic2", "", "the second icon")
}
