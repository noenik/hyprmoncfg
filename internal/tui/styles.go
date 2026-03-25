package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	title       lipgloss.Style
	header      lipgloss.Style
	focused     lipgloss.Style
	activePane  lipgloss.Style
	inactive    lipgloss.Style
	statusOK    lipgloss.Style
	statusError lipgloss.Style
	help        lipgloss.Style
	warning     lipgloss.Style
}

func newStyles() styles {
	return styles{
		title:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("25")).Padding(0, 1),
		header:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("111")),
		focused:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")),
		activePane:  lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69")).Padding(0, 1),
		inactive:    lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1),
		statusOK:    lipgloss.NewStyle().Foreground(lipgloss.Color("84")).Bold(true),
		statusError: lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true),
		help:        lipgloss.NewStyle().Foreground(lipgloss.Color("247")),
		warning:     lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
	}
}
