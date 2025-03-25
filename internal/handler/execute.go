package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ptdewey/blueprinter/internal/config"
	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/utils"
	"github.com/ptdewey/blueprinter/pkg/flags"
)

// Special reserved template variables
var templateVars = map[string]any{
	"current_year": time.Now().Year(),
	// Expand as necessary
}

func execTemplate(in *io.Reader, tmplPath string, item data.Item) error {
	cfg := config.Config()
	if !cfg.PopulateTemplates {
		return nil
	}

	templateVars["filename"] = item.Title()

	// Global level populate options
	templateVars = utils.MergeMaps(templateVars, cfg.TemplateVars)

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
		dir, _ := os.Getwd()
		templateVars["dir"] = filepath.Base(dir)
		// FIX: parentDir only seems to work some of the time
		templateVars["parentDir"] = filepath.Base(filepath.Dir(dir))

		templateVars["id"] = flags.BlueprinterFlags.ID
		if flags.BlueprinterFlags.ID == "" {
			templateVars["id"] = item.OutputName()
		}
		templateVars["tags"] = flags.BlueprinterFlags.Tags

		year, month, day := time.Now().Date()
		templateVars["date"] = fmt.Sprintf("%s %d, %d", month, day, year)

		var err error
		*in, err = data.ExecuteTemplate(tmplPath, templateVars)
		if err != nil {
			return err
		}
	}

	return nil
}
