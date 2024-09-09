package main

import (
	"fmt"
	"os"

	"blueprinter/internal/config"
	"blueprinter/internal/data"
	"blueprinter/internal/model"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.ParseConfig()

	items := data.GetItems(cfg.TemplateSources)

	m := model.Model{
		List:            list.New(items, list.NewDefaultDelegate(), 0, 0),
		TemplateSources: cfg.TemplateSources,
	}
	m.List.Title = "Available Templates"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
