package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type Item struct {
	title string
	desc  string
	path  string
	ext   string
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

func GetItems(templateSources []string) []list.Item {
	var out []list.Item
	for _, src := range templateSources {
		items, err := getDirContents(src)
		if err != nil {
			return nil
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

	var out []list.Item
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error reading entry info: ", err)
			continue
		}

		ext := filepath.Ext(entry.Name())

		// TODO: allow some sort of autonaming of output determined from some config info/file? (i.e. for gitignore-blueprints, save as .gitignore)
		var t string
		if entry.IsDir() {
			if strings.HasSuffix(entry.Name(), "blueprints") {
				tempOut, err := getDirContents(filepath.Join(dir, entry.Name()))
				if err != nil {
					fmt.Println("Error reading subdirectory info: ", err)
					continue
				}
				out = append(out, tempOut...)
				continue // skip adding '...blueprints' directories to output list
			} else {
				t = "dir"
			}
		} else {
			t = ext
		}

		out = append(out, Item{
			title: entry.Name(),
			desc:  fmt.Sprintf("Type: %s | Mode: %s | Size: %d bytes\n", t, info.Mode(), info.Size()),
			path:  filepath.Join(dir, entry.Name()),
			ext:   ext,
		})
	}

	return out, nil
}
