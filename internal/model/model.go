package model

import (
	"blueprinter/internal/handler"
	"blueprinter/internal/ui"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	List            list.Model
	TemplateSources []string
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
			if selectedItem != nil {
				// TODO: get full path to selected item
				// - might need to do this via a config file or cli flags (both?)
				// - the selectedItem type is not modifiable, and it only provides means to get the filter key

				var dst string
				if len(os.Args) < 2 {
					dst = "./" + selectedItem.FilterValue()
				} else {
					dst = os.Args[1] + selectedItem.FilterValue()
				}

				// CHANGE: replace this with non-hardcoded version
				// FIX: figure out how to iterate through all template sources
				src := fmt.Sprintf("%s/%s", m.TemplateSources[0], selectedItem.FilterValue())
				_, err := handler.CopySelectedItem(src, dst)
				if err != nil {
					fmt.Println("Error copying selected item:", err)
					return m, nil
				}

				return m, tea.Quit
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := ui.DocStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return ui.DocStyle.Render(m.List.View())
}
