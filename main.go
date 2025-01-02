package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {
	width, err := getTerminalWidth()
	if err != nil {
		fmt.Println("Error detecting terminal width: ", err)
	}
	fmt.Println("W:", width)
}

func getTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	return width, err
}

func colorPrint(text string, color string) {
	color = strings.ToLower(color)
	colors := map[string]string{
		"black": "30",
		"red": "31",
		"green": "32",
		"yellow": "33",
		"blue": "34",
		"magenta": "35",
		"cyan": "36",
		"white": "37",
	}
	code, ok := colors[color]
	if ok {
		fmt.Printf("\033[%sm%s\033[0m", code, text)
	} else {
		fmt.Print(text)
	}
}
