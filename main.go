package main

import (
	"fmt"
	"os"

	"blueprinter/internal/data"
	"blueprinter/internal/model"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	items := data.GetItems()

	// TODO: change delegate to not use default styles
	m := model.Model{List: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Available Templates"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
