package flagx

import (
	"errors"
	"flag"
	"reflect"
	"strings"
	"testing"
)

func Test_AliasedFlagSet(t *testing.T) {
	var (
		str1  string
		str2  string
		int1  int
		int2  int
		bool1 bool
		bool2 bool
	)

	fs := flag.NewFlagSet("test", flag.ExitOnError)

	fs.StringVar(&str1, "str1", "", "usage str1")
	fs.IntVar(&int1, "int1", 0, "usage int1")
	AliasedStringVar(fs, &str2, "str2,s2", "valstr2", "usage str2")
	AliasedIntVar(fs, &int2, "int2,i2", 10, "usage int2")
	AliasedBoolVar(fs, &bool2, "bool2,b2,alias-bool2", true, "usage bool2")
	fs.BoolVar(&bool1, "bool1", false, "usage bool1")

	var buf strings.Builder

	fs.SetOutput(&buf)
	fs.PrintDefaults()

	got := buf.String()

	for _, want := range []string{"-str1 string", "-int1 int", "-str2 string", "-s2 string",
		"-int2 int", "-i2 int", "-bool2", "-b2", "-alias-bool2", "-bool1"} {
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
