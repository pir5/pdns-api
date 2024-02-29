package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pir5/pdns-api/cmd"
)

// @title PDNS-API
// @version 1.0
// @description This is PDNS RESTful API Server.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey ID
// @in header
// @name PIR5-ID
// @securityDefinitions.apikey Secret
// @in header
// @name PIR5-SECRET
// @tag.name domains
// @tag.name records
// @tag.name viron

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'pdns-api help'.
var commands = []*cmd.Command{
	cmd.CmdServer,
}

func main() {
	cmdFlags := cmd.GlobalFlags{}
	cmdFlags.ConfPath = flag.String("config", "/etc/pdns-api/api.toml", "config file path")
	flag.Usage = usage
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hostname, _ := os.Hostname()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(slog.String("server_host", hostname))
	slog.SetDefault(logger)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()

			if err := cmd.Run(&cmdFlags, args); err != nil {

				fmt.Println(err)
				os.Exit(2)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "pdns-api: unknown subcommand %q\nRun ' pdns-api help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `pdns-api is a tool for 

Usage:

	pdns-api command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "pdns-api help [command]" for more information about a command.

`

var helpTemplate = `usage: pdns-api {{.UsageLine}}

{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'pdns-api help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: pdns-api help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'pdns-api help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'pdns-api help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'pdns-api help'.\n", arg)
	os.Exit(2) // failed at 'pdns-api help cmd'
}
