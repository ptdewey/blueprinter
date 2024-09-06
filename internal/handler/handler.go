package handler

import (
	"io"
	"os"
)

func CopySelectedItem(src string, dst string) (bool, error) {
	in, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer in.Close()

	out, err := os.Open(dst)
	if err != nil {
		return false, err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return false, err
	}

	if err = out.Sync(); err != nil {
		return false, err
	}

	return true, nil
}
