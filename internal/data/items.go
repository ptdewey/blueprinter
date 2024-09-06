package data

import (
	"blueprinter/internal/ui"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

// TODO: allow alternative search options (possibly search by file type?w)

// TODO: get items from target directory (search subdirectories?)
// - allow both directories and files to be targeted by copy operation
// - allow creating a project from an entire subdir, but also allow searching subdirs
//   - if only one is possible, go with copy entire dir (for presentations and template projects)
//   - if both are possible, probably make subdir search opt-in for certain dirs (i.e. python, R dirs, whereas go would want entire dir copied)
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

func GetItems() []list.Item {
	out, err := getDirContents("/home/patrick/dotfiles/templates")
	if err != nil {
		return nil
	}

	return out
}
