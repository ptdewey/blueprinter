package data

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type Item struct {
	title      string
	desc       string
	path       string
	ext        string
	dirPath    string
	outputName string
	extras     []templateConfig
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	return i.desc
}

// TODO: allow alternative search options (possibly search by file type?)
func (i Item) FilterValue() string {
	return i.title
}

func (i Item) Path() string {
	return i.path
}

func (i Item) Ext() string {
	return i.ext
}

func (i Item) DirPath() string {
	return i.dirPath
}

func (i Item) OutputName() string {
	return i.outputName
}

func (i Item) Extras() []templateConfig {
	return i.extras
}

func GetItems(templateSources []string) []list.Item {
	var out []list.Item
	for _, src := range templateSources {
		items, err := getDirContents(src)
		if err != nil {
			continue
		}
		out = append(out, items...)
	}

	return out
}

func getDirContents(dir string) ([]list.Item, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	blueprint := blueprint{}

	var out []list.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error reading entry info: ", err)
			continue
		}

		ext := filepath.Ext(entry.Name())

		// Read hidden blueprint file in directory if it exists
		if strings.Contains(entry.Name(), ".blueprint.toml") {
			blueprint, err = parseBlueprint(filepath.Join(dir, ".blueprint.toml"))
			if err != nil {
				fmt.Println("Error parsing .blueprint.toml: ", err)
			}
			continue
		}

		var t string
		if entry.IsDir() {
			if strings.HasSuffix(entry.Name(), "blueprints") {
				tempOut, err := getDirContents(filepath.Join(dir, entry.Name()))
				if err != nil {
					fmt.Println("Error reading subdirectory info: ", err)
					continue
				}
				out = append(out, tempOut...)
				continue // Skip adding `...blueprints` directories to output list
			} else {
				t = "dir"
			}
		} else {
			t = ext
		}

		if len(blueprint.Ignore) > 0 {
			if slices.Contains(blueprint.Ignore, entry.Name()) {
				continue
			}
		}

		out = append(out, Item{
			title:      entry.Name(),
			desc:       fmt.Sprintf("Type: %s | Mode: %s | Size: %d bytes\n", t, info.Mode(), info.Size()),
			path:       filepath.Join(dir, entry.Name()),
			ext:        ext,
			dirPath:    dir,
			outputName: blueprint.OutputName,
			extras:     blueprint.Extras,
		})
	}

	return out, nil
}
