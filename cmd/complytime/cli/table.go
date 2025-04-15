package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

// This example was adapted from charmbracelet/bubbles
// The original example can be found here:
// https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's update the %s component: %s.", m.table.SelectedRow()[1], m.table.SelectedRow()[2]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// parameterTable instantiates a separate table
// parameterTable instantiates a separate table.
// Choose to update parameters via command flags.
func parameterTable() {
	columns := []table.Column{
		// Parameter Headers
		{Title: "Framework ID", Width: 20},
		{Title: "Component Type", Width: 20},
		{Title: "Parameters", Width: 20},
	}
	rows := []table.Row{
		{"anssi_bp28_minimal", "Parameter", "param_1"},
		{"anssi_bp28_minimal", "Parameter", "param_2"},
		{"anssi_bp28_minimal", "Parameter", "param_3"},
		{"anssi_bp28_minimal", "Parameter", "param_4"},
		{"anssi_bp28_minimal", "Parameter", "param_5"},
		{"anssi_bp28_minimal", "Parameter", "param_6"},
		{"anssi_bp28_minimal", "Parameter", "param_7"},
		{"anssi_bp28_minimal", "Parameter", "param_8"},
		{"anssi_bp28_minimal", "Parameter", "param_9"},
		{"anssi_bp28_minimal", "Parameter", "param_1"},
		{"anssi_bp28_minimal", "Parameter", "param_2"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// controlsTable instantiates a separate table.
// Choose to update controls via command flags.
func controlsTable() {
	columns := []table.Column{
		// Control Headers
		{Title: "Framework ID", Width: 20},
		{Title: "Component Type", Width: 20},
		{Title: "Controls", Width: 20},
	}
	rows := []table.Row{
		{"anssi_bp28_minimal", "Control", "control_1"},
		{"anssi_bp28_minimal", "Control", "control_2"},
		{"anssi_bp28_minimal", "Control", "control_3"},
		{"anssi_bp28_minimal", "Control", "control_4"},
		{"anssi_bp28_minimal", "Control", "control_5"},
		{"anssi_bp28_minimal", "Control", "control_6"},
		{"anssi_bp28_minimal", "Control", "control_7"},
		{"anssi_bp28_minimal", "Control", "control_8"},
		{"anssi_bp28_minimal", "Control", "control_9"},
		{"anssi_bp28_minimal", "Control", "control_n"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

// rulesTable instantiates a separate table.
// Choose to update rules via command flags.
func rulesTable() {
	columns := []table.Column{
		// Rule Headers
		{Title: "Framework ID", Width: 20},
		{Title: "Component Type", Width: 20},
		{Title: "Rules", Width: 20},
	}
	rows := []table.Row{
		{"anssi_bp28_minimal", "Rule", "rule_1"},
		{"anssi_bp28_minimal", "Rule", "rule_2"},
		{"anssi_bp28_minimal", "Rule", "rule_3"},
		{"anssi_bp28_minimal", "Rule", "rule_4"},
		{"anssi_bp28_minimal", "Rule", "rule_5"},
		{"anssi_bp28_minimal", "Rule", "rule_6"},
		{"anssi_bp28_minimal", "Rule", "rule_7"},
		{"anssi_bp28_minimal", "Rule", "rule_8"},
		{"anssi_bp28_minimal", "Rule", "rule_9"},
		{"anssi_bp28_minimal", "Rule", "rule_n"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"

}

// ComponentTable calls the functions for each option.
// Available components of the assessment plan will print in table.
// Updates to assessment-plan components will be present in table
func ComponentTable() {
	controlsTable()
	rulesTable()
	parameterTable()
}
