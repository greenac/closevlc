package main

import (
	"os/exec"
	"fmt"
	"strings"
)

func getLines(data []byte) [][]byte {
	lines := make([][]byte, 0)
	var currentLine *[]byte
	for _, b := range data {
		if currentLine == nil {
			l := make([]byte, 0)
			currentLine = &l
		}

		if b == '\n' {
			lines = append(lines, *currentLine)
			currentLine = nil
		} else {
			*currentLine = append(*currentLine, b)
		}
	}

	return lines
}

func vlcLines(lines *[][]byte) [][]byte {
	const target = "/Applications/VLC.app/Contents/MacOS/VLC"
	vlcl := make([][]byte, 0)
	for _, line := range *lines {
		if strings.Contains(string(line), target) {
			vlcl = append(vlcl, line)
		}
	}

	return vlcl
}

func main()  {
	cmd := exec.Command("ps", "ax")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to make command: ps -ax error:", err)
		return
	}

	lines := getLines(output)
	targetLines := vlcLines(&lines)
	for i, l := range targetLines {
		cleanedLine := make([]rune, 0, len(l))
		for j, value := range string(l) {
			if j == 0 && value == ' ' {
				continue
			}

			cleanedLine = append(cleanedLine, value)
		}

		targetLines[i] = []byte(string(cleanedLine))
	}

	pids := make([][]byte, len(targetLines))
	for i, line := range targetLines {
		pid := make([]byte, 0)
		for _, value := range line {
			if value == ' ' {
				break
			}

			pid = append(pid, value)
		}

		pids[i] = pid
	}

	for _, pid := range pids {
		pidStr := string(pid)
		cmd := exec.Command("kill", "-9", pidStr)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Failed to kill:", pidStr, "error:", err)
			continue
		}

		fmt.Println("killed process:", pidStr)
	}
}
