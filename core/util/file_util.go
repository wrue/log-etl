package util

import (
	"bufio"
	"os"
)

func ReadToStrArray(filePath string) []string {
	file, _ := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()
	buf := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		line, _, err := buf.ReadLine()
		if err == nil {
			lines = append(lines, string(line))
		} else {
			break
		}
	}
	return lines
}
