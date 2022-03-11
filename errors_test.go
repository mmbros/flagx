package flagx

import (
	"errors"
	"testing"
)

func Test_wrapNameError(t *testing.T) {

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name: "not-nil",
			err:  ErrCommandNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapNameError(tt.err, tt.name)

			if !errors.Is(err, tt.err) {
				t.Errorf("wrapNameError: got = %q, want %v", err, tt.err)
				return
			}
		})
	}
}

func Test_wrapNameErrorString(t *testing.T) {

	// err := wrapNameError(nil, "nil")
	// if err != nil {

	// }

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name: "not-nil",
			err:  ErrCommandNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapNameErrorString(tt.err, tt.name, "string")

			if !errors.Is(err, tt.err) {
				t.Errorf("wrapNameErrorString: got = %q, want %v", err, tt.err)
				return
			}
		})
	}
}

func Test_wrapErrorf(t *testing.T) {

	// err := wrapNameError(nil, "nil")
	// if err != nil {

	// }

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name: "not-nil",
			err:  ErrCommandNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapErrorf(tt.err, "%s: error %q", tt.name, tt.err)

			if !errors.Is(err, tt.err) {
				t.Errorf("wrapErrorf: got = %q, want %v", err, tt.err)
				return
			}
		})
	}
}
