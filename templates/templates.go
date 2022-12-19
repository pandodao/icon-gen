package templates

import _ "embed"

var (
	//go:embed rings.png
	RingsTemplate string

	//go:embed rings-lp.svg
	RingsLpTemplate string

	//go:embed fswap-lp.svg
	FswapLpTemplate string
)
