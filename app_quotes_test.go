package flagx_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mmbros/flagx"
)

const (
	defaultConfigType = "yaml"
	defaultMode       = "1"
	defaultWorkers    = 1
)

type appArgs struct {
	config     string
	configType string
	database   string
	dryrun     bool
	isins      []string
	proxy      string
	sources    []string
	workers    int
	mode       string
}

const (
	usageApp = `Usage:
    %s <command> [options]

Available commands:
    get      Get the quotes of the specified isins
    sources  Show available sources
    tor      Checks if Tor network will be used
`
	usageGet = `Usage:
    %s [options]

Options:
    -c, --config      path     config file
        --config-type string   used if config file does not have the extension in the name;
                               accepted values are: YAML, TOML and JSON 
    -i, --isins       strings  list of isins to get the quotes
    -n, --dry-run              perform a trial run with no request/updates made
    -p, --proxy       url      default proxy
    -s, --sources     strings  list of sources to get the quotes from
    -w, --workers     int      number of workers (default 1)
    -d, --database    dns      sqlite3 database used to save the quotes
    -m, --mode        char     result mode: "1" first success or last error (default)
                                            "U" all errors until first success 
                                            "A" all 
`

	usageTor = `Usage:
	%s [options]

Checks if Tor network will be used to get the quote.

To use the Tor network the proxy must be defined through:
	1. proxy argument parameter
	2. proxy config file parameter
	3. HTTP_PROXY, HTTPS_PROXY and NOPROXY environment variables.

Options:
    -c, --config      path    config file (default is $HOME/.quotes.yaml)
	    --config-type string  used if config file does not have the extension in the name;
	                          accepted values are: YAML, TOML and JSON 
    -p, --proxy       url     proxy to test the Tor network
`

	usageSources = `Usage:
	%s

Prints list of available sources.
`
)

// argsQuotes contains the argument passed to the app.
// It is global in order to test the values
var argsQuotes *appArgs

func initApp() *flagx.Command {

	app := &flagx.Command{
		ParseExec: runQuotesApp,
		SubCmd: map[string]*flagx.Command{
			"get,g": {
				ParseExec: runQuotesGet,
			},
			"tor,t": {
				ParseExec: runQuotesTor,
			},
			"sources,s": {
				ParseExec: runQuotesSources,
			},
		},
	}

	return app
}

func runQuotesApp(name string, arguments []string) error {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageApp, name)
	}

	err := fs.Parse(arguments)
	return err
}

func runQuotesGet(name string, arguments []string) error {
	// it is used a module level declaration for test porpouses.
	// normally do: argsQuotes := &appArgs{}
	argsQuotes = &appArgs{}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageGet, name)
	}
	flagx.AliasedStringVar(fs, &argsQuotes.config, "config,c", "", "config file")
	flagx.AliasedStringVar(fs, &argsQuotes.configType, "config-type", defaultConfigType, "used if config file does not have the extension in the name; accepted values are: YAML, TOML and JSON")
	flagx.AliasedBoolVar(fs, &argsQuotes.dryrun, "dry-run,n", false, "perform a trial run with no request/updates made")
	flagx.AliasedStringVar(fs, &argsQuotes.proxy, "proxy,p", "", "default proxy")
	flagx.AliasedIntVar(fs, &argsQuotes.workers, "workers,w", defaultWorkers, "number of workers")
	flagx.AliasedStringVar(fs, &argsQuotes.database, "database,d", "", "sqlite3 database used to save the quotes")
	flagx.AliasedStringVar(fs, &argsQuotes.mode, "mode,m", defaultMode, `result mode: "1" first success or last error (default), "U" all errors until first success, "A" all`)
	flagx.AliasedStringsVar(fs, &argsQuotes.isins, "isins,i", "list of isins to get the quotes")
	flagx.AliasedStringsVar(fs, &argsQuotes.sources, "sources,s", "list of sources to get the quotes from")

	err := fs.Parse(arguments)
	return err
}

func runQuotesTor(name string, arguments []string) error {
	// it is used a module level declaration for test porpouses.
	// normally do: argsQuotes := &appArgs{}
	argsQuotes = &appArgs{}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageTor, name)
	}
	flagx.AliasedStringVar(fs, &argsQuotes.config, "config,c", "", "config file")
	flagx.AliasedStringVar(fs, &argsQuotes.configType, "config-type", defaultConfigType, "used if config file does not have the extension in the name; accepted values are: YAML, TOML and JSON")
	flagx.AliasedStringVar(fs, &argsQuotes.proxy, "proxy,p", "", "default proxy")

	err := fs.Parse(arguments)
	return err
}

