package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"

	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

var highlightedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func search(query string) error {
	lipgloss.SetColorProfile(termenv.TrueColor)

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	lines, err := filter(data, query, isatty.IsTerminal(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

func filter(data []byte, query string, highlight bool) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if query == "" {
		return lines, nil
	}

	matches := fuzzy.Find(query, lines)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found for query: %s", query)
	}

	out := make([]string, 0, len(matches))

	for _, match := range matches {
		var b strings.Builder

		for i := 0; i < len(match.Str); i++ {
			if highlight && slices.Contains(match.MatchedIndexes, i) {
				b.WriteString(highlightedStyle.Render(string(match.Str[i])))
			} else {
				b.WriteString(string(match.Str[i]))
			}
		}

		out = append(out, b.String())
	}

	return out, nil
}
