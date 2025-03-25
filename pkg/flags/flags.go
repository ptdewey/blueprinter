package flags

import "flag"

type Flags struct {
	Input    string
	Output   string
	Verbose  bool
	Force    bool
	NoCopy   bool
	DumpList bool
	ID       string
	Tags     string
}

var BlueprinterFlags Flags

func Parse() *Flags {
	flag.StringVar(&BlueprinterFlags.Input, "i", "", "-i <path-to-input-file>")
	flag.StringVar(&BlueprinterFlags.Output, "o", "", "-o <path-to-output-file>")
	flag.BoolVar(&BlueprinterFlags.Verbose, "v", false, "Print output path (verbose)")
	flag.BoolVar(&BlueprinterFlags.Force, "f", false, "Force creation of new file")
	flag.BoolVar(&BlueprinterFlags.NoCopy, "no-copy", false, "Do not create new file, output target path")
	flag.BoolVar(&BlueprinterFlags.DumpList, "dump-list", false, "Dump list contents to stdout")
	flag.StringVar(&BlueprinterFlags.ID, "id", "", "Set {{ .id }} template variable")
	flag.StringVar(&BlueprinterFlags.Tags, "t", "", "Set {{ .tags }} template variable")

	flag.Parse()
	return &BlueprinterFlags
}
