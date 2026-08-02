package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestFilter(t *testing.T) {
	data := []byte("golang\npython\nrust\n")

	got, err := filter(data, "go", false)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(got, "golang") {
		t.Errorf("expected golang in results, got %v", got)
	}

	if _, err := filter(data, "zzz", false); err == nil {
		t.Error("expected error when no matches found")
	}

	all, err := filter(data, "", false)
	if err != nil {
		t.Fatal(err)
	}

	if len(all) != 3 {
		t.Errorf("expected all 3 lines, got %v", all)
	}
}

func TestFilterHighlight(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	got, err := filter([]byte("golang\n"), "go", true)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got[0], "\x1b[") {
		t.Errorf("expected ANSI highlight codes, got %q", got[0])
	}

	plain, err := filter([]byte("golang\n"), "go", false)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(plain[0], "\x1b[") {
		t.Errorf("expected no ANSI codes when not highlighting, got %q", plain[0])
	}
}
