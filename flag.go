package flagx

import (
	"flag"
	"fmt"
	"strings"
)

// IsPassed checks if flag was provided.
// It returns true if at least one alias of the flag was found.
// names is the comma separated aliases of the flag.

func IsPassed(fs *flag.FlagSet, names string) bool {
	found := false

	fs.Visit(func(f *flag.Flag) {
		anames := splitTrimSpace(names, ",")
		for _, name := range anames {
			if f.Name == name {
				found = true
				return
			}
		}
	})
	return found
}

// AliasedStringVar defines a string flag with specified names, default value, and usage string.
// The `names` argument is the comma separated aliases of the flag.
// The specified usage string is used for the primary flag name only.
// The usage string of a secondary flag name specifies that it is an alias of the primary name.
// The argument p points to a string variable in which to store the value of the flag.
func AliasedStringVar(fs *flag.FlagSet, p *string, names string, value string, usage string) {
	anames := splitTrimSpace(names, ",")
	for j, name := range anames {
		if j == 1 {
			// redefine usage for the aliased names
			usage = "alias of \"" + anames[0] + "\""
		}
		fs.StringVar(p, name, value, usage)
	}
}

// AliasedIntVar defines an int flag with specified names, default value, and usage string.
// The specified usage string is used for the primary flag name only.
// The usage string of a secondary flag name specifies that it is an alias of the primary name.
// The argument p points to an int variable in which to store the value of the flag.
func AliasedIntVar(fs *flag.FlagSet, p *int, names string, value int, usage string) {
	anames := splitTrimSpace(names, ",")
	for j, name := range anames {
		if j == 1 {
			// redefine usage for the aliased names
			usage = "alias of \"" + anames[0] + "\""
		}
		fs.IntVar(p, name, value, usage)
	}
}

// AliasedBoolVar defines a bool flag with specified names, default value, and usage string.
// The specified usage string is used for the primary flag name only.
// The usage string of a secondary flag name specifies that it is an alias of the primary name.
// The argument p points to a bool variable in which to store the value of the flag.
func AliasedBoolVar(fs *flag.FlagSet, p *bool, names string, value bool, usage string) {
	anames := splitTrimSpace(names, ",")
	for j, name := range anames {
		if j == 1 {
			// redefine usage for the aliased names
			usage = "alias of \"" + anames[0] + "\""
		}
		fs.BoolVar(p, name, value, usage)
	}
}

// astring type is an array of string implementing the flag.Value interface.
type astring []string

// String method of flag.Value interface.
func (o *astring) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*o, ", "))
}

// Set method of flag.Value interface.
func (o *astring) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		if s := strings.TrimSpace(v); s != "" {
			*o = append(*o, s)
		}
	}
	return nil
}

// AliasedStringsVar defines a []string flag with specified names, and usage string.
// The specified usage string is used for the primary flag name only.
// The usage string of a secondary flag name specifies that it is an alias of the primary name.
// The argument p points to a []string variable in which to store the value of the flag.
// Note: no default value is given.
func AliasedStringsVar(fs *flag.FlagSet, p *[]string, names string, usage string) {
	var ss *astring = (*astring)(p)
	anames := splitTrimSpace(names, ",")
	for j, name := range anames {
		if j == 1 {
			// redefine usage for the aliased names
			usage = "alias of \"" + anames[0] + "\""
		}
		fs.Var(ss, name, usage)
	}
}
