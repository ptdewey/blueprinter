package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopySelectedItem(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Println("Error getting source info: ", err)
		return err
	}

	if srcInfo.IsDir() {
		err := copyDirectory(src, dst)
		if err != nil {
			fmt.Println("Error copying directory: ", err)
			return err
		}
	} else {
		err := copyFile(src, dst)
		if err != nil {
			fmt.Println("Error copying file: ", err)
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// TODO: this currently overwrites the file if it already exists
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func copyDirectory(src, dst string) error {
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

		return copyFile(path, dstPath)
	})

	return err
}
