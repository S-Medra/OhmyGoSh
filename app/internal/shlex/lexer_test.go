package shlex

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Basic arguments",
			input:    `echo hello world`,
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "Multiple spaces",
			input:    `echo   multiple   spaces  `,
			expected: []string{"echo", "multiple", "spaces"},
		},
		{
			name:     "Single quotes",
			input:    `echo 'hello world'`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:    "Single quote with backslash (literal behavior)",
			input:   `echo 'it\'s'`,
			wantErr: true, // In shell, 'it\' ends the quote, leaving s' unmatched
		},
		{
			name:     "Double quotes",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "Escaped quote in double quotes",
			input:    `echo "A \" inside"`,
			expected: []string{"echo", `A " inside`},
		},
		{
			name:     "Escaped backslash in double quotes",
			input:    `echo "A \\ escapes itself"`,
			expected: []string{"echo", `A \ escapes itself`},
		},
		{
			name:     "Literal backslash in double quotes",
			input:    `echo "no escape \n here"`,
			expected: []string{"echo", `no escape \n here`},
		},
		{
			name:     "Backslash outside quotes",
			input:    `echo hello\ world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "Escaped spaces",
			input:    `echo before\ \ \ after`,
			expected: []string{"echo", "before   after"},
		},
		{
			name:     "Mixed quotes and concatenation",
			input:    `echo ""/""`,
			expected: []string{"echo", "/"},
		},
		{
			name:    "Unmatched double quote",
			input:   `echo "unmatched`,
			wantErr: true,
		},
		{
			name:    "Trailing backslash",
			input:   `echo hello\`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Split(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Split() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
