# Blueprinter

Blueprinter is a command-line tool built with Bubble Tea that allows users to browse template files and directories from configured source locations and copy them to a target destination. It provides an interactive list of available templates from which users can select, and then copy the selected item to their current working directory or a specified target directory.

![Blueprinter Interface](./assets/screenshot-1.png)

Blueprinter allows fuzzy search and selection from files in configured template directories.

![Blueprinter Filtered Interface](./assets/screenshot-2.png)


## Installation

Blueprinter can be installed using `go install`

```bash
go install github.com/ptdewey/blueprinter@latest
```

It can also be built from source if desired
```bash
git clone https://github.com/ptdewey/blueprinter.git
cd blueprinter
go build
```

## Usage

Blueprinter is run by calling the executable from the command line.

```bash
# installed with go get (or from source with executable in PATH)
blueprinter

# built from source (not in PATH)
./blueprinter
```

To specify an output location, add an argument for the desired output location.

```bash
# add selected template to the 'example' directory
./blueprinter ./example
```

## Configuration

Blueprinter reads configuration from JSON files, which specify the directories where your template files are located. If no configuration file is found, it falls back to default directories like `~/Templates` or `~/Documents/Templates`. Here's how the configuration system works:

### Configuration Files

By default, Blueprinter looks for one of the following configuration files in either the current Git repository root (if available) or in the home directory:

- `blueprinter.json`
- `.blueprinter.json`
- `.blueprinterrc`
- `.blueprinterrc.json`
- `blueprinterrc.json`

These files should contain a JSON object with a `template-sources` key, which lists the directories where your templates are stored.

### Example Configuration (`blueprinter.json`)

```json
{
  "template-sources": [
    "~/Templates",
    "~/Documents/MyCustomTemplates"
  ]
}
```

In this example, Blueprinter will look for templates in `~/Templates` and `~/Documents/MyCustomTemplates`.

### Default Configuration

If no configuration file is found, Blueprinter will try to use the following default directories within the user's home directory:

- `~/Templates`
- `~/Documents/Templates`
- `~/Documents/templates`

If one of these directories exists, it will be used as the default template source. If none of the default directories exist and no configuration file is present, Blueprinter will panic and exit with an error.

