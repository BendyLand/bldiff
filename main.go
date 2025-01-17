package main

import (
	"fmt"
	"golang.org/x/term"
	"math"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		file1, file2 := getBothFileContents(os.Args[1], os.Args[2])
		file1 = addLineNumbers(file1)
		file2 = addLineNumbers(file2)
		width, err := getTerminalWidth()
		longest := getLongestFileLength(file1, file2)
		if err != nil {
			fmt.Println("Error detecting terminal width:", err)
			os.Exit(1)
		}
		width = int(float64(width) * 0.45)
		file1 = normalizeFileLength(file1, width, longest)
		file2 = normalizeFileLength(file2, width, longest)
		printHalves(file1, file2)
	} else {
		fmt.Println("Usage: bldiff <file1> <file2>")
	}
}

func getTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}

func getLongestFileLength(file1 string, file2 string) int {
	lines1 := strings.Split(file1, "\n")
	lines2 := strings.Split(file2, "\n")
	return int(math.Max(float64(len(lines1)), float64(len(lines2))))

}

func colorPrint(text string, color string) {
	color = strings.ToLower(color)
	colors := map[string]string{
		"black":   "30",
		"red":     "31",
		"green":   "32",
		"yellow":  "33",
		"blue":    "34",
		"magenta": "35",
		"cyan":    "36",
		"white":   "37",
	}
	code, ok := colors[color]
	if ok {
		fmt.Printf("\033[%sm%s\033[0m", code, text)
	} else {
		fmt.Print(text)
	}
}

func getFileContents(path string) (string, error) {
	text, err := os.ReadFile(path)
	return string(text), err
}

func getBothFileContents(path1 string, path2 string) (string, string) {
	text1, err1 := getFileContents(path1)
	if err1 != nil {
		fmt.Printf("Error reading file '%s': %v", path1, err1)
		os.Exit(1)
	}
	text2, err2 := getFileContents(path2)
	if err2 != nil {
		fmt.Printf("Error reading file '%s': %v", path2, err2)
		os.Exit(1)
	}
	return text1, text2
}

func addLineNumbers(file string) string {
	result := ""
	lines := strings.Split(file, "\n")
	length := len(lines)
	padLen := len(fmt.Sprintf("%d", length)) - 1
	padding := " "
	for range padLen {
		padding += " "
	}
	for i, line := range lines {
		result += fmt.Sprintf("%d%s%s\n", i+1, padding, line)
	}
	return result
}

func normalizeFileLength(file string, maxWidth int, minLines int) string {
	lines := strings.Split(file, "\n")
	parts := make([]string, len(lines))
	for i, line := range lines {
		if len(line) > maxWidth {
			parts[i] = line[:maxWidth]
		} else {
			for len(line) < maxWidth {
				line += " "
			}
			parts[i] = line
		}
		// this is for the line directly after the numbers
		if i == len(lines)-1 {
			temp := "\033[31mX\033[0m "
			// the color code is an extra 9 chars
			for len(temp)-9 < maxWidth {
				temp += " "
			}
			parts[i] = temp
		}
	}
	// this is for the rest of the missing lines
	for len(parts) < minLines {
		temp := "\033[31mX\033[0m "
		for len(temp)-9 < maxWidth {
			temp += " "
		}
		parts = append(parts, temp)
	}
	return strings.Join(parts, "\n")
}

func printHalves(file1 string, file2 string) {
	lines1 := strings.Split(file1, "\n")
	lines2 := strings.Split(file2, "\n")
	length := len(lines1) - 1
	fmt.Println("")
	for i := range length {
		temp1 := strings.Trim(extractLineContents(lines1[i]), " ")
		temp2 := strings.Trim(extractLineContents(lines2[i]), " ")
		extraLine := strings.Contains(temp1, "X") || strings.Contains(temp2, "X")
		switch {
		case temp1 == temp2:
			colorPrint(lines1[i], "green")
			colorPrint(" | ", "green")
			colorPrint(lines2[i], "green")
			fmt.Println("")
		case checkSimilar(temp1, temp2) && !extraLine:
			colorPrint(lines1[i], "yellow")
			colorPrint(" | ", "yellow")
			colorPrint(lines2[i], "yellow")
			fmt.Println("")
		default:
			colorPrint(lines1[i], "red")
			colorPrint(" | ", "red")
			colorPrint(lines2[i], "red")
			fmt.Println("")
		}
	}
	fmt.Println("")
}

// This function will extract the contents of a line past the number, 
// unless it starts with an 'X', which will be prepended to the line.
func extractLineContents(s string) string {
	result := ""
	start := false
	for _, c := range s {
		if !start && c == 'X' {
			result += string('X')
			continue
		}
		if start {
			result += string(c)
		}
		if c == ' ' {
			start = true
		}
	}
	return result
}

func checkSimilar(line1 string, line2 string) bool {
	len1 := len(line1)
	len2 := len(line2)
	// empty lines shouldn't match literally everything
	if (len1 == 0 || len2 == 0) && len1 != len2 {
		return false
	}
	switch {
	case strings.Contains(line1, line2):
		return true
	case strings.Contains(line2, line1):
		return true
	default:
		return false
	}
}
