package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

type listModel struct {
	items    []registry.Tool
	cursor   int
	vp       viewport.Model
	helpText string
	states   map[string]toolState
}

func newListModel() listModel {
	items := registry.AllTools()
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().PaddingLeft(1)
	return listModel{
		items:    items,
		cursor:   0,
		vp:       vp,
		helpText: "↑/↓ navigasi • Enter pilih • r refresh • q/Esc keluar",
		states:   snapshotStates(items),
	}
}

func (m listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m listModel) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("📦 ArchStack — Development Tools Manager"))
	b.WriteString("\n")
	b.WriteString(HeaderStyle.Render("  NAMA" + padRight("", 7) + "SUMBER" + padRight("", 4) + "STATUS" + padRight("", 10) + "SERVICE" + padRight("", 8) + "DESKRIPSI"))
	b.WriteString("\n")

	for i, tool := range m.items {
		st := m.states[tool.Name]

		prefix := "  "
		if i == m.cursor {
			prefix = CursorStyle.Render("▶ ")
		} else {
			prefix = "  "
		}

		name := tool.Name
		if i == m.cursor {
			name = SelectedStyle.Render(tool.Name)
		} else {
			name = NormalStyle.Render(tool.Name)
		}

		source := DimStyle.Render(string(tool.From))

		var status string
		if st.installed {
			status = InstalledStyle.Render("installed")
		} else {
			status = NotInstalledStyle.Render("not installed")
		}

		var svc string
		if tool.Service != "" {
			switch st.svc {
			case system.ServiceRunning:
				svc = ServiceRunningStyle.Render("running")
			case system.ServiceStopped:
				svc = ServiceStoppedStyle.Render("stopped")
			default:
				svc = ServiceNoneStyle.Render("-")
			}
		} else {
			svc = ServiceNoneStyle.Render("-")
		}

		desc := DimStyle.Render(tool.Desc)

		b.WriteString(prefix + name + padRight("", 12-len(tool.Name)) +
			source + padRight("", 10) +
			status + padRight("", 14) +
			svc + padRight("", 10) +
			desc + "\n")
	}

	b.WriteString(HelpStyle.Render(m.helpText))
	return b.String()
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return ""
	}
	return strings.Repeat(" ", n-len(s))
}