func runQuotesSources(name string, arguments []string) error {
	// it is used a module level declaration for test porpouses.
	// normally do: argsQuotes := &appArgs{}
	argsQuotes = &appArgs{}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageSources, name)
	}

	err := fs.Parse(arguments)
	return err
}

func Test_QuotesApp(t *testing.T) {
	const AppName = "QUOTES"

	tests := []struct {
		name        string
		args        string
		wantErr     error
		wantErrMsg  string
		wantOutput  string
		wantOptions *appArgs
	}{
		{
			name: "empty command line",
			args: "",
		},
		{
			name:       "app help",
			args:       "-h",
			wantOutput: fmt.Sprintf(usageApp, AppName),
			wantErr:    flag.ErrHelp,
		},
		{
			name:       "app unknown option",
			args:       "-x",
			wantErrMsg: "flag provided but not defined: -x",
		},
		{
			name:       "app help + unknown option",
			args:       "-h -x",
			wantOutput: fmt.Sprintf(usageApp, AppName),
			wantErr:    flag.ErrHelp,
		},
		{
			name:       "app unknown option + help",
			args:       "-x -h",
			wantErrMsg: "flag provided but not defined: -x",
		},
		{
			name:    "app unknown command",
			args:    "cmd",
			wantErr: flagx.ErrCommandNotFound,
		},
		{
			name:       "app get help",
			args:       "get -h",
			wantOutput: fmt.Sprintf(usageGet, fmt.Sprintf("%s get", AppName)),
			wantErr:    flag.ErrHelp,
		},
		{
			name:       "app g help",
			args:       "g -h",
			wantOutput: fmt.Sprintf(usageGet, fmt.Sprintf("%s get", AppName)),
			wantErr:    flag.ErrHelp,
		},
		{
			name: "app get with -c flag",
			args: "get -c /path/to/file",
			wantOptions: &appArgs{
				workers:    defaultWorkers,
				mode:       defaultMode,
				configType: defaultConfigType,
				config:     "/path/to/file",
			},
		},
		{
			name: "app get with --config flag",
			args: "get --config /path/to/file",
			wantOptions: &appArgs{
				workers:    defaultWorkers,
				mode:       defaultMode,
				configType: defaultConfigType,
				config:     "/path/to/file",
			},
		},
		{
			name: "app get with --config-type flag",
			args: "get --config-type json",
			wantOptions: &appArgs{
				workers:    defaultWorkers,
				mode:       defaultMode,
				configType: "json",
			},
		},
		{
			name: "get dry-run, isins and workers",
			args: "get -i isin1,isin2 --isins isin3 --w=5 --dry-run=1",
			wantOptions: &appArgs{
				workers:    5,
				mode:       defaultMode,
				configType: defaultConfigType,
				isins:      []string{"isin1", "isin2", "isin3"},
				dryrun:     true,
			},
		},
		{
			name:       "tor help",
			args:       "t -h",
			wantOutput: fmt.Sprintf(usageTor, fmt.Sprintf("%s tor", AppName)),
			wantErr:    flag.ErrHelp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var out strings.Builder
			flag.CommandLine.SetOutput(&out)

			app := initApp()

			os.Args = strings.Split(AppName+" "+tt.args, " ")
			err := flagx.Run(app)

			if err != nil {
				if (tt.wantErr == nil) && (tt.wantErrMsg == "") {
					t.Errorf("error: got %q, want nil", err)
					return
				}
			}

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("error: got nil, want %q", tt.wantErrMsg)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("error: got %q, want %q", err, tt.wantErr)
					return
				}
			}

			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("error message: got nil, want %q", tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error message: got %q, want %q", err, tt.wantErrMsg)
				}
			}

			if tt.wantOutput != "" {
				got := out.String()
				if !strings.Contains(got, tt.wantOutput) {
					t.Errorf("output: got %q, want %q", got, tt.wantOutput)
				}
			}

			if tt.wantOptions != nil {
				if !reflect.DeepEqual(argsQuotes, tt.wantOptions) {
					t.Errorf("options: got %v, want %v", argsQuotes, tt.wantOptions)
				}
			}
		})
	}
}
