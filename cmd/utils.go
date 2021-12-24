package cmd

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fox-one/mixin-sdk-go"
)

const (
	SIZE = 520
)

type (
	Asset struct {
		Asset    *mixin.Asset
		AssetID  string
		Icon     string
		ColorHex string
	}
)

func (asset *Asset) Load(ctx context.Context) error {
	if asset.Asset != nil {
		return nil
	}

	if asset.AssetID == "" {
		return fmt.Errorf("empty asset id")
	}
	ass, err := mixin.ReadNetworkAsset(ctx, asset.AssetID)
	if err != nil {
		return err
	}

	asset.Asset = ass
	return nil
}

func (asset *Asset) Image(ctx context.Context) (image.Image, error) {
	var (
		err    error
		reader io.Reader
	)

	if asset.Icon != "" {
		reader, err = os.Open(asset.Icon)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := http.Get(strings.ReplaceAll(asset.Asset.IconURL, "=s128", ""))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		reader = resp.Body
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
