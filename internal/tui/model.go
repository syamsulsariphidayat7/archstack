package tui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

type state int

const (
	stateList state = iota
	stateSubmenu
	statePrompts
	stateConfirm
	stateDone
)

type actionDoneMsg struct {
	err error
}

type Model struct {
	state   state
	list    listModel
	tool    *registry.Tool
	answers map[string]string

	submenuCursor int
	submenuItems  []string

	confirmMsg string
	execErr    error
	execMsg    string

	promptModel *promptFlowModel
}

func NewModel() Model {
	return Model{
		state:   stateList,
		list:    newListModel(),
		answers: make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case actionDoneMsg:
		m.execErr = msg.err
		m.state = stateDone
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	default:
		return m.handleOther(msg)
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.state {
	case stateList:
		return m.updateList(msg)
	case stateSubmenu:
		return m.updateSubmenu(msg)
	case stateConfirm:
		return m.updateConfirm(msg)
	case stateDone:
		if msg.String() == "enter" || msg.String() == "esc" || msg.String() == "q" {
			m.state = stateList
			m.execErr = nil
			m.tool = nil
			m.refreshStates()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleOther(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == statePrompts {
		return m.updatePrompts(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "r":
		m.list.states = snapshotStates(m.list.items)
		return m, nil
	case "enter":
		tool := m.list.items[m.list.cursor]
		m.tool = &tool
		installed := system.IsInstalled(tool.Binary, tool.Pkg)
		if installed {
			m.submenuItems = buildSubmenu(&tool)
			m.submenuCursor = 0
			m.state = stateSubmenu
			return m, nil
		}
		if len(tool.Prompts) > 0 {
			m.answers = registry.GetCachedAnswers(tool.Name)
			m.state = statePrompts
			m.promptModel = nil
			pm := initialPromptModel(tool.Prompts, 0, m.answers)
			m.promptModel = &pm
			return m, nil
		}
		return m, execInstall(m.tool, nil)
	case "up", "k":
		if m.list.cursor > 0 {
			m.list.cursor--
		}
	case "down", "j":
		if m.list.cursor < len(m.list.items)-1 {
			m.list.cursor++
		}
	}
	return m, nil
}

func (m *Model) refreshStates() {
	m.list.states = snapshotStates(m.list.items)
}

func buildSubmenu(tool *registry.Tool) []string {
	items := []string{"Uninstall " + tool.Name}
	if tool.Service != "" {
		state, _ := system.ServiceStatus(tool.Service)
		switch state {
		case system.ServiceRunning:
			items = append(items, "Stop service "+tool.Service)
		default:
			items = append(items, "Start service "+tool.Service)
		}
	}
	items = append(items, "Batal")
	return items
}

func (m Model) updateSubmenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.submenuCursor > 0 {
			m.submenuCursor--
		}
	case "down", "j":
		if m.submenuCursor < len(m.submenuItems)-1 {
			m.submenuCursor++
		}
	case "enter":
		sel := m.submenuItems[m.submenuCursor]
		if sel == "Batal" {
			m.state = stateList
			return m, nil
		}
		if strings.HasPrefix(sel, "Uninstall") {
			m.state = stateConfirm
			m.confirmMsg = fmt.Sprintf("Hapus %s? (y/n)", m.tool.Name)
			return m, nil
		}
		if strings.Contains(sel, "Start") {
			return m, execRun(m.tool)
		}
		if strings.Contains(sel, "Stop") {
			return m, execStop(m.tool)
		}
	case "q", "esc":
		m.state = stateList
		return m, nil
	}
	return m, nil
}

func (m Model) updatePrompts(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.promptModel == nil {
		return m, nil
	}

	pm, cmd := m.promptModel.Update(msg)
	var ok bool
	m.promptModel, ok = pm.(*promptFlowModel)
	if !ok {
		return m, nil
	}

	if m.promptModel.current >= len(m.promptModel.prompts) {
		answers := m.promptModel.answers
		for k, v := range answers {
			m.answers[k] = v
		}
		registry.SetCachedAnswers(m.tool.Name, answers)
		return m, execInstall(m.tool, answers)
	}

	return m, cmd
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, execUninstall(m.tool)
	case "n", "N", "esc", "q":
		m.state = stateList
		return m, nil
	}
	return m, nil
}

func execInstall(tool *registry.Tool, answers map[string]string) tea.Cmd {
	var cmd *exec.Cmd
	switch tool.From {
	case registry.SourcePacman:
		cmd = exec.Command("sudo", "pacman", "-S", "--needed", "--noconfirm", tool.Pkg)
	case registry.SourceYay:
		cmd = exec.Command("yay", "-S", "--needed", "--noconfirm", tool.Pkg)
	}
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return actionDoneMsg{err: err}
	})
}

func execUninstall(tool *registry.Tool) tea.Cmd {
	state, _ := system.ServiceStatus(tool.Service)
	if tool.Service != "" && state == system.ServiceRunning {
		stopCmd := exec.Command("sudo", "systemctl", "stop", tool.Service)
		_ = stopCmd.Run()
	}
	cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", tool.Pkg)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return actionDoneMsg{err: err}
	})
}

func execRun(tool *registry.Tool) tea.Cmd {
	cmd := exec.Command("sudo", "systemctl", "enable", "--now", tool.Service)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return actionDoneMsg{err: err}
	})
}

func execStop(tool *registry.Tool) tea.Cmd {
	cmd := exec.Command("sudo", "systemctl", "stop", tool.Service)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return actionDoneMsg{err: err}
	})
}

func (m Model) View() string {
	switch m.state {
	case stateList:
		return AppStyle.Render(m.list.View())
	case stateSubmenu:
		return m.renderSubmenu()
	case statePrompts:
		if m.promptModel != nil {
			return m.promptModel.View()
		}
		return AppStyle.Render("Memuat pertanyaan...")
	case stateConfirm:
		return AppStyle.Render(m.renderConfirm())
	case stateDone:
		return AppStyle.Render(m.renderDone())
	}
	return ""
}

func (m Model) renderSubmenu() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(fmt.Sprintf("Aksi untuk %s", m.tool.Name)))
	b.WriteString("\n\n")
	for i, item := range m.submenuItems {
		prefix := "  "
		if i == m.submenuCursor {
			prefix = CursorStyle.Render("▶ ")
		}
		if i == m.submenuCursor {
			b.WriteString(prefix + SelectedStyle.Render(item) + "\n")
		} else {
			b.WriteString(prefix + NormalStyle.Render(item) + "\n")
		}
	}
	b.WriteString("\n" + DimStyle.Render("↑/↓ navigasi • Enter pilih • Esc kembali"))
	return b.String()
}

func (m Model) renderConfirm() string {
	var b strings.Builder
	b.WriteString(ErrorStyle.Render(m.confirmMsg))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("y/Y untuk ya • n/N/Enter untuk tidak"))
	return b.String()
}

func (m Model) renderDone() string {
	var b strings.Builder
	if m.execErr != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("[error] %s", m.execErr)))
	} else {
		b.WriteString(SuccessStyle.Render("[done] Berhasil"))
	}
	b.WriteString("\n\n" + DimStyle.Render("Tekan Enter/Esc untuk kembali"))
	return b.String()
}

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
