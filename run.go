package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const debounceDuration = 30 * time.Millisecond

type commandOutputMsg struct {
	seq   int
	lines []string
}

type commandErrorMsg struct {
	seq int
	err error
}

type commandCloseMsg struct {
	output string
}

type debounceTickMsg struct {
	seq   int
	event string
}

func (m model) debounce(event string, seq int) tea.Cmd {
	return tea.Tick(debounceDuration, func(t time.Time) tea.Msg {
		return debounceTickMsg{seq: seq, event: event}
	})
}

func (m model) runCommand(event string, seq int) tea.Cmd {
	return func() tea.Msg {
		selectedLine := ""
		if m.selected >= 0 && m.selected < len(m.lines) {
			selectedLine = strings.TrimSpace(m.lines[m.selected])
		}

		cmd := exec.Command("sh", "-c", strings.Join(
			[]string{
				fmt.Sprintf("export prompt=%q", m.input.Value()),
				fmt.Sprintf("export event=%q", event),
				fmt.Sprintf("export sel_index=%d", m.selected+1),
				fmt.Sprintf("export sel_line=%q", selectedLine),
				m.command,
			},
			";"),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return commandErrorMsg{
				seq: seq,
				err: fmt.Errorf("error running command: %w\n%s", err, string(output)),
			}
		}

		if event == "close" {
			finalOutput := string(output)

			if *selection {
				lines := splitLinesTerminator(finalOutput)

				if m.selected >= 0 && m.selected < len(lines) {
					finalOutput = lines[m.selected]
				}
			}

			return commandCloseMsg{output: finalOutput}
		}

		return commandOutputMsg{
			seq:   seq,
			lines: splitLinesTerminator(string(output)),
		}
	}
}
