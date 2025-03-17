package data

import (
	"os"

	"github.com/BurntSushi/toml"
)

type templateConfig struct {
	TargetTemplate    string   `toml:"target_template"`
	ExtraTemplates    []string `toml:"extra_templates"`
	ExtraDestinations []string `toml:"extra_destinations"`
}

// `.blueprint.toml` Specification
type blueprint struct {
	OutputName string           `toml:"output_name"`
	Ignore     []string         `toml:"ignore"`
	Extras     []templateConfig `toml:"template_config"`
	// Add any other local config attributes
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
