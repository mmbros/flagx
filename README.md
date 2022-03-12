# flagx

[![Build Status](https://travis-ci.com/mmbros/flagx.svg?branch=main)](https://app.travis-ci.com/github/mmbros/flagx)
[![GoDoc](https://godoc.org/github.com/mmbros/flagx?status.svg)](https://godoc.org/github.com/mmbros/flagx)
[![Report Card](https://goreportcard.com/badge/github.com/mmbros/flagx)](https://goreportcard.com/report/github.com/mmbros/flagx)

Package `flagx` implements some utilities for command-line flag parsing.
It extends the functionalities of the standard `flag` package.

Main features:

- sub commands management
- alias of command and flag names
- array of string flag type
- check if a flag was passed

For example the next code defines an `app` Command instance with a sub-command with name `action` and aliases `act`, `ac` and `a`. Note that only the names of the sub-commands are defined; the command name itself is not defined in the Command type. The name of the root command is obtained from the `os.Args[0]` parameter.

```golang
var app = &flagx.Command{
    ParseExec: runApp,
    SubCmd: map[string]*flagx.Command{
        "action,act,ac,a": {
            ParseExec: runAction,
        },
    },
}
```

To execute the `app` Command instance with the command line arguments call the `Run` module function.

```golang
err := flagx.Run(app)
```

The `Run` module function in turn calls the `ParseExec` function of the `app` command or the `action` sub-command based on the command-line arguments.
Each `ParseExec` function first parse the passed arguments, then execute the specific work.

```golang
func runAction(name string, arguments []string) error {
    var params []string
    fs := flag.NewFlagSet(name, flag.ContinueOnError)
    flagx.AliasedStringsVar(fs, &params, "params,p", "description of the parameters")

    err := fs.Parse(arguments)

    if err == nil {
        err = execAction(params) // execute the specific work
    }
    return err
}
```

The `AliasedStringsVar` function defines an array of strings flag with name `params` and alias `p`.
The command-line

```shell
app action --params str1,str2 -p str3 -params str4 -p str5,str6
```

will execute the `execAction` function with `[]string{"str1", "str2", "str3", "str4", "str5", "str6"}` strings.

See test for more informations and usage examples.
