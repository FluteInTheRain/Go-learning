package main

import (
	"errors"
	"fmt"
	"os"
)

func readFileSafely(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("đọc file %q thất bại: %w", path, err)
	}
	return string(data), nil
}

func main() {
	path := "./week02/lesson02/exercise01/hello.txt"
	content, err := readFileSafely(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("⚠️ File %q chua duoc tao\n", path)
			os.Exit(1)
		} else {
			fmt.Println("Loi: ", err)
		}
	}
	fmt.Printf("Đọc thành công %q:\n%s\n", path, content)
}
