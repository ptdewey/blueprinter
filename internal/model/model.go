package model

import (
	"blueprinter/internal/data"
	"blueprinter/internal/handler"
	"blueprinter/internal/ui"
	"fmt"
	"os"
	"path/filepath"

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
			if selectedItem == nil {
				return m, nil
			}

			item, ok := selectedItem.(data.Item)
			if !ok {
				fmt.Println("Error: could not type assert to data.Item")
				return m, nil
			}

			var dst string
			if len(os.Args) < 2 {
				cwd, err := os.Getwd()
				if err != nil {
					fmt.Println("Error getting current working directory: ", err)
					return m, nil
				}

				dst = filepath.Join(cwd, item.Title())
			} else {
				dst = filepath.Join(os.Args[1], item.Title())
			}

			if err := handler.CopySelectedItem(item.Path(), dst); err != nil {
				fmt.Println("Error copying selected item:", err)
				return m, nil
			}

			return m, tea.Quit
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
