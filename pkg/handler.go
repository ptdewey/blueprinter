package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/handler"
)

func GetItems(templateSources []string) []list.Item {
	return data.GetItems(templateSources)
}

func MatchItem(items []list.Item, name string) (*data.Item, error) {
	matches := []data.Item{}
	for _, li := range items {
		i := li.(data.Item)
		if i.Title() == name || i.Path() == name {
			matches = append(matches, i)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("Failed to find matching template with name '%s'\n", name)
	} else if len(matches) == 1 {
		return &matches[0], nil
	}

	// TODO: require more specificity if matches are found?
	return nil, fmt.Errorf("Found more than one matching item with name '%s'.\n%v\n", name, matches)
}

func CopyItem(item *data.Item, dst string, verbose bool) error {
	if dst == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory: ", err)
			return err
		}

		if item.OutputName() == "" {
			dst = filepath.Join(cwd, item.Title())
		} else {
			dst = filepath.Join(cwd, item.OutputName())
		}
	}

	if err := handler.CopySelectedItem(*item, item.Path(), dst); err != nil {
		fmt.Println("Error copying selected item:", err)
		return err
	}

	if verbose {
		fmt.Println(dst)
	}

	return nil
}
