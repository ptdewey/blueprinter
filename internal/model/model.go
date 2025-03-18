package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/handler"
	"github.com/ptdewey/blueprinter/internal/ui"

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

				if item.OutputName() == "" {
					dst = filepath.Join(cwd, item.Title())
				} else {
					dst = filepath.Join(cwd, item.OutputName())
				}
			} else {
				dst = os.Args[1]
			}

			if err := handler.CopySelectedItem(item.Path(), dst, item); err != nil {
				fmt.Println("Error copying selected item:", err)
				return m, nil
			}

			// Copy any extra specified template files
			for _, et := range item.Blueprint().Extras {
				if et.TargetTemplate != item.Title() {
					continue
				}

				cwd, err := os.Getwd()
				if err != nil {
					fmt.Println("Error getting current working directory: ", err)
					return m, nil
				}

				for i, t := range et.ExtraTemplates {
					var dst string
					if len(et.ExtraDestinations) > i && et.ExtraDestinations[i] != "" {
						dst = filepath.Join(cwd, et.ExtraDestinations[i])
					} else {
						dst = filepath.Join(cwd, t)
					}

					src := filepath.Join(item.DirPath(), t)
					if err := handler.CopySelectedItem(src, dst, item); err != nil {
						fmt.Println("Error copying additional template files for selected item:", err)
						return m, nil
					}
				}
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := ui.BlueprinterStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return ui.BlueprinterStyle.Render(m.List.View())
}
