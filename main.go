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

func main() {
	cfg := config.ParseConfig()

	items := data.GetItems(cfg.TemplateSources)

	input := flag.String("i", "", "-i <path-to-input-file>")
	dst := flag.String("o", "", "-o <path-to-output-file>")
	verbose := flag.Bool("v", false, "Print output path (verbose)")
	force := flag.Bool("f", false, "Force creation of new file")
	noCopy := flag.Bool("no-copy", false, "Do not create new file, output target path")
	flag.Parse()

	// Determine if Blueprinter should run in CLI or TUI mode
	if len(flag.Args()) < 1 && *input == "" {
		m := ui.Model{
			List:            list.New(items, list.NewDefaultDelegate(), 0, 0),
			TemplateSources: cfg.TemplateSources,
			Flags: ui.Flags{
				Output:    *dst,
				Verbose:   *verbose,
				NoCopy:    *noCopy,
				ForceCopy: *force,
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
		name := input
		if *input == "" {
			*name = flag.Args()[0]
		}

		item, err := pkg.MatchItem(items, *name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !*noCopy {
			if err := pkg.CopyItem(item, *dst, *force, *verbose); err != nil {
				fmt.Println(err)
			}
		} else {
			target, err := pkg.TargetPath(item, *dst)
			if err != nil {
				return
			}
			fmt.Println(target)
		}
	}
}
