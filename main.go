package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

var (
	resultErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1"))
)

var (
	command     string
	finalOutput string
)

var (
	initial     = pflag.StringP("initial", "i", "", "Initial prompt text")
	placeholder = pflag.StringP("placeholder", "p", "", "Placeholder text")
	selection   = pflag.BoolP("selection", "s", false, "Only return the selected item")
)

type model struct {
	w, h int

	input textinput.Model

	selected int
	lines    []string
	// list  list.Model

	err error
}

var usageStyle = lipgloss.NewStyle().Width(80).PaddingLeft(2)

var usageMain = usageStyle.Render(strings.TrimSpace(`
%[1]s [OPTIONS] COMMAND
%[1]s search QUERY
  
This is a simple interactive command line tool that allows you to run a command at every key press. It can be used to create interactive prompts or menus.

The search subcommand performs a fuzzy search on the stdin lines and prints the results highlighting the matched characters. If no query is provided, it will print all the lines.
`)) + "\n\n"

var usageCommand = usageStyle.Render(strings.TrimSpace(`
This command will be executed at every key press so mind the performance. It has access to the following environment variables:

  $prompt     the current prompt text
  $event      the event that triggered the command (open, key, select, close)
  $sel_line   the selected line
  $sel_index  the index of the selected item, starting from 1
`)) + "\n"

func main() {
	pflag.CommandLine.Init(os.Args[0], pflag.ContinueOnError)
	pflag.Usage = func() {
		w := os.Stderr

		fmt.Fprintln(w, "Usage:")
		fmt.Fprintf(w, usageMain, os.Args[0])

		fmt.Fprintln(w, "Options:")
		pflag.PrintDefaults()
		fmt.Fprintln(w)

		fmt.Fprintln(w, "Command:")
		fmt.Fprintln(w, usageCommand)
	}

	err := pflag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}

		log.Fatalf("error parsing flags: %v", err)
	}

	if pflag.Arg(0) == "search" {
		if err := search(pflag.Arg(1)); err != nil {
			log.Fatalf("error searching: %v", err)
			return
		}

		return
	}

	command = pflag.Arg(0)
	if command == "" {
		panic("No command provided")
	}

	m := model{
		input: textinput.New(),
	}

	m.input.Prompt = "> "
	m.input.SetValue(*initial)
	m.input.Placeholder = *placeholder
	m.input.Focus()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if finalOutput != "" {
		if finalOutput[len(finalOutput)-1] == '\n' {
			finalOutput = finalOutput[:len(finalOutput)-1]
		}

		fmt.Println(finalOutput)
	}
}

func (m model) Init() tea.Cmd {
	return m.runCommand("open")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyUp:
			m.selected = max(0, m.selected-1)

			return m, m.runCommand("select")

		case tea.KeyDown:
			m.selected = min(len(m.lines)-1, m.selected+1)

			return m, m.runCommand("select")

		case tea.KeyEnter:
			return m, m.runCommand("close")

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		default:
			m.input, _ = m.input.Update(msg)

			return m, m.runCommand("key")
		}

	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		// m.input.Width = m.w
		return m, nil

	case commandOutputMsg:
		m.lines = msg.lines
		m.err = nil

		m.selected = min(m.selected, len(m.lines)-1)

		return m, nil

	case commandErrorMsg:
		m.lines = nil
		m.err = msg.err
		return m, nil

	case commandCloseMsg:
		return m, tea.Quit

	}

	m.input, cmd = m.input.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	v := strings.Builder{}

	v.WriteString(
		lipgloss.NewStyle().
			Width(m.w-2).
			Height(1).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			Render(
				m.input.View(),
			),
	)

	v.WriteString("\n")

	content := ""

	if m.err != nil {
		content = resultErrorStyle.Render(m.err.Error())
	} else if len(m.lines) == 0 {
		content = "No output"
	} else {
		items := []string{}

		for i, line := range m.lines {
			formattedItem := lipgloss.NewStyle().Width(m.w - 10).Render(line)

			if i == m.selected {
				formattedItem = lipgloss.JoinHorizontal(lipgloss.Top, "▶ ", formattedItem)
			} else {
				formattedItem = lipgloss.JoinHorizontal(lipgloss.Top, "  ", formattedItem)
			}

			items = append(items, formattedItem)
		}

		content = lipgloss.JoinVertical(lipgloss.Left, items...)

		// only show lines that fit in the screen
		maxHeight := m.h - 5
		if strings.Count(content, "\n") > maxHeight {
			contentLines := strings.Split(content, "\n")
			content = strings.Join(contentLines[:maxHeight], "\n")
		}

		content = strings.TrimRight(content, "\n ")
	}

	v.WriteString(
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(
				lipgloss.NewStyle().
					Width(m.w-4).
					Height(m.h-5).
					Padding(0, 1).
					Border(lipgloss.RoundedBorder()).
					Render(content),
			),
	)

	return v.String()
}

type commandOutputMsg struct{ lines []string }

type commandErrorMsg struct{ err error }

type commandCloseMsg struct{}

func (m model) runCommand(event string) tea.Cmd {
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
				command,
			},
			";"),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return commandErrorMsg{
				err: fmt.Errorf("error running command: %w\n%s", err, string(output)),
			}
		}

		if event == "close" {
			finalOutput = string(output)

			if *selection {
				lines := splitLinesTerminator(finalOutput)

				if m.selected >= 0 && m.selected < len(lines) {
					finalOutput = lines[m.selected]
				}
			}

			return commandCloseMsg{}
		}

		return commandOutputMsg{
			lines: splitLinesTerminator(string(output)),
		}
	}
}

func splitLinesTerminator(s string) []string {
	// s = strings.ReplaceAll(s, "\n", "\n~")

	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// for i, line := range lines {
	// 	lines[i] = strings.ReplaceAll(line, " ", "·")
	// }

	return lines
}
