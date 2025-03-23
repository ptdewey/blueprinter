package ui

import (
	"fmt"

	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/pkg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Flags struct {
	Output  string
	Verbose bool
}

type Model struct {
	List            list.Model
	TemplateSources []string
	Flags           Flags
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.List.FilterState() == list.Filtering {
				m.List.ResetFilter()
				return m, nil
			}
		case "enter":
			selectedItem := m.List.SelectedItem()
			if selectedItem == nil {
				return m, nil
			}

			item, ok := selectedItem.(data.Item)
			if !ok {
				fmt.Println("Error: could not type assert to data.Item")
				return m, nil
			}

			if err := pkg.CopyItem(&item, m.Flags.Output, m.Flags.Verbose); err != nil {
				return m, nil
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := blueprinterStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return blueprinterStyle.Render(m.List.View())
}
