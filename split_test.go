package main

import (
	"reflect"
	"testing"
)

func TestSplitLinesTerminator(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{}},
		{"a", []string{"a"}},
		{"a\n", []string{"a"}},
		{"a\nb", []string{"a", "b"}},
		{"a\nb\n", []string{"a", "b"}},
		{"\n", []string{""}},
	}

	for _, tt := range tests {
		got := splitLinesTerminator(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("splitLinesTerminator(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
