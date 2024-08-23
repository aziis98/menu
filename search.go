package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/sahilm/fuzzy"
)

var highlightedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func search(query string) error {
	lipgloss.SetColorProfile(termenv.TrueColor)

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if query == "" {
		for _, line := range lines {
			fmt.Println(line)
		}

		return nil
	}

	matches := fuzzy.Find(query, lines)
	if len(matches) == 0 {
		return fmt.Errorf("no matches found for query: %s", query)
	}

	for _, match := range matches {
		for i := 0; i < len(match.Str); i++ {
			if slices.Contains(match.MatchedIndexes, i) {
				fmt.Print(highlightedStyle.Render(string(match.Str[i])))
			} else {
				fmt.Print(string(match.Str[i]))
			}
		}

		fmt.Println()
	}

	return nil
}
