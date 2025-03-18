package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ptdewey/blueprinter/internal/data"
)

func CopySelectedItem(src string, dst string, item data.Item) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Println("Error getting source info: ", err)
		return err
	}

	if srcInfo.IsDir() {
		err := copyDirectory(src, dst, item)
		if err != nil {
			fmt.Println("Error copying directory: ", err)
			return err
		}
		return nil
	}

	err = copyFile(src, dst, item)
	if err != nil {
		fmt.Println("Error copying file: ", err)
		return err
	}

	return nil
}

func copyFile(src, dst string, item data.Item) error {
	var in io.Reader
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.(*os.File).Close()

	// TODO: this currently overwrites the file if it already exists
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Call template population handler
	handleTemplatePopulation(&in, src, item)

	if _, err = io.Copy(out, in); err != nil {
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

		return copyFile(path, dstPath, item)
	})

	return err
}
