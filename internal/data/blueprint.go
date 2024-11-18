package data

import (
	"os"

	"github.com/BurntSushi/toml"
)

type blueprint struct {
	OutputName string `toml:"output_name"`
	// TODO: add any other local config attributes
}

func parseBlueprint(blueprintPath string) (blueprint, error) {
	contents, err := os.ReadFile(blueprintPath)
	if err != nil {
		return blueprint{}, err
	}

	var out blueprint
	err = toml.Unmarshal(contents, &out)
	if err != nil {
		return blueprint{}, err
	}

	return out, nil
}
