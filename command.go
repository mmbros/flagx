package flagx

import (
	"errors"
	"os"
	"path"
	"strings"
)

// flagx defined inner errors
var (
	ErrInvalidCommandName = errors.New("invalid command name")
	ErrNoExecFunc         = errors.New("exec function undefined")
	ErrCommandNotFound    = errors.New("command not found")
)

// ParseExecFunc is the signature of the function that is called
// when a Command is executed.
//
// cmdfullname: the full name of the command (example: "app cmd subcmd")
// arguments:   the arguments of the command
type ParseExecFunc func(fullname string, arguments []string) error

// Command represents a node of the commands tree.
// Each node has the function to be called if the command is executed
// and the children sub-commands.
type Command struct {
	SubCmd    map[string]*Command // sub-commands of the command
	ParseExec ParseExecFunc       // function to be executed by the command
}

// handleSubCmd checks if the command must be executed
// or if a sub-command must be (recursivelly) called.
//
// fullname is the join of the ancestors or self command names, starting from root command.
// example: cmdfullname = "appname cmd1 subcmd11"
func (cmd *Command) handleSubCmd(fullname string, arguments []string) error {

	var arg0 string
	if len(arguments) > 0 {
		arg0 = arguments[0]
	}

	if arg0 == "" || strings.HasPrefix(arg0, "-") || (len(cmd.SubCmd) == 0) {
		// if no argument is passed
		// or the first argument begin with "-"
		// or the command has no subcommand
		// then parse the current command

		if cmd.ParseExec == nil {
			return wrapNameError(ErrNoExecFunc, fullname)
		}

		return cmd.ParseExec(fullname, arguments)
	}

	// arg0 must be the name of a sub command
	for names, subcmd := range cmd.SubCmd {
		ns := splitTrimSpace(names, ",")
		if len(ns) == 0 {
			return wrapNameErrorString(ErrInvalidCommandName, fullname, names)
		}

		if contains(ns, arg0) {
			// parse the subcommand
			return subcmd.handleSubCmd(fullname+" "+ns[0], arguments[1:])
		}
	}

	return wrapNameErrorString(ErrCommandNotFound, fullname, arg0)
}

// ParseExec execute the `root` Command with the command line arguments.
// The name of the `root` command is obtained from the `os.Args[0]` argument.
func ParseExec(root *Command) error {
	appname := path.Base(os.Args[0])

	return root.handleSubCmd(appname, os.Args[1:])
}
