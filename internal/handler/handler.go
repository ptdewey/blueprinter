package handler

import (
	"fmt"
	"io"
	"os"
)

// FIX: does not currently work with directories
func CopySelectedItem(src string, dst string) (bool, error) {
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("Error opening source file: ", err)
		return false, err
	}
	defer in.Close()

	// TODO: this currently overwrites the file if it already exists
	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("Error creating destination file: ", err)
		return false, err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		fmt.Println("Error copying data: ", err)
		return false, err
	}

	if err = out.Sync(); err != nil {
		fmt.Println("Error syncing output: ", err)
		return false, err
	}

	return true, nil
}
