package flagx

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_AliasedFlagSet(t *testing.T) {
	var (
		str1   string
		str2   string
		int1   int
		int2   int
		bool1  bool
		bool2  bool
		int641 int64
		int642 int64
	)

	fs := flag.NewFlagSet("test", flag.ExitOnError)

	fs.StringVar(&str1, "str1", "", "usage str1")
	fs.IntVar(&int1, "int1", 0, "usage int1")
	AliasedStringVar(fs, &str2, "str2,s2", "valstr2", "usage str2")
	AliasedIntVar(fs, &int2, "int2,i2", 10, "usage int2")
	AliasedBoolVar(fs, &bool2, "bool2,b2,alias-bool2", true, "usage bool2")
	fs.BoolVar(&bool1, "bool1", false, "usage bool1")

	fs.Int64Var(&int641, "int641", 0, "usage int641")
	AliasedInt64Var(fs, &int642, "int642,i642", 99999999999, "usage int642")

	var buf strings.Builder

	fs.SetOutput(&buf)
	fs.PrintDefaults()

	got := buf.String()

	for _, want := range []string{"-str1 string", "-int1 int", "-str2 string", "-s2 string",
		"-int2 int", "-i2 int", "-bool2", "-b2", "-alias-bool2", "-bool1",
		"-int641 int", "-int642 int", "-i642 int"} {
		if !strings.Contains(got, want) {
			t.Errorf("PrintDefaults(): got %v, want substring %v", got, want)
		}
	}

}

func Test_AliasedStringsVar(t *testing.T) {

	tests := []struct {
		name        string
		args        string
		wantErr     error
		wantErrMsg  string
		wantOutput  string
		wantStrings []string
	}{
		{
			name: "empty command line",
			args: "",
		},
		{
			name:        "one arg with one string",
			args:        "-a s1",
			wantStrings: []string{"s1"},
		},
		{
			name:        "one arg with two strings",
			args:        "--astr s1,s2 ",
			wantStrings: []string{"s1", "s2"},
		},
		{
			name:        "two args with one string",
			args:        "-a s1 --astr s2 ",
			wantStrings: []string{"s1", "s2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var astr []string
			var out strings.Builder
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			fs.SetOutput(&out)

			AliasedStringsVar(fs, &astr, "a,astr", "usage astr")

			args := splitTrimSpace(tt.args, " ")
			err := fs.Parse(args)

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

			if tt.wantStrings != nil {
				if !reflect.DeepEqual(astr, tt.wantStrings) {
					t.Errorf("strings: got %v, want %v", astr, tt.wantStrings)
				}
			}
		})
	}
}

func Test_IsPassed_Int(t *testing.T) {

	const defValue = 99

	tests := []struct {
		name       string
		args       string
		flagNames  string
		wantPassed bool
		wantValue  int
	}{
		{
			name:       "empty command line",
			args:       "",
			flagNames:  "i,int",
			wantPassed: false,
			wantValue:  defValue,
		},
		{
			name:       "arg --other",
			args:       "--other 12",
			flagNames:  "i,int",
			wantPassed: false,
			wantValue:  defValue,
		},
		{
			name:       "arg --int",
			args:       fmt.Sprintf("--int %v", defValue),
			flagNames:  "i,int",
			wantPassed: true,
			wantValue:  defValue,
		},
		{
			name:       "arg -i",
			args:       "-i 12",
			flagNames:  "int,i",
			wantPassed: true,
			wantValue:  12,
		},
		{
			name:       "multiple win the last",
			args:       "-i 12 -o 13 --int 14 --other 15",
			flagNames:  "int,i",
			wantPassed: true,
			wantValue:  14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var value int
			var other int

			var out strings.Builder
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			fs.SetOutput(&out)

			AliasedIntVar(fs, &value, "i,int", defValue, "usage int")
			AliasedIntVar(fs, &other, "o,other", defValue, "usage other")

			args := splitTrimSpace(tt.args, " ")
			err := fs.Parse(args)

			if err != nil {
				t.Errorf("error: got %q, want nil", err)
				return
			}

			gotPassed := IsPassed(fs, tt.flagNames)

			if gotPassed != tt.wantPassed {
				t.Errorf("IsFlagPassed(%q): got %v, want %v", tt.flagNames, gotPassed, tt.wantPassed)
			}

			if value != tt.wantValue {
				t.Errorf("Value(%q): got %v, want %v", tt.flagNames, value, tt.wantValue)
			}
		})
	}
}
