package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ptdewey/blueprinter/internal/config"
	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/ui"
	"github.com/ptdewey/blueprinter/pkg"
	"github.com/ptdewey/blueprinter/pkg/flags"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.ParseConfig()

	items := data.GetItems(cfg.TemplateSources)

	f := flags.Parse()

	if f.DumpList {
		for _, item := range items {
			fmt.Println(item.(data.Item).Path())
		}
		return
	}

	// Determine if Blueprinter should run in CLI or TUI mode
	if len(flag.Args()) < 1 && f.Input == "" {
		m := ui.Model{
			List:            list.New(items, list.NewDefaultDelegate(), 0, 0),
			TemplateSources: cfg.TemplateSources,
		}
		m.List.Title = "Available Templates"

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	} else {
		// Run in CLI mode
		name := f.Input
		if name == "" {
			name = flag.Args()[0]
		}

		item, err := pkg.MatchItem(items, name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !f.NoCopy {
			if err := pkg.CopyItem(item, f.Output, f.Force, f.Verbose); err != nil {
				fmt.Println(err)
			}
		} else {
			target, err := pkg.TargetPath(item, f.Output)
			if err != nil {
				return
			}
			fmt.Println(target)
		}
	}
}
