package data

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

// TODO: add file extension (ext) to item (to allow filtering by ft?)
type Item struct {
	title string
	desc  string
	path  string
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	return i.desc
}

func (i Item) FilterValue() string {
	return i.title
}

func (i Item) Path() string {
	return i.path
}

// TODO: allow alternative search options (possibly search by file type?)
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
			t = "d"
		} else {
			// TODO: show file extension rather than 'f'?
			t = "f"
		}

		out[i] = Item{
			title: entry.Name(),
			desc:  fmt.Sprintf("Type: %s | Mode: %s | Size: %d bytes\n", t, info.Mode(), info.Size()),
			path:  filepath.Join(dir, entry.Name()),
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
