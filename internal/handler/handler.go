package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ptdewey/blueprinter/internal/data"
)

func CopySelectedItem(item data.Item, src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Println("Error getting source info: ", err)
		return err
	}

	if srcInfo.IsDir() {
		if err := copyDirectory(src, dst, item); err != nil {
			fmt.Println("Error copying directory: ", err)
			return err
		}
		return nil
	}

	if err := copyFile(item, src, dst); err != nil {
		fmt.Println("Error copying file: ", err)
		return err
	}

	if err := copyExtraItems(item); err != nil {
		return err
	}

	return nil
}

func copyExtraItems(item data.Item) error {
	for _, et := range item.Blueprint().Extras {
		if et.TargetTemplate != item.Title() {
			continue
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory: ", err)
			return err
		}

		for i, t := range et.ExtraTemplates {
			var dst string
			if len(et.ExtraDestinations) > i && et.ExtraDestinations[i] != "" {
				dst = filepath.Join(cwd, et.ExtraDestinations[i])
			} else {
				dst = filepath.Join(cwd, t)
			}

			src := filepath.Join(item.DirPath(), t)
			if err := copyFile(item, src, dst); err != nil {
				fmt.Println("Error copying additional template files for selected item:", err)
				return err
			}
		}
	}
	return nil
}

func copyFile(item data.Item, src string, dst string) error {
	var in io.Reader
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.(*os.File).Close()

	// FIX: this currently overwrites the file if it already exists
	// - Should it stop and warn the user? Config option?
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := execTemplate(&in, src, item); err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func copyDirectory(src, dst string, item data.Item) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(item, path, dstPath)
	})

	return err
}
