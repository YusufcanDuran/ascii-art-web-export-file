package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var asciiArt []string
	var group []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		group = append(group, line)
		if len(group) == 9 {
			asciiArt = append(asciiArt, strings.Join(group, "\n"))
			group = []string{}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return asciiArt
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		return
	}
	var asciiArt []string
	style := strings.ToLower(os.Args[2])

	switch style {
	case "standard":
		asciiArt = readFile("standard.txt")
	case "shadow":
		asciiArt = readFile("shadow.txt")
	case "thinkertoy":
		asciiArt = readFile("thinkertoy.txt")
	default:
		fmt.Printf("Invalid style '%s'. Choose 'standard', 'shadow', or 'thinkertoy'.\n", style)
		return
	}

	input := os.Args[1]
	input = strings.ReplaceAll(input, "\\n", "\n")
	PrintAscii(input, style, asciiArt)
}
func PrintAscii(input string, style string, ascii []string) {
	for i, satir := range strings.Split(input, "\n") { // Her kelime için ASCII sanatını yazdır
		if satir == "" {
			if i != 0 {
				fmt.Println()
			}
			continue
		}
		var i int
		if style == "standard" {
			i = 0
		} else if style == "shadow" || style == "thinkertoy" {
			i = 1
		}
		for line := i; line < i+8; line++ {
			for _, char := range satir {
				asciiCode := int(char)
				if asciiCode >= 32 && asciiCode <= 126 {
					asciiIndex := asciiCode - 32
					if asciiIndex < 0 || asciiIndex >= len(ascii) {
						fmt.Printf("No ASCII art representation found for character '%c'\n", char)
						continue
					}
					lines := strings.Split(ascii[asciiIndex], "\n")
					if line < len(lines) {
						fmt.Print(lines[line])
					}
				} else if asciiCode == 13 {
					fmt.Print("")
				} else {
					fmt.Printf("Invalid character: %c\n", char)
				}
			}
			if line != 0 {
				fmt.Println()
			}
		}
	}

}
