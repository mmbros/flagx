package flagx_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	flagx "github.com/mmbros/flagx"
)

// const appname string = "app"

type options struct {
	move  bool
	path  string
	files []string
}

func (o *options) Clear() {
	*o = options{}
}

var opt = &options{}

var app = &flagx.Command{
	ParseExec: runApp,
	SubCmd: map[string]*flagx.Command{
		"import,imp,im,i": {
			ParseExec: runImport,
		},
	},
}

func runApp(name string, arguments []string) error {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <command>", name)
	}

	err := fs.Parse(arguments)
	return err
}

func runImport(name string, arguments []string) error {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <command>", name)
	}
	flagx.AliasedBoolVar(fs, &opt.move, "move,m", false, "move the files, instead of copying")
	flagx.AliasedStringVar(fs, &opt.path, "path,p", "", "path of the file")
	flagx.AliasedStringsVar(fs, &opt.files, "files,f", "files to import")

	err := fs.Parse(arguments)
	return err
}

func Test_App(t *testing.T) {

	tests := []struct {
		name        string
		args        string
		wantErr     error
		wantErrMsg  string
		wantOutput  string
		wantOptions *options
	}{
		{
			name: "empty command line",
			args: "",
		},
		{
			name:       "app help",
			args:       "-h",
			wantOutput: "usage: app",
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
			wantOutput: "usage: app",
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
			name:       "app import help",
			args:       "import -h",
			wantOutput: "usage: app import",
			wantErr:    flag.ErrHelp,
		},
		{
			name:       "app imp help",
			args:       "imp -h",
			wantOutput: "usage: app import",
			wantErr:    flag.ErrHelp,
		},
		{
			name:       "app im help",
			args:       "im -h",
			wantOutput: "usage: app import",
			wantErr:    flag.ErrHelp,
		},
		{
			name: "app import with options",
			args: "import --move -p /path/to/file",
			wantOptions: &options{
				move: true,
				path: "/path/to/file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// re-initilize global options
			opt.Clear()

			var out strings.Builder
			flag.CommandLine.SetOutput(&out)

			os.Args = strings.Split("app "+tt.args, " ")
			err := flagx.Run(app)

			if err != nil {
				if (tt.wantErr == nil) && (tt.wantErrMsg == "") {
					t.Errorf("error: got %q, want nil", err)
					return
				}
			}

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("error: got nil, want %q", tt.wantErr)
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
				if !reflect.DeepEqual(opt, tt.wantOptions) {
					t.Errorf("options: got %v, want %v", opt, tt.wantOptions)
				}
			}

		})
	}
}
