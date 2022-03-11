package flagx

import (
	"reflect"
	"testing"
)

func Test_splitTrimSpace(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{"", ","},
			want: []string{},
		},
		{
			name: "one name",
			args: args{",  ,  , cmd1  ,,  ,  ", ","},
			want: []string{"cmd1"},
		},
		{
			name: "two names",
			args: args{"  ,cmd1 , ,, cmd2 ", ","},
			want: []string{"cmd1", "cmd2"},
		},
		{
			name: "three names",
			args: args{", ,   ,,  cmd1 ,   ,cmd2  ,, cmd3 ,,   ,  ", ","},
			want: []string{"cmd1", "cmd2", "cmd3"},
		},
		{
			name: "args",
			args: args{"-o option  --flag    -w   10  ", " "},
			want: []string{"-o", "option", "--flag", "-w", "10"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitTrimSpace(tt.args.s, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitTrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		a []string
		x string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "nil array and empty string",
			args: args{nil, ""},
			want: false,
		},
		{
			name: "empty array and empty string",
			args: args{[]string{}, ""},
			want: false,
		},
		{
			name: "array with empty string and empty string",
			args: args{[]string{""}, ""},
			want: true,
		},
		{
			name: "a",
			args: args{[]string{"str1"}, "str1"},
			want: true,
		},
		{
			name: "b",
			args: args{[]string{"str1", "str2"}, "str3"},
			want: false,
		},
		{
			name: "c",
			args: args{[]string{"str1", "str2"}, "STR1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.a, tt.args.x); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
