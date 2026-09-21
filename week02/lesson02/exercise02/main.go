package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func readFileSafely(path string) (int, error) {
	data, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("đọc file %q thất bại: %w", path, err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("Scan file %q that bai: %w", path, err)
	}
	return lineNumber, err
}

func main() {
	path := "./week02/lesson02/exercise02/hello.txt"
	lineNumber, err := readFileSafely(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("⚠️ File %q chua duoc tao\n", path)
		} else {
			fmt.Println("Loi: ", err)
		}
		return
	}
	fmt.Printf("So dong trong file %q: %d\n", path, lineNumber)
}
