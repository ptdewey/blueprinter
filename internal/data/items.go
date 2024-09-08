package data

import (
	"blueprinter/internal/ui"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

// TODO: allow alternative search options (possibly search by file type?w)
func getDirContents(dir string) ([]list.Item, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	out := make([]list.Item, len(entries))
	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error reading entry info: ", err)
			continue
		}

		var t string
		if entry.IsDir() {
			t = "dir "
		} else {
			t = "file"
		}

		out[i] = ui.Item{
			Name: entry.Name(),
			Desc: fmt.Sprintf("Type: %s | Size: %d bytes | Mode: %s \n", t, info.Size(), info.Mode()),
		}
	}

	return out, nil
}

func GetItems(templateSources []string) []list.Item {
	var out []list.Item
	for _, src := range templateSources {
		items, err := getDirContents(src)
		if err != nil {
			return nil
		}
		for _, item := range items {
			out = append(out, item)
		}
	}

	return out
}
