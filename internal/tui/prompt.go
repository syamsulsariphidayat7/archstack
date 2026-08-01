package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/syamsulsariphidayat7/archstack/internal/registry"
)

type promptSelectModel struct {
	question    string
	options     []string
	recommended string
	cursor      int
	done        bool
	selected    string
}

func newPromptSelectModel(p registry.Prompt) promptSelectModel {
	cur := 0
	for i, o := range p.Options {
		if o == p.Recommended {
			cur = i
			break
		}
	}
	return promptSelectModel{
		question:    p.Question,
		options:     p.Options,
		recommended: p.Recommended,
		cursor:      cur,
	}
}

func (m promptSelectModel) Update(msg tea.Msg) (promptSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.done = true
			m.selected = m.options[m.cursor]
		}
	}
	return m, nil
}

func (m promptSelectModel) View() string {
	var b strings.Builder
	b.WriteString(PromptStyle.Render(m.question))
	b.WriteString("\n\n")
	for i, opt := range m.options {
		prefix := "  "
		if i == m.cursor {
			prefix = CursorStyle.Render("▶ ")
		}
		label := opt
		if opt == m.recommended {
			label = opt + " " + RecommendedLabel.Render("(disarankan)")
		}
		if i == m.cursor {
			b.WriteString(prefix + SelectedStyle.Render(label) + "\n")
		} else {
			b.WriteString(prefix + NormalStyle.Render(label) + "\n")
		}
	}
	b.WriteString("\n" + DimStyle.Render("↑/↓ navigasi • Enter pilih"))
	return b.String()
}

type promptTextModel struct {
	question    string
	recommended string
	input       textinput.Model
	done        bool
	answer      string
}

func newPromptTextModel(p registry.Prompt) promptTextModel {
	ti := textinput.New()
	ti.Placeholder = p.Recommended
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	return promptTextModel{
		question:    p.Question,
		recommended: p.Recommended,
		input:       ti,
	}
}

func (m promptTextModel) Update(msg tea.Msg) (promptTextModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			val := strings.TrimSpace(m.input.Value())
			if val == "" {
				val = m.recommended
			}
			m.answer = val
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m promptTextModel) View() string {
	var b strings.Builder
	b.WriteString(PromptStyle.Render(m.question))
	b.WriteString("\n")
	b.WriteString("Disarankan: " + RecommendedLabel.Render(m.recommended) + "\n")
	b.WriteString(m.input.View())
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("Enter untuk konfirmasi"))
	return b.String()
}

type promptResult struct {
	answers map[string]string
	err     error
}

func runPrompts(prompts []registry.Prompt, prevAnswers map[string]string) tea.Cmd {
	return func() tea.Msg {
		answers := make(map[string]string)
		for k, v := range prevAnswers {
			answers[k] = v
		}

		p := tea.NewProgram(initialPromptModel(prompts, 0, answers))
		finalModel, err := p.Run()
		if err != nil {
			return promptResult{err: err}
		}
		m := finalModel.(promptFlowModel)
		return promptResult{answers: m.answers, err: nil}
	}
}

type promptFlowModel struct {
	prompts     []registry.Prompt
	current     int
	answers     map[string]string
	selectModel *promptSelectModel
	textModel   *promptTextModel
}

func initialPromptModel(prompts []registry.Prompt, idx int, answers map[string]string) promptFlowModel {
	m := promptFlowModel{
		prompts: prompts,
		current: idx,
		answers: answers,
	}
	if idx < len(prompts) {
		p := prompts[idx]
		switch p.Type {
		case registry.PromptSelect:
			sm := newPromptSelectModel(p)
			m.selectModel = &sm
		case registry.PromptText:
			tm := newPromptTextModel(p)
			m.textModel = &tm
		}
	}
	return m
}

func (m promptFlowModel) Init() tea.Cmd { return nil }

func (m promptFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.current >= len(m.prompts) {
		return m, tea.Quit
	}

	p := m.prompts[m.current]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch p.Type {
	case registry.PromptSelect:
		if m.selectModel == nil {
			break
		}
		sm, cmd := m.selectModel.Update(msg)
		m.selectModel = &sm
		if sm.done {
			m.answers[p.Key] = sm.selected
			m.current++
			if m.current >= len(m.prompts) {
				return m, tea.Quit
			}
			next := m.prompts[m.current]
			switch next.Type {
			case registry.PromptSelect:
				nsm := newPromptSelectModel(next)
				m.selectModel = &nsm
				m.textModel = nil
			case registry.PromptText:
				ntm := newPromptTextModel(next)
				m.textModel = &ntm
				m.selectModel = nil
			}
			return m, nil
		}
		return m, cmd
	case registry.PromptText:
		if m.textModel == nil {
			break
		}
		tm, cmd := m.textModel.Update(msg)
		m.textModel = &tm
		if tm.done {
			m.answers[p.Key] = tm.answer
			m.current++
			if m.current >= len(m.prompts) {
				return m, tea.Quit
			}
			next := m.prompts[m.current]
			switch next.Type {
			case registry.PromptSelect:
				nsm := newPromptSelectModel(next)
				m.selectModel = &nsm
				m.textModel = nil
			case registry.PromptText:
				ntm := newPromptTextModel(next)
				m.textModel = &ntm
				m.selectModel = nil
			}
			return m, nil
		}
		return m, cmd
	}

	return m, nil
}

func (m promptFlowModel) View() string {
	if m.current >= len(m.prompts) {
		return ""
	}
	d := m.prompts[m.current]
	switch d.Type {
	case registry.PromptSelect:
		if m.selectModel != nil {
			return AppStyle.Render(m.selectModel.View())
		}
	case registry.PromptText:
		if m.textModel != nil {
			return AppStyle.Render(m.textModel.View())
		}
	}
	return "Loading..."
}
