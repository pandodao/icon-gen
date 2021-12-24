package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
)

type (
	Params struct {
		outputPath  string
		outputAsset Asset
		firstAsset  Asset
		secondAsset Asset
	}
)

var params Params

func (p *Params) loadAssets(ctx context.Context) error {
	if err := p.outputAsset.Load(ctx); err != nil {
		return err
	}

	if p.firstAsset.AssetID != "" {
		if err := p.firstAsset.Load(ctx); err != nil {
			return err
		}
	}

	if p.secondAsset.AssetID != "" {
		if err := p.secondAsset.Load(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *Params) generate(ctx context.Context, genImageFunc func(ctx context.Context) (*image.RGBA, error)) error {
	if p.outputAsset.AssetID == "" {
		return fmt.Errorf("empty output asset")
	}

	if err := p.loadAssets(ctx); err != nil {
		return err
	}

	folder := path.Join(p.outputPath, p.outputAsset.AssetID)
	os.MkdirAll(folder, os.ModePerm)

	{
		bts, _ := json.MarshalIndent(map[string]interface{}{
			"asset_id": p.outputAsset.AssetID,
			"chain_id": p.outputAsset.Asset.ChainID,
			"cmc_id":   "",
		}, "", "    ")
		if err := ioutil.WriteFile(path.Join(folder, "index.json"), []byte(bts), os.ModePerm); err != nil {
			return err
		}
	}

	image, err := genImageFunc(ctx)
	if err != nil {
		return err
	}

	out, err := os.Create(path.Join(folder, "icon.png"))
	if err != nil {
		return err
	}
	defer out.Close()

	if err = png.Encode(out, image); err != nil {
		return err
	}
	return nil
}
