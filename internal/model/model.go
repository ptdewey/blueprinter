package model

import (
	"blueprinter/internal/ui"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	List list.Model
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
			if m.List.FilterState() == list.Filtering {
				m.List.ResetFilter()
			}

			selectedItem := m.List.SelectedItem()
			if selectedItem != nil {
				// TODO: get full path to selected item
				// - might need to do this via a config file or cli flags (both?)
				// - the selectedItem type is not modifiable, and it only provides means to get the filter key

				fmt.Println(selectedItem.FilterValue())

				var dst string
				if len(os.Args) < 2 {
					dst = "./"
				} else {
					dst = os.Args[1]
				}
				fmt.Println(dst) // TODO: remove this

				// handler.CopySelected(src, dst)

				// TODO: figure out how to write copied file path to stdout

				// return m, tea.Quit
			}
			// TODO: maybe do something besides quit in case of no selected item?

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
