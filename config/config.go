package config

type (
	Config struct {
		Template Template `json:"template"`
	}

	Template struct {
		Fswap   string `json:"fswap"`
		Rings   string `json:"rings"`
		RingsLP string `json:"rings_lp"`
	}
)
