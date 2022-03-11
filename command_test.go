package flagx

import (
	"errors"
	"strings"
	"testing"
)

var errApp = errors.New("error app")
var errCmd1 = errors.New("error 1")
var errCmd2 = errors.New("error 2")

func cmdAppExec(name string, arguments []string) error {
	return wrapNameError(errApp, name)
}

func cmdCmd1Exec(name string, arguments []string) error {
	return wrapNameError(errCmd1, name)
}

func cmdCmd2Exec(name string, arguments []string) error {
	return wrapNameError(errCmd2, name)
}

func TestCommand_parseExec(t *testing.T) {
	app := &Command{
		SubCmd: map[string]*Command{
			"cmd0,c0": {},
			"cmd1,alias11,alias12": {
				ParseExec: cmdCmd1Exec,
			},
			"cmd2,alias21,alias22": {
				ParseExec: cmdCmd2Exec,
				SubCmd: map[string]*Command{
					"sub21": {
						ParseExec: cmdCmd1Exec,
					},
				},
			},
		},
		ParseExec: cmdAppExec,
	}

	type args struct {
		name      string
		arguments string
	}
	tests := []struct {
		name       string
		cmd        *Command
		args       args
		wantErr    error
		wantErrMsg string
	}{
		{
			name:    "root",
			cmd:     app,
			args:    args{"app", ""},
			wantErr: errApp,
		},
		{
			name:       "unknown root subcommand",
			cmd:        app,
			args:       args{"app", "cmd-unknown"},
			wantErr:    ErrCommandNotFound,
			wantErrMsg: `app: command not found "cmd-unknown"`,
		},
		{
			name: "app cmd with invalid name ",
			cmd: &Command{
				SubCmd: map[string]*Command{
					"  ,  ,  ": {},
				},
			},
			args:       args{"app", "xxxx"},
			wantErr:    ErrInvalidCommandName,
			wantErrMsg: `app: invalid command name "  ,  ,  "`,
		},
		{
			name:       "app cmd with no exec function",
			cmd:        app,
			args:       args{"app", "cmd0"},
			wantErr:    ErrNoExecFunc,
			wantErrMsg: "app cmd0: exec function undefined",
		},
		{
			name:    "app cmd1",
			cmd:     app,
			args:    args{"app", "cmd1"},
			wantErr: errCmd1,
		},
		{
			name:    "app alias11",
			cmd:     app,
			args:    args{"app", "alias11"},
			wantErr: errCmd1,
		},
		{
			name:    "app cmd2",
			cmd:     app,
			args:    args{"app", "cmd2"},
			wantErr: errCmd2,
		},
		{
			name:    "app alias22",
			cmd:     app,
			args:    args{"app", "alias22"},
			wantErr: errCmd2,
		},
		{
			name:    "app alias21 sub21",
			cmd:     app,
			args:    args{"app", "alias21 sub21"},
			wantErr: errCmd1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := splitTrimSpace(tt.args.arguments, " ")
			err := tt.cmd.handleSubCmd(tt.args.name, args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Command.parseExec() error = %q, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrMsg != "" {
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Command.parseExec() error = %q, wantErrMsg %v", err, tt.wantErrMsg)
				}
			}
		})
	}
}
