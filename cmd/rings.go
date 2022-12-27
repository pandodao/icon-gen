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
	"context"
	_ "embed"
	"image"
	"image/draw"
	"strings"

	"github.com/nfnt/resize"
	"github.com/pandodao/icon-gen/templates"
	"github.com/spf13/cobra"
)

// ringsCmd represents the rings command
var ringsCmd = &cobra.Command{
	Use:   "rings",
	Short: "generate rings token icons",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if err := params.generate(ctx, func(ctx context.Context) (*image.RGBA, error) {
			return drawRingsToken(ctx, params.firstAsset)
		}); err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("rings token generated")
	},
}

func drawRingsToken(ctx context.Context, asset Asset) (*image.RGBA, error) {
	const ICON_SIZE = 472

	rgba := image.NewRGBA(image.Rect(0, 0, SIZE, SIZE))

	{
		icon, err := asset.Image(ctx)
		if err != nil {
			return nil, err
		}

		icon = resize.Resize(ICON_SIZE, ICON_SIZE, icon, resize.Bilinear)
		iconSize := icon.Bounds().Size()

		offset := image.Point{(SIZE - ICON_SIZE) / 2, (SIZE - ICON_SIZE) / 2}
		draw.Draw(
			rgba,
			image.Rectangle{offset, offset.Add(iconSize)},
			icon,
			image.Point{},
			draw.Over,
		)
	}

	{
		mask, _, err := image.Decode(strings.NewReader(templates.RingsTemplate))
		if err != nil {
			return nil, err
		}
		draw.Draw(rgba, rgba.Bounds(), mask, image.Point{0, 0}, draw.Over)
	}
	return rgba, nil
}

func init() {
	rootCmd.AddCommand(ringsCmd)

	ringsCmd.Flags().StringVar(&params.outputAsset.AssetID, "a", "", "output asset id")
	ringsCmd.Flags().StringVar(&params.outputPath, "o", "outputs", "output path")
	ringsCmd.Flags().StringVar(&params.firstAsset.AssetID, "a1", "", "the first asset id")
	ringsCmd.Flags().StringVar(&params.firstAsset.ColorHex, "c1", "", "the first color (in hex)")
	ringsCmd.Flags().StringVar(&params.firstAsset.Icon, "ic1", "", "the first icon")
}
