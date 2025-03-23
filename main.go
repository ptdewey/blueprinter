package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ptdewey/blueprinter/internal/config"
	"github.com/ptdewey/blueprinter/internal/data"
	"github.com/ptdewey/blueprinter/internal/ui"
	"github.com/ptdewey/blueprinter/pkg"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	input   string
	dst     string
	verbose bool
)

func main() {
	cfg := config.ParseConfig()

	items := data.GetItems(cfg.TemplateSources)

	flag.StringVar(&input, "i", "", "-i <path-to-input-file>")
	flag.StringVar(&dst, "o", "", "-o <path-to-output-file>")
	flag.BoolVar(&verbose, "v", false, "-v")
	flag.Parse()

	// Determine if Blueprinter should run in CLI or TUI mode
	if len(flag.Args()) < 1 && input == "" {
		m := ui.Model{
			List:            list.New(items, list.NewDefaultDelegate(), 0, 0),
			TemplateSources: cfg.TemplateSources,
			Flags: ui.Flags{
				Output:  dst,
				Verbose: verbose,
			},
		}
		m.List.Title = "Available Templates"

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	} else {
		// Run in CLI mode
		if input != "" {
			item, err := pkg.MatchItem(items, input)
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := pkg.CopyItem(item, dst, verbose); err != nil {
				fmt.Println(err)
			}
		} else {
			item, err := pkg.MatchItem(items, flag.Args()[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := pkg.CopyItem(item, dst, verbose); err != nil {
				fmt.Println(err)
			}
		}
	}
}
