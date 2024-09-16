package data

import (
	"encoding/json"
	"os"
)

type blueprint struct {
	OutputName string `json:"outputName"`
	// TODO: add any other local config attributes
}

func parseBlueprint(blueprintPath string) (blueprint, error) {
	contents, err := os.ReadFile(blueprintPath)
	if err != nil {
		return blueprint{}, err
	}

	var out blueprint
	err = json.Unmarshal(contents, &out)
	if err != nil {
		return blueprint{}, err
	}

	return out, nil
}
