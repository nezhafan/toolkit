package toolkit

import (
	"bufio"
	"os"
	"path/filepath"
)

const (
	defaultMaxLineSize = 16 << 20
)

// 文件是否存在（不含权限）
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// 打开文件 (文件或文件夹不存在则创建)
func FileOpen(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if !FileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775)
}

// 按行读取文件
func FileReadByLine(filename string, fn func(string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), defaultMaxLineSize)

	for scanner.Scan() {
		if err := fn(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}
