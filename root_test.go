package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{name: "root help", want: []string{"Available Commands:", "btree"}},
		{name: "explicit help", args: []string{"--help"}, want: []string{"btree"}},
		{name: "demo help", args: []string{"btree", "--help"}, want: []string{"btrees btree"}},
		{name: "demo", args: []string{"btree"}, want: []string{
			"Insert 6\n[10]\n  [5 6]\n  [20]",
			"Search for 17\nvisit [10 20]: follow child 1",
			"visit [12 17]: found",
			"Search for 15", "leaf reached; not found",
			"Search for 10\nvisit [10 20]: found",
		}},
		{name: "bplus help", args: []string{"bplus", "--help"}, want: []string{"btrees bplus"}},
		{name: "bplus demo", args: []string{"bplus"}, want: []string{
			"branch [10 20]", "leaf [10 12 17]",
			"Leaf chain: [5 6 7] -> [10 12 17] -> [20 30] -> nil",
			"Range [6, 20] (inclusive)", "follow next leaf", "Result: [6 7 10 12 17 20]",
		}},
		{name: "bplus unexpected argument", args: []string{"bplus", "extra"}, wantErr: true},
		{name: "bplus unknown flag", args: []string{"bplus", "--missing"}, wantErr: true},
		{name: "unknown command", args: []string{"missing"}, wantErr: true},
		{name: "unexpected argument", args: []string{"btree", "extra"}, wantErr: true},
		{name: "unknown flag", args: []string{"btree", "--missing"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newRootCommand()
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			root.SetArgs(tt.args)
			err := root.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v; wantErr = %v", err, tt.wantErr)
			}
			for _, want := range tt.want {
				if !strings.Contains(output.String(), want) {
					t.Errorf("output missing %q: %s", want, output.String())
				}
			}
			if tt.name != "demo" && tt.name != "bplus demo" && strings.Contains(output.String(), "Insert 10") {
				t.Fatal("demonstration ran unexpectedly")
			}
		})
	}
}
