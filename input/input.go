package input

import (
	"bufio"
	"os"
	"strings"
)

func Simple() [][]string {
	var in [][]string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		in = append(in, strings.Fields(line))
	}
	return in
}
