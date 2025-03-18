package data

import (
	"bytes"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
)

type templateConfig struct {
	TargetTemplate    string   `toml:"target_template"`
	ExtraTemplates    []string `toml:"extra_templates"`
	ExtraDestinations []string `toml:"extra_destinations"`

	PopulateTemplate bool           `toml:"populate_template"`
	TemplateVars     map[string]any `toml:"template_vars"`
}

// `.blueprint.toml` Specification
type blueprint struct {
	OutputName string           `toml:"output_name"`
	Ignore     []string         `toml:"ignore"`
	Extras     []templateConfig `toml:"template_config"`

	// DOC: document hierarchy of vars
	PopulateTemplates bool           `toml:"populate_templates"` // TODO: decide if multiple template flags are necessary
	TemplateVars      map[string]any `toml:"template_vars"`
}

func ExecuteTemplate(tmplPath string, templateVars map[string]any) (*bytes.Buffer, error) {
	contents, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(tmplPath).Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, templateVars); err != nil {
		return nil, err
	}

	return &out, nil
}

func parseBlueprint(blueprintPath string) (blueprint, error) {
	var out blueprint
	_, err := toml.DecodeFile(blueprintPath, &out)
	if err != nil {
		return blueprint{}, err
	}

	return out, nil
}
