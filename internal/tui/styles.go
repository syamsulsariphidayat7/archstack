package tui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		PaddingBottom(1)

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("247")).
		PaddingBottom(1)

	CursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("63"))

	SelectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("15"))

	NormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	InstalledStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	NotInstalledStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203"))

	ServiceRunningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	ServiceStoppedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203"))

	ServiceNoneStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	DimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingTop(1)

	PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	RecommendedLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Italic(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	Separator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		SetString(" │ ").String()
)
