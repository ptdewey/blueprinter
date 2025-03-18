package handler

import (
	"io"

	"github.com/ptdewey/blueprinter/internal/config"
	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/utils"
)

func handleTemplatePopulation(in *io.Reader, tmplPath string, item data.Item) error {
	// Global level populate options
	cfg := config.Config()
	if !cfg.PopulateTemplates {
		return nil
	}

	templateVars := cfg.TemplateVars

	// Template directory level options
	blueprint := item.Blueprint()
	templateVars = utils.MergeMaps(templateVars, blueprint.TemplateVars)

	// Template file level options
	for _, tc := range blueprint.Extras {
		if tc.TargetTemplate == item.Title() {
			templateVars = utils.MergeMaps(templateVars, tc.TemplateVars)
		}
	}

	if templateVars != nil {
		var err error
		*in, err = data.ExecuteTemplate(tmplPath, templateVars)
		if err != nil {
			return err
		}
	}

	return nil
}
